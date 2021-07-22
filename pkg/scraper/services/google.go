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

// GoogleConfig is a set of options used by GoogleScraper to perform/filter/format search results.
type GoogleConfig struct {
	AspectRatio string
	Compact     bool
	ImageColor  string
	ImageSize   string
	ImageType   string
	Query       string
	SafeSearch  string
}

// GoogleScraper represents scraper for google image search.
type GoogleScraper struct {
	Config *GoogleConfig
}

// GoogleResult is a set of attributes that defines an image result.
type GoogleResult struct {
	Height       int    `json:"oh"`
	Width        int    `json:"ow"`
	ReferenceURL string `json:"ru"`
	ThumbnailURL string `json:"tu"`
	Title        string `json:"pt"`
	URL          string `json:"ou"`
}

func (g GoogleScraper) makeFilterString() (string, error) {
	filters := make([]string, 0)

	if g.Config.AspectRatio != "" {
		switch g.Config.AspectRatio {
		case params.AspectRatioTall:
			filters = append(filters, "iar:t")
		case params.AspectRatioSquare:
			filters = append(filters, "iar:s")
		case params.AspectRationWide:
			filters = append(filters, "iar:w")
		case params.AspectRatioPanoramic:
			filters = append(filters, "iar:xw")
		default:
			return "", fmt.Errorf("--aspect-ratio: invalid value %s", g.Config.AspectRatio)
		}
	}

	if g.Config.ImageColor != "" {
		switch g.Config.ImageColor {
		case params.ColorFull:
			filters = append(filters, "ic:color")
		case params.ColorBlackAndWhite:
			filters = append(filters, "ic:gray")
		case params.ImageTypeTransparent:
			filters = append(filters, "ic:trans")
		case params.ColorRed:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:red")
		case params.ColorOrange:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:orange")
		case params.ColorYellow:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:yellow")
		case params.ColorGreen:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:green")
		case params.ColorTeal:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:teal")
		case params.ColorBlue:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:blue")
		case params.ColorPurple:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:purple")
		case params.ColorPink:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:pink")
		case params.ColorWhite:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:white")
		case params.ColorGray:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:gray")
		case params.ColorBlack:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:black")
		case params.ColorBrown:
			filters = append(filters, "ic:specific")
			filters = append(filters, "isc:brown")
		default:
			return "", fmt.Errorf("--image-color: invalid value %s", g.Config.ImageColor)
		}
	}

	if g.Config.ImageSize != "" {
		switch g.Config.ImageSize {
		case params.ImageSizeLarge:
			filters = append(filters, "isz:l")
		case params.ImageSizeMedium:
			filters = append(filters, "isz:m")
		case params.ImageSizeIcon:
			filters = append(filters, "isz:i")
		default:
			return "", fmt.Errorf("--image-size: invalid value %s", g.Config.ImageSize)
		}
	}

	if g.Config.ImageType != "" {
		switch g.Config.ImageType {
		case params.ImageTypeFace:
			filters = append(filters, "itp:face")
		case params.ImageTypePhoto:
			filters = append(filters, "itp:photo")
		case params.ImageTypeClipArt:
			filters = append(filters, "itp:clipart")
		case params.ImageTypeLineDrawing:
			filters = append(filters, "itp:lineart")
		case params.ImageTypeAnimated:
			filters = append(filters, "itp:animated")
		default:
			return "", fmt.Errorf("--image-type: invalid value %s", g.Config.ImageType)
		}
	}

	return strings.Join(filters, ","), nil
}

// Scrape parses and returns the results from google image search if successful. An error is returned otherwise.
func (g GoogleScraper) Scrape() ([]interface{}, int, error) {
	filter, err := g.makeFilterString()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	paramIjn := -1
	paramStart := 0
	qParams := map[string]string{
		"tbm":     "isch",
		"asearch": "ichunk",
		"ijn":     strconv.Itoa(paramIjn),
		"start":   strconv.Itoa(paramStart),
		"q":       g.Config.Query,
		"hl":      "en",
		"async":   "_id:rg_s,_pms:s,_fmt:pc",
		"tbs":     filter,
	}

	switch g.Config.SafeSearch {
	case params.SafeSearchOff:
		qParams["safe"] = "images"
	case params.SafeSearchOn:
		qParams["safe"] = "active"
	default:
		return nil, 0, fmt.Errorf("--safe-search: invalid value %s", g.Config.SafeSearch)
	}

	hasMore := true
	items := make([]interface{}, 0)
	itemsURLCache := make(map[string]bool)
	pages := 0

	for hasMore {
		newCount := 0
		pages += 1
		paramIjn += 1
		paramStart = paramIjn * 100
		qParams["ijn"] = strconv.Itoa(paramIjn)
		qParams["start"] = strconv.Itoa(paramStart)
		page, err := util.DownloadWebpage("https://www.google.com/search", http.StatusOK, nil, qParams)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		r := bytes.NewReader(page)

		doc, err := goquery.NewDocumentFromReader(r)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		sel := doc.Find("div.rg_meta.notranslate")

		for i := range sel.Nodes {
			element := sel.Eq(i)
			gr := &GoogleResult{}
			err = json.Unmarshal([]byte(element.Text()), gr)

			if err != nil {
				return nil, 0, fmt.Errorf("cannot json.Unmarshal on match")
			}

			if !itemsURLCache[gr.URL] {
				items = append(items, gr)
				itemsURLCache[gr.URL] = true
				newCount++
			}
		}

		hasMore = newCount > 0
	}

	return items, pages, nil
}
