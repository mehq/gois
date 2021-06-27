package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mzbaulhaque/gomage/internal"
	"net/http"
	"net/url"
	"strings"
)

// Google is used to scrape data from google search engine.
type Google struct {
	Client *http.Client
	Opts   *internal.Options
}

type googleInfoItem struct {
	URL string `json:"ou"`
}

func (g Google) makeFilterString() string {
	filters := make([]string, 0)

	if g.Opts.Gif {
		filters = append(filters, "ift:gif")
	}

	if g.Opts.Gray {
		filters = append(filters, "ic:gray")
	}

	if g.Opts.Height > 0 && g.Opts.Width > 0 {
		filters = append(filters, "isz:ex")
		filters = append(filters, fmt.Sprintf("iszh:%d", g.Opts.Height))
		filters = append(filters, fmt.Sprintf("iszw:%d", g.Opts.Width))
	}

	return strings.Join(filters, ",")
}

// Scrape is the entrypoint.
func (g Google) Scrape() []string {
	params := &url.Values{}
	params.Set("tbm", "isch")
	params.Set("asearch", "ichunk")
	params.Set("ijn", "0")
	params.Set("start", "0")
	params.Set("q", g.Opts.Query)
	params.Set("hl", "en")
	params.Set("async", "_id:rg_s,_pms:s,_fmt:pc")
	params.Set("tbs", g.makeFilterString())

	if g.Opts.Safe {
		params.Set("safe", "active")
	} else {
		params.Set("safe", "images")
	}

	hasMore := true
	itemCache := make(map[string]bool)
	items := make([]string, 0)

	for hasMore {
		newCount := 0
		res, err := g.Client.Do(internal.MakeRequest("GET", "https://www.google.com/search", params, nil))

		if err != nil {
			panic(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			panic(err)
		}

		doc.Find("div.rg_meta.notranslate").Each(func(_ int, element *goquery.Selection) {
			rawInfo := element.Text()
			parsedInfo := &googleInfoItem{}
			err = json.Unmarshal([]byte(rawInfo), parsedInfo)

			if err != nil {
				panic(err)
			}

			if !itemCache[parsedInfo.URL] {
				items = append(items, parsedInfo.URL)
				itemCache[parsedInfo.URL] = true
				newCount++
			}
		})

		_ = res.Body.Close()

		if g.Opts.TestMode {
			break
		}

		hasMore = newCount > 0
		params.Set("ijn", internal.Itoa(internal.Atoi(params.Get("ijn"))+1))
		params.Set("start", internal.Itoa(internal.Atoi(params.Get("ijn"))*100))
	}

	return items
}
