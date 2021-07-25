package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/mzbaulhaque/gois/internal/util"
)

// YahooConfig is a set of options used by YahooScraper to perform/filter/format search results.
type YahooConfig struct {
	AspectRatio string
	Compact     bool
	ImageColor  string
	ImageSize   string
	ImageType   string
	Query       string
	SafeSearch  string
}

// YahooScraper represents scraper for yahoo image search.
type YahooScraper struct {
	Config *YahooConfig
}

// YahooResult is a set of attributes that defines an image result.
type YahooResult struct {
	Height       int    `json:"-"`
	Width        int    `json:"-"`
	RawHeight string `json:"h"`
	RawWidth string `json:"w"`
	ReferenceURL string `json:"rurl"`
	ThumbnailURL string `json:"ith"`
	Title        string `json:"alt"`
	URL          string `json:"ourl"`
}

type yahooResponseMeta struct {
	Count int `json:"count"`
	First int `json:"first"`
	Last int `json:"last"`
	Total int `json:"total"`
}

type yahooResponse struct {
	HTML string `json:"html"`
	Meta yahooResponseMeta `json:"meta"`
}

func (y YahooScraper) makeFilterString() (string, error) {
	filters := make([]string, 0)

	return strings.Join(filters, ","), nil
}

// Scrape parses and returns the results from yahoo image search if successful. An error is returned otherwise.
func (y YahooScraper) Scrape() ([]interface{}, int, error) {
	qParams := map[string]string{
		"b": "0",
		"ei": "UTF-8",
		"fr": "sfp",
		"fr2": "sb-top-images.search",
		"o": "js",
		"p":       y.Config.Query,
	}

	hasMore := true
	items := make([]interface{}, 0)
	itemsURLCache := make(map[string]bool)
	pages := 0

	switch y.Config.SafeSearch {
	case "", params.SafeSearchOn:
		util.SetCookie(
			"sB",
			"vm=r",
			"/",
			".images.search.yahoo.com",
			"https://images.search.yahoo.com/search/images",
		)
	case params.SafeSearchOff:
		util.SetCookie(
			"sB",
			"vm=p",
			"/",
			".images.search.yahoo.com",
			"https://images.search.yahoo.com/search/images",
		)
	}

	for hasMore {
		newCount := 0
		pages += 1
		qParams["b"] = strconv.Itoa(((pages - 1) * 60) + 1)
		page, err := util.DownloadWebpage("https://images.search.yahoo.com/search/images", http.StatusOK, nil, qParams)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		res := &yahooResponse{}
		err = json.Unmarshal(page, res)

		if err != nil {
			return nil, 0, fmt.Errorf("cannot json.Unmarshal on response")
		}

		r := bytes.NewReader([]byte(res.HTML))

		doc, err := goquery.NewDocumentFromReader(r)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		sel := doc.Find("li.ld")

		for i := range sel.Nodes {
			element := sel.Eq(i)
			dataAttr, exists := element.Attr("data")

			if !exists {
				return nil, 0, fmt.Errorf("data attribute does not exist")
			}

			yr := &YahooResult{}
			err = json.Unmarshal([]byte(dataAttr), yr)

			if err != nil {
				return nil, 0, fmt.Errorf("cannot json.Unmarshal on data attribute")
			}

			yr.Height, err = strconv.Atoi(yr.RawHeight)

			if err != nil {
				return nil, 0, fmt.Errorf("cannot convert height '%s' to integer", yr.RawHeight)
			}

			yr.Width, err = strconv.Atoi(yr.RawWidth)

			if err != nil {
				return nil, 0, fmt.Errorf("cannot convert width '%s' to integer", yr.RawWidth)
			}

			if !itemsURLCache[yr.URL] {
				items = append(items, yr)
				itemsURLCache[yr.URL] = true
				newCount++
			}
		}

		hasMore = res.Meta.Last < res.Meta.Total
	}

	return items, pages, nil
}
