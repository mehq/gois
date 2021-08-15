package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
)

type YandexConfig struct {
	Compact     bool
	ImageColor  string
	ImageSize   string
	ImageType   string
	Orientation string
	Query       string
	SafeSearch  string
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

func (y YandexScraper) makeFilters() (map[string]string, error) {
	filters := map[string]string{}

	switch y.Config.ImageColor {
	case "", params.ParamAll:
	case params.ColorFull:
		filters["icolor"] = "color"
	case params.ColorBlackAndWhite:
		filters["icolor"] = "gray"
	case params.ColorRed:
		filters["icolor"] = "red"
	case params.ColorOrange:
		filters["icolor"] = "orange"
	case params.ColorYellow:
		filters["icolor"] = "yellow"
	case params.ColorCyan:
		filters["icolor"] = "cyan"
	case params.ColorGreen:
		filters["icolor"] = "green"
	case params.ColorBlue:
		filters["icolor"] = "blue"
	case params.ColorViolet:
		filters["icolor"] = "violet"
	case params.ColorWhite:
		filters["icolor"] = "white"
	case params.ColorBlack:
		filters["icolor"] = "black"
	default:
		return nil, fmt.Errorf("--image-color: invalid value %s", y.Config.ImageColor)
	}

	switch y.Config.ImageSize {
	case "", params.ParamAll:
	case params.ImageSizeSmall:
		filters["isize"] = "small"
	case params.ImageSizeMedium:
		filters["isize"] = "medium"
	case params.ImageSizeLarge:
		filters["isize"] = "large"
	default:
		return nil, fmt.Errorf("--image-size: invalid value %s", y.Config.ImageSize)
	}

	switch y.Config.ImageType {
	case "", params.ParamAll:
	case params.ImageTypePhoto:
		filters["type"] = "photo"
	case params.ImageTypeClipArt:
		filters["type"] = "clipart"
	case params.ImageTypeLineDrawing:
		filters["type"] = "lineart"
	case params.ImageTypeFace:
		filters["type"] = "face"
	case params.ImageTypeDemotivational:
		filters["type"] = "demotivator"
	default:
		return nil, fmt.Errorf("--image-type: invalid value %s", y.Config.ImageType)
	}

	switch y.Config.Orientation {
	case "", params.ParamAll:
	case params.OrientationLandscape:
		filters["iorient"] = "horizontal"
	case params.OrientationPortrait:
		filters["iorient"] = "vertical"
	case params.AspectRatioSquare:
		filters["iorient"] = "square"
	default:
		return nil, fmt.Errorf("--orientation: invalid value %s", y.Config.Orientation)
	}

	return filters, nil
}

func (y YandexScraper) setSafeSearchPreference() error {
	q := y.Config.Query

	page, err := util.DownloadWebpage(
		"https://yandex.com/images/search",
		http.StatusOK,
		nil,
		map[string]string{
			"text": q,
		},
	)

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	var safeSearchValue string

	switch y.Config.SafeSearch {
	case params.SafeSearchOff:
		safeSearchValue = "0"
	case "", params.SafeSearchOn:
		safeSearchValue = "2"
	case params.SafeSearchModerate:
		safeSearchValue = "1"
	default:
		return fmt.Errorf("--safe-search: invalid value %s", y.Config.SafeSearch)
	}

	yandexUID, err := util.SearchRegex(`"yandexuid"\s*:\s*"([^"]+)"`, string(page), "yandexuid")

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	qParams := map[string]string{
		"save":      "1",
		"retpath":   fmt.Sprintf("https://yandex.com/images/search?text=%s", url.QueryEscape(q)),
		"yandexuid": yandexUID,
		"family":    safeSearchValue,
	}

	_, err = util.DownloadWebpage("https://yandex.com/images/customize", http.StatusOK, nil, qParams)

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (y YandexScraper) Scrape() ([]interface{}, int, error) {
	filters, err := y.makeFilters()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	err = y.setSafeSearchPreference()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	qParams := map[string]string{
		"format": "json",
		"rpt":    "image",
		"text":   y.Config.Query,
		"p":      "",
		"request": `{"blocks":[{"block":"serp-controller","params":{},"version":2},{"block":"serp-list_infinite_yes",` +
			`"params":{"initialPageNum":0},"version":2}]}`,
	}

	for k, v := range filters {
		qParams[k] = v
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
