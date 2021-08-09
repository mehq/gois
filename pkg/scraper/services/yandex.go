package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/mzbaulhaque/gois/internal/util"
)

type YandexConfig struct {
	Compact bool
	Query   string
}

type YandexScraper struct {
	Config *YandexConfig
}

type YandexResult struct {
	Height       int
	Width        int
	RawHeight    string
	RawWidth     string
	ReferenceURL string
	ThumbnailURL string
	Title        string
	URL          string
}

type yandexImage struct {
	Height int    `json:"h"`
	URL    string `json:"url"`
	Width  int    `json:"w"`
}

type yandexRawResult struct {
	SerpItem struct {
		Duplicates []yandexImage `json:"dups"`
		Previews   []yandexImage `json:"preview"`
		Snippet    struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"snippet"`
		Thumb struct {
			URL string `json:"url"`
		} `json:"thumb"`
		URL string `json:"img_href"`
	} `json:"serp-item"`
}

type yandexResponse struct {
	Blocks []struct {
		Params struct {
			LastPage int `json:"lastPage"`
		} `json:"params"`
		HTML string `json:"html"`
	} `json:"blocks"`
}

func (y YandexScraper) Scrape() ([]interface{}, int, error) {
	qParams := map[string]string{
		"format": "json",
		"rpt":    "image",
		"text":   y.Config.Query,
		"p":      "",
		"request": `{"blocks":[{"block":"serp-controller","params":{},"version":2},{"block":"serp-list_infinite_yes",` +
			`"params":{"initialPageNum":0},"version":2}]}`,
	}
	hasMore := true
	items := make([]interface{}, 0)
	itemsURLCache := make(map[string]bool)
	pages := 0

	for hasMore {
		newCount := 0
		pages += 1
		qParams["p"] = strconv.Itoa(pages - 1)
		page, err := util.DownloadWebpage(
			"https://yandex.com/images/search",
			http.StatusOK,
			nil,
			qParams,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		res := &yandexResponse{}
		err = json.Unmarshal(page, res)

		if err != nil {
			return nil, 0, fmt.Errorf("cannot json.Unmarshal on response")
		}

		serpListBlock := res.Blocks[1]

		r := bytes.NewReader([]byte(strings.ReplaceAll(serpListBlock.HTML, "[object Object]", "")))

		doc, err := goquery.NewDocumentFromReader(r)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		sel := doc.Find("div.serp-item")

		for i := range sel.Nodes {
			element := sel.Eq(i)
			bemAttr, exists := element.Attr("data-bem")

			if !exists {
				return nil, 0, fmt.Errorf("attribute 'data-bem' does not exist on target element")
			}

			rawResult := &yandexRawResult{}
			err := json.Unmarshal([]byte(bemAttr), rawResult)

			if err != nil {
				return nil, 0, fmt.Errorf("cannot json.Unmarshal on 'data-bem' attribute")
			}

			if !itemsURLCache[rawResult.SerpItem.URL] {
				height := 0
				width := 0
				downloadURL := rawResult.SerpItem.URL

				for _, image := range append(rawResult.SerpItem.Previews, rawResult.SerpItem.Duplicates...) {
					if image.URL == downloadURL {
						height = image.Height
						width = image.Width

						break
					}
				}

				items = append(items, YandexResult{
					Height:       height,
					ReferenceURL: rawResult.SerpItem.Snippet.URL,
					ThumbnailURL: "https:" + rawResult.SerpItem.Thumb.URL,
					Title:        rawResult.SerpItem.Snippet.Title,
					URL:          downloadURL,
					Width:        width,
				})
				itemsURLCache[rawResult.SerpItem.URL] = true
				newCount++
			}
		}

		hasMore = pages-1 < serpListBlock.Params.LastPage
	}

	return items, pages, nil
}
