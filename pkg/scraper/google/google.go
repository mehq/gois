package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/mzbaulhaque/gois/internal/util/conv"
	"github.com/mzbaulhaque/gois/internal/util/request"
)

// Config is a set of options used by Google.
type Config struct {
	Compact  bool
	Explicit bool
	GIF      bool
	Gray     bool
	Query    string
	Height   int
	Width    int
}

// Google is used to scrape data from google search engine.
type Google struct {
	Client *http.Client
	Config *Config
}

type imageInfo struct {
	Height       int    `json:"oh"`
	Width        int    `json:"ow"`
	ReferenceURL string `json:"ru"`
	ThumbnailURL string `json:"tu"`
	Title        string `json:"pt"`
	URL          string `json:"ou"`
}

func (g *Google) init() error {
	if g.Client == nil {
		g.Client = request.NewHTTPClient()
	}

	if g.Config == nil {
		return errors.New("empty configuration for bing scraper")
	}

	return nil
}

func (g Google) makeFilterString() string {
	filters := make([]string, 0)

	if g.Config.GIF {
		filters = append(filters, "ift:gif")
	}

	if g.Config.Gray {
		filters = append(filters, "ic:gray")
	}

	if g.Config.Height > 0 && g.Config.Width > 0 {
		filters = append(filters, "isz:ex")
		filters = append(filters, fmt.Sprintf("iszh:%d", g.Config.Height))
		filters = append(filters, fmt.Sprintf("iszw:%d", g.Config.Width))
	}

	return strings.Join(filters, ",")
}

// Scrape is the entrypoint.
func (g Google) Scrape() ([]interface{}, error) {
	err := g.init()
	if err != nil {
		return nil, err
	}

	params := &url.Values{}
	params.Set("tbm", "isch")
	params.Set("asearch", "ichunk")
	params.Set("ijn", "0")
	params.Set("start", "0")
	params.Set("q", g.Config.Query)
	params.Set("hl", "en")
	params.Set("async", "_id:rg_s,_pms:s,_fmt:pc")
	params.Set("tbs", g.makeFilterString())

	if !g.Config.Explicit {
		params.Set("safe", "active")
	} else {
		params.Set("safe", "images")
	}

	hasMore := true
	items := make([]interface{}, 0)
	itemsURLCache := make(map[string]bool)

	for hasMore {
		newCount := 0
		res, err := g.Client.Do(request.NewRequest("GET", "https://www.google.com/search", params))

		if err != nil {
			return nil, fmt.Errorf("failed to make request to bing")
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			return nil, fmt.Errorf("failed to parse bing response")
		}

		sel := doc.Find("div.rg_meta.notranslate")

		for i := range sel.Nodes {
			element := sel.Eq(i)
			rawInfo := element.Text()
			parsedInfo := &imageInfo{}
			err = json.Unmarshal([]byte(rawInfo), parsedInfo)

			if err != nil {
				panic(err)
			}

			if !itemsURLCache[parsedInfo.URL] {
				items = append(items, parsedInfo)
				itemsURLCache[parsedInfo.URL] = true
				newCount++
			}
		}

		_ = res.Body.Close()
		hasMore = newCount > 0

		params.Set("ijn", conv.Itoa(conv.Atoi(params.Get("ijn"))+1))
		params.Set("start", conv.Itoa(conv.Atoi(params.Get("ijn"))*100))
	}

	return items, nil
}
