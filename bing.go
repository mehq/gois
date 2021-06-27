package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type bingInfoItem struct {
	URL string `json:"murl"`
}

// Bing is used to scrape data from bing search engine.
type Bing struct {
	client *http.Client
	opts   *Options
}

func (b Bing) makeFilterString() string {
	filters := make([]string, 0)

	if b.opts.gif {
		filters = append(filters, "filterui:photo-animatedgif")
	}

	if b.opts.gray {
		filters = append(filters, "filterui:color2-bw")
	}

	if b.opts.height > 0 && b.opts.width > 0 {
		fmt.Println("setting height")
		filters = append(filters, fmt.Sprintf("filterui:imagesize-custom_%d_%d", b.opts.width, b.opts.height))
	}

	return strings.Join(filters, "+")
}

func (b Bing) turnSafeSearchOff() {
	params := &url.Values{}
	params.Set("q", b.opts.query)

	res, err := b.client.Do(MakeRequest("GET", "https://www.bing.com/images/search", params, nil))

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

	_, err = b.client.Do(MakeRequest("GET", "https://www.bing.com/settings.aspx", params, nil))

	if err != nil {
		panic(err)
	}
}

// Scrape is the entrypoint.
func (b Bing) Scrape() []string {
	if !b.opts.safe {
		b.turnSafeSearchOff()
	}

	params := &url.Values{}
	params.Set("q", b.opts.query)
	params.Set("first", "0")
	params.Set("count", "150")
	params.Set("relp", "150")
	params.Set("qft", b.makeFilterString())

	hasMore := true
	itemCache := make(map[string]bool)
	items := make([]string, 0)

	for hasMore {
		newCount := 0
		res, err := b.client.Do(MakeRequest("GET", "https://www.bing.com/images/async", params, nil))

		if err != nil {
			panic(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			panic(err)
		}

		doc.Find("a.iusc").Each(func(_ int, element *goquery.Selection) {
			rawInfo, exists := element.Attr("m")

			if !exists {
				panic("Does not exist")
			}

			parsedInfo := &bingInfoItem{}

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

		if b.opts.testMode {
			break
		}

		hasMore = newCount > 0
		params.Set("first", Itoa(Atoi(params.Get("first"))+150))
	}

	return items
}
