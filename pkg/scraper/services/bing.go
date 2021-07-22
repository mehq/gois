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
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
)

// BingConfig is a set of options used by BingScraper to perform/filter/format search results.
type BingConfig struct {
	AspectRatio  string
	Compact      bool
	ImageColor   string
	ImageSize    string
	ImageType    string
	PeopleFilter string
	Query        string
	SafeSearch   string
}

// BingScraper represents scraper for bing image search.
type BingScraper struct {
	Config *BingConfig
}

// BingResult is a set of attributes that defines an image result.
type BingResult struct {
	Height       int    `json:"-"`
	Width        int    `json:"-"`
	ReferenceURL string `json:"purl"`
	ThumbnailURL string `json:"turl"`
	Title        string `json:"t"`
	URL          string `json:"murl"`
}

func (b BingScraper) makeFilterString() (string, error) {
	filters := make([]string, 0)

	if b.Config.AspectRatio != "" {
		switch b.Config.AspectRatio {
		case params.AspectRatioSquare:
			filters = append(filters, "filterui:aspect-square")
		case params.AspectRationWide:
			filters = append(filters, "filterui:aspect-wide")
		case params.AspectRatioTall:
			filters = append(filters, "filterui:aspect-tall")
		default:
			return "", fmt.Errorf("--aspect-ratio: invalid value %s", b.Config.AspectRatio)
		}
	}

	if b.Config.ImageColor != "" {
		switch b.Config.ImageColor {
		case params.ColorFull:
			filters = append(filters, "filterui:color2-color")
		case params.ColorBlackAndWhite:
			filters = append(filters, "filterui:color2-bw")
		case params.ColorRed:
			filters = append(filters, "filterui:color2-FGcls_RED")
		case params.ColorOrange:
			filters = append(filters, "filterui:color2-FGcls_ORANGE")
		case params.ColorYellow:
			filters = append(filters, "filterui:color2-FGcls_YELLOW")
		case params.ColorGreen:
			filters = append(filters, "filterui:color2-FGcls_GREEN")
		case params.ColorTeal:
			filters = append(filters, "filterui:color2-FGcls_TEAL")
		case params.ColorBlue:
			filters = append(filters, "filterui:color2-FGcls_BLUE")
		case params.ColorPurple:
			filters = append(filters, "filterui:color2-FGcls_PURPLE")
		case params.ColorPink:
			filters = append(filters, "filterui:color2-FGcls_PINK")
		case params.ColorBrown:
			filters = append(filters, "filterui:color2-FGcls_BROWN")
		case params.ColorBlack:
			filters = append(filters, "filterui:color2-FGcls_BLACK")
		case params.ColorGray:
			filters = append(filters, "filterui:color2-FGcls_GRAY")
		case params.ColorWhite:
			filters = append(filters, "filterui:color2-FGcls_WHITE")
		default:
			return "", fmt.Errorf("--image-color: invalid value %s", b.Config.ImageColor)
		}
	}

	if b.Config.ImageSize != "" {
		switch b.Config.ImageSize {
		case params.ImageSizeSmall:
			filters = append(filters, "filterui:imagesize-small")
		case params.ImageSizeMedium:
			filters = append(filters, "filterui:imagesize-medium")
		case params.ImageSizeLarge:
			filters = append(filters, "filterui:imagesize-large")
		case params.ImageSizeExtraLarge:
			filters = append(filters, "filterui:imagesize-wallpaper")
		default:
			return "", fmt.Errorf("--image-size: invalid value %s", b.Config.ImageSize)
		}
	}

	if b.Config.ImageType != "" {
		switch b.Config.ImageType {
		case params.ImageTypePhoto:
			filters = append(filters, "filterui:photo-photo")
		case params.ImageTypeClipArt:
			filters = append(filters, "filterui:photo-clipart")
		case params.ImageTypeLineDrawing:
			filters = append(filters, "filterui:photo-linedrawing")
		case params.ImageTypeAnimated:
			filters = append(filters, "filterui:photo-animatedgif")
		case params.ImageTypeTransparent:
			filters = append(filters, "filterui:photo-transparent")
		default:
			return "", fmt.Errorf("--image-type: invalid value %s", b.Config.ImageType)
		}
	}

	if b.Config.PeopleFilter != "" {
		switch b.Config.PeopleFilter {
		case params.ImageTypeFace:
			filters = append(filters, "filterui:face-face")
		case params.OrientationPortrait:
			filters = append(filters, "filterui:face-portrait")
		default:
			return "", fmt.Errorf("--people-filter: invalid value %s", b.Config.PeopleFilter)
		}
	}

	return strings.Join(filters, "+"), nil
}

func (b BingScraper) setSafeSearchSetting() error {
	q := b.Config.Query

	page, err := util.DownloadWebpage(
		"https://www.bing.com/images/search",
		http.StatusOK,
		nil,
		map[string]string{
			"q": q,
		},
	)

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	r := bytes.NewReader(page)

	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	guid, exists := doc.Find("input#GUID").Attr("value")

	if !exists {
		return fmt.Errorf("guid not found")
	}

	ru, exists := doc.Find("input#ru").Attr("value")

	if !exists {
		return fmt.Errorf("ru not found")
	}

	var safeSearchOption string

	switch b.Config.SafeSearch {
	case params.SafeSearchOff:
		safeSearchOption = "off"
	case params.SafeSearchOn:
		safeSearchOption = "strict"
	case params.SafeSearchModerate:
		safeSearchOption = "demote"
	default:
		return fmt.Errorf("--safe-search: invalid value %s", b.Config.SafeSearch)
	}

	qParams := map[string]string{
		"pref_sbmt":    "1",
		"adlt_set":     safeSearchOption,
		"adlt_confirm": "1",
		"GUID":         guid,
		"is_child":     "0",
		"ru":           ru,
	}

	_, err = util.DownloadWebpage("https://www.bing.com/settings.aspx", http.StatusOK, nil, qParams)

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

// Scrape parses and returns the results from bing image search if successful. An error is returned otherwise.
func (b BingScraper) Scrape() ([]interface{}, int, error) {
	filter, err := b.makeFilterString()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	err = b.setSafeSearchSetting()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	paramFirst := -150
	qParams := map[string]string{
		"count": "150",
		"first": strconv.Itoa(paramFirst),
		"q":     b.Config.Query,
		"qft":   filter,
		"relp":  "150",
	}

	hasMore := true
	items := make([]interface{}, 0)
	itemsURLCache := make(map[string]bool)
	pages := 0

	for hasMore {
		newCount := 0
		pages += 1
		paramFirst += 150
		qParams["first"] = strconv.Itoa(paramFirst)
		page, err := util.DownloadWebpage(
			"https://www.bing.com/images/async",
			http.StatusOK,
			nil,
			qParams,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		r := bytes.NewReader(page)

		doc, err := goquery.NewDocumentFromReader(r)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		sel := doc.Find("a.iusc")

		for i := range sel.Nodes {
			element := sel.Eq(i)
			mAttr, exists := element.Attr("m")

			if !exists {
				return nil, 0, fmt.Errorf("attribute 'm' does not exist on target element")
			}

			br := &BingResult{}
			err = json.Unmarshal([]byte(strings.ReplaceAll(mAttr, "&quot;", "\"")), br)

			if err != nil {
				return nil, 0, fmt.Errorf(
					"cannot json.Unmarshal on match '%s'",
					strings.ReplaceAll(mAttr, "&quot;", "\""),
				)
			}

			hrefAttr, exists := element.Attr("href")

			if !exists {
				return nil, 0, fmt.Errorf("attribute 'href' does not exist on target element")
			}

			height, err := util.SearchRegex("exph=([0-9]+)", hrefAttr, "height", true)

			if err != nil {
				return nil, 0, fmt.Errorf("%v", err)
			}

			width, err := util.SearchRegex("expw=([0-9]+)", hrefAttr, "width", true)

			if err != nil {
				return nil, 0, fmt.Errorf("%v", err)
			}

			br.Height, _ = strconv.Atoi(height)
			br.Width, _ = strconv.Atoi(width)

			if !itemsURLCache[br.URL] {
				items = append(items, br)
				itemsURLCache[br.URL] = true
				newCount++
			}
		}

		hasMore = newCount > 0
	}

	return items, pages, nil
}
