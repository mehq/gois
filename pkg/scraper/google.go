package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/mzbaulhaque/gois/internal/util"
	"net/http"
	"strconv"
	"strings"
)

// GoogleConfig is a set of options used by Google.
type GoogleConfig struct {
	AspectRatio string
	Compact  bool
	FileType string
	ImageColor string
	ImageSize string
	ImageType string
	License string
	Query    string
	Region string
	SafeSearch string
}

// GoogleScraper is used to scrape data from google search engine.
type GoogleScraper struct {
	Config *GoogleConfig
}

// GoogleResult is
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

	if g.Config.AspectRatio == "tall" {
		filters = append(filters, "iar:t")
	} else if g.Config.AspectRatio == "square" {
		filters = append(filters, "iar:s")
	} else if g.Config.AspectRatio == "wide" {
		filters = append(filters, "iar:w")
	} else if g.Config.AspectRatio == "panoramic" {
		filters = append(filters, "iar:xw")
	} else if g.Config.AspectRatio != "" {
		return "", fmt.Errorf("--aspect-ratio: invalid value %s", g.Config.AspectRatio)
	}

	if g.Config.FileType == "jpg" {
		filters = append(filters, "ift:jpg")
	} else if g.Config.FileType == "gif" {
		filters = append(filters, "ift:gif")
	} else if g.Config.FileType == "png" {
		filters = append(filters, "ift:png")
	} else if g.Config.FileType == "bmp" {
		filters = append(filters, "ift:bmp")
	} else if g.Config.FileType == "svg" {
		filters = append(filters, "ift:svg")
	} else if g.Config.FileType == "webp" {
		filters = append(filters, "ift:webp")
	} else if g.Config.FileType == "ico" {
		filters = append(filters, "ift:ico")
	} else if g.Config.FileType == "raw" {
		filters = append(filters, "ift:raw")
	} else if g.Config.FileType != "" {
		return "", fmt.Errorf("--file-type: invalid value %s", g.Config.FileType)
	}

	if g.Config.ImageColor == "full-color" {
		filters = append(filters, "ic:color")
	} else if g.Config.ImageColor == "black-white" {
		filters = append(filters, "ic:gray")
	} else if g.Config.ImageColor == "transparent" {
		filters = append(filters, "ic:trans")
	} else if g.Config.ImageColor == "red" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:red")
	} else if g.Config.ImageColor == "orange" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:orange")
	} else if g.Config.ImageColor == "yellow" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:yellow")
	} else if g.Config.ImageColor == "green" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:green")
	} else if g.Config.ImageColor == "teal" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:teal")
	} else if g.Config.ImageColor == "blue" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:blue")
	} else if g.Config.ImageColor == "purple" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:purple")
	} else if g.Config.ImageColor == "pink" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:pink")
	} else if g.Config.ImageColor == "white" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:white")
	} else if g.Config.ImageColor == "gray" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:gray")
	} else if g.Config.ImageColor == "black" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:black")
	} else if g.Config.ImageColor == "brown" {
		filters = append(filters, "ic:specific")
		filters = append(filters, "isc:brown")
	} else if g.Config.ImageColor != "" {
		return "", fmt.Errorf("--image-color: invalid value %s", g.Config.ImageColor)
	}

	if g.Config.ImageSize == "large" {
		filters = append(filters, "isz:l")
	} else if g.Config.ImageSize == "medium" {
		filters = append(filters, "isz:m")
	} else if g.Config.ImageSize == "icon" {
		filters = append(filters, "isz:i")
	} else if g.Config.ImageSize == "qsvga" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:qsvga")
	} else if g.Config.ImageSize == "vga" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:vga")
	} else if g.Config.ImageSize == "svga" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:svga")
	} else if g.Config.ImageSize == "xga" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:xga")
	} else if g.Config.ImageSize == "2mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:2mp")
	} else if g.Config.ImageSize == "4mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:4mp")
	} else if g.Config.ImageSize == "6mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:6mp")
	} else if g.Config.ImageSize == "8mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:8mp")
	} else if g.Config.ImageSize == "10mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:10mp")
	} else if g.Config.ImageSize == "12mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:12mp")
	} else if g.Config.ImageSize == "15mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:15mp")
	} else if g.Config.ImageSize == "20mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:20mp")
	} else if g.Config.ImageSize == "40mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:40mp")
	} else if g.Config.ImageSize == "70mp" {
		filters = append(filters, "isz:lt")
		filters = append(filters, "islt:70mp")
	} else if g.Config.ImageSize != "" {
		return "", fmt.Errorf("--image-size: invalid value %s", g.Config.ImageSize)
	}

	if g.Config.ImageType == "face" {
		filters = append(filters, "itp:face")
	} else if g.Config.ImageType == "photo" {
		filters = append(filters, "itp:photo")
	} else if g.Config.ImageType == "clip-art" {
		filters = append(filters, "itp:clipart")
	} else if g.Config.ImageType == "line-drawing" {
		filters = append(filters, "itp:lineart")
	} else if g.Config.ImageType == "animated" {
		filters = append(filters, "itp:animated")
	} else if g.Config.ImageType != "" {
		return "", fmt.Errorf("--image-type: invalid value %s", g.Config.ImageType)
	}

	if g.Config.License == "creative-commons" {
		filters = append(filters, "sur:cl")
	} else if g.Config.License == "commercial" {
		filters = append(filters, "sur:ol")
	} else if g.Config.License != "" {
		return "", fmt.Errorf("--license: invalid value %s", g.Config.License)
	}

	return strings.Join(filters, ","), nil
}

// Scrape is the entrypoint.
func (g GoogleScraper) Scrape() ([]interface{}, int, error) {
	filter, err := g.makeFilterString()

	if err != nil {
		return nil, 0, err
	}

	paramIjn := -1
	paramStart := 0
	headers := map[string]string{
		"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"upgrade-insecure-requests": "1",
	}
	params := map[string]string{
		"tbm": "isch",
		"asearch": "ichunk",
		"ijn": strconv.Itoa(paramIjn),
		"start": strconv.Itoa(paramStart),
		"q": g.Config.Query,
		"hl": "en",
		"async": "_id:rg_s,_pms:s,_fmt:pc",
		"tbs": filter,
	}

	if g.Config.SafeSearch == "on" {
		params["safe"] = "active"
	} else if g.Config.SafeSearch == "off" {
		params["safe"] = "images"
	} else if g.Config.SafeSearch != "" {
		return nil, 0, fmt.Errorf("--safe-search: invalid value %s", g.Config.SafeSearch)
	}

	if g.Config.Region != "" {
		params["cr"] = fmt.Sprintf("country%s", g.Config.Region)
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
		params["ijn"] = strconv.Itoa(paramIjn)
		params["start"] = strconv.Itoa(paramStart)
		page, err := util.DownloadWebpage("https://www.google.com/search", http.StatusOK, &headers, &params)

		if err != nil {
			return nil, 0, err
		}

		matches, err := util.SearchRegexMultiple("notranslate\"[^>]*>([^<]+)", string(page), "links", false)

		if err != nil {
			return nil, 0, err
		}

		for _, match := range matches {
			gr := &GoogleResult{}
			err = json.Unmarshal([]byte(match), gr)

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

	return items, 0, nil
}
