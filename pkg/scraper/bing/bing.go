package bing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/mzbaulhaque/gois/internal/util/conv"
	"github.com/mzbaulhaque/gois/internal/util/request"
)

// Config is a set of options used by Bing.
type Config struct {
	Compact  bool
	Explicit bool
	GIF      bool
	Gray     bool
	Query    string
	Height   int
	Width    int
}

// Bing is used to scrape data from bing search engine.
type Bing struct {
	Client *http.Client
	Config *Config
}

type imageInfo struct {
	Height       int    `json:"-"`
	Width        int    `json:"-"`
	ReferenceURL string `json:"purl"`
	ThumbnailURL string `json:"turl"`
	Title        string `json:"t"`
	URL          string `json:"murl"`
}

var heightWidthRe = regexp.MustCompile(`exph=([0-9]+)[^0-9]+expw=([0-9]+)[^0-9]`)

func (b *Bing) init() error {
	if b.Client == nil {
		b.Client = request.NewHTTPClient()
	}

	if b.Config == nil {
		return errors.New("empty configuration for bing scraper")
	}

	return nil
}

func (b Bing) makeFilterString() string {
	filters := make([]string, 0)

	if b.Config.GIF {
		filters = append(filters, "filterui:photo-animatedgif")
	}

	if b.Config.Gray {
		filters = append(filters, "filterui:color2-bw")
	}

	if b.Config.Height > 0 && b.Config.Width > 0 {
		filters = append(filters, fmt.Sprintf("filterui:imagesize-custom_%d_%d", b.Config.Width, b.Config.Height))
	}

	return strings.Join(filters, "+")
}

func (b Bing) turnSafeSearchOff() {
	params := &url.Values{}
	params.Set("q", b.Config.Query)

	res, err := b.Client.Do(request.NewRequest("GET", "https://www.bing.com/images/search", params))

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			panic(err)
		}
	}(res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		panic(err)
	}

	guid, exists := doc.Find("input#GUID").Attr("value")

	if !exists {
		panic("guid not found")
	}

	ru, exists := doc.Find("input#ru").Attr("value")

	if !exists {
		panic("ru not found")
	}

	params.Del("q")
	params.Set("pref_sbmt", "1")
	params.Set("adlt_set", "off")
	params.Set("adlt_confirm", "1")
	params.Set("GUID", guid)
	params.Set("is_child", "0")
	params.Set("ru", ru)

	_, err = b.Client.Do(request.NewRequest("GET", "https://www.bing.com/settings.aspx", params))

	if err != nil {
		panic(err)
	}
}

// Scrape is the entrypoint.
func (b Bing) Scrape() ([]interface{}, error) {
	err := b.init()
	if err != nil {
		return nil, err
	}

	if b.Config.Explicit {
		b.turnSafeSearchOff()
	}

	params := &url.Values{}
	params.Set("q", b.Config.Query)
	params.Set("first", "0")
	params.Set("count", "150")
	params.Set("relp", "150")
	params.Set("qft", b.makeFilterString())

	hasMore := true
	items := make([]interface{}, 0)
	itemsURLCache := make(map[string]bool)

	for hasMore {
		newCount := 0
		res, err := b.Client.Do(request.NewRequest("GET", "https://www.bing.com/images/async", params))

		if err != nil {
			return nil, fmt.Errorf("failed to make request to bing")
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			return nil, fmt.Errorf("failed to parse bing response")
		}

		sel := doc.Find("a.iusc")

		for i := range sel.Nodes {
			element := sel.Eq(i)
			rawInfo, exists := element.Attr("m")

			if !exists {
				panic("Does not exist")
			}

			parsedInfo := &imageInfo{}

			if !b.Config.Compact {
				href, exists := element.Attr("href")

				if !exists {
					panic("href not found")
				}

				h, w := parseHeightAndWidth(href)
				parsedInfo.Height = h
				parsedInfo.Width = w
			}

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

		params.Set("first", conv.Itoa(conv.Atoi(params.Get("first"))+150))
	}

	return items, nil
}

func parseHeightAndWidth(href string) (int, int) {
	matches := heightWidthRe.FindAllStringSubmatch(href, -1)

	if len(matches) != 1 {
		panic("parseHeightAndWidth() error")
	}

	height := conv.Atoi(matches[0][1])
	width := conv.Atoi(matches[0][2])

	return height, width
}
