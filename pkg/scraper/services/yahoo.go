package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"github.com/mehq/gois/pkg/scraper/params"

	"github.com/mehq/gois/internal/util"
)

// YahooConfig is a set of options used by YahooScraper to perform/filter/format search results.
type YahooConfig struct {
	Compact    bool
	ImageColor string
	ImageSize  string
	ImageType  string
	Query      string
	SafeSearch string
}

// YahooScraper represents scraper for yahoo image search.
type YahooScraper struct {
	Config *YahooConfig
}

// YahooResult is a set of attributes that defines an image result.
type YahooResult struct {
	Height       int    `json:"-"`
	Width        int    `json:"-"`
	RawHeight    string `json:"h"`
	RawWidth     string `json:"w"`
	ReferenceURL string `json:"rurl"`
	ThumbnailURL string `json:"ith"`
	Title        string `json:"alt"`
	URL          string `json:"ourl"`
}

type yahooResponseMeta struct {
	Count int `json:"count"`
	First int `json:"first"`
	Last  int `json:"last"`
	Total int `json:"total"`
}

type yahooResponse struct {
	HTML string            `json:"html"`
	Meta yahooResponseMeta `json:"meta"`
}

func (y YahooScraper) makeFilters() (map[string]string, error) {
	filters := map[string]string{}

	switch y.Config.ImageColor {
	case "", params.ParamAll:
	case params.ColorBlackAndWhite:
		filters["imgc"] = "bw"
	case params.ColorRed:
		filters["imgc"] = "red"
	case params.ColorOrange:
		filters["imgc"] = "orange"
	case params.ColorYellow:
		filters["imgc"] = "yellow"
	case params.ColorGreen:
		filters["imgc"] = "green"
	case params.ColorTeal:
		filters["imgc"] = "teal"
	case params.ColorBlue:
		filters["imgc"] = "blue"
	case params.ColorPurple:
		filters["imgc"] = "purple"
	case params.ColorPink:
		filters["imgc"] = "pink"
	case params.ColorWhite:
		filters["imgc"] = "white"
	case params.ColorGray:
		filters["imgc"] = "gray"
	case params.ColorBlack:
		filters["imgc"] = "black"
	case params.ColorBrown:
		filters["imgc"] = "brown"
	default:
		return nil, fmt.Errorf("--image-color: invalid value %s", y.Config.ImageColor)
	}

	switch y.Config.ImageSize {
	case "", params.ParamAll:
	case params.ImageSizeSmall:
		filters["imgsz"] = "small"
	case params.ImageSizeMedium:
		filters["imgsz"] = "medium"
	case params.ImageSizeLarge:
		filters["imgsz"] = "large"
	default:
		return nil, fmt.Errorf("--image-size: invalid value %s", y.Config.ImageSize)
	}

	switch y.Config.ImageType {
	case "", params.ParamAll:
	case params.ImageTypePhoto:
		filters["imgty"] = "photo"
	case params.ImageTypeGraphic:
		filters["imgty"] = "graphics"
	case params.ImageTypeAnimated:
		filters["imgty"] = "gif"
	case params.ImageTypeFace:
		filters["imgty"] = "face"
	case params.OrientationPortrait:
		filters["imgty"] = "portrait"
	case params.ImageTypeNonPortrait:
		filters["imgty"] = "nonportrait"
	case params.ImageTypeClipArt:
		filters["imgty"] = "clipart"
	case params.ImageTypeLineDrawing:
		filters["imgty"] = "linedrawing"
	default:
		return nil, fmt.Errorf("--image-type: invalid value %s", y.Config.ImageType)
	}

	return filters, nil
}

// Scrape parses and returns the results from yahoo image search if successful. An error is returned otherwise.
func (y YahooScraper) Scrape() ([]interface{}, int, error) {
	filters, err := y.makeFilters()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	qParams := map[string]string{
		"b":   "0",
		"ei":  "UTF-8",
		"fr":  "sfp",
		"fr2": "sb-top-images.search",
		"o":   "js",
		"p":   y.Config.Query,
	}

	for k, v := range filters {
		qParams[k] = v
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
