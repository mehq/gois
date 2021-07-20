package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/mzbaulhaque/gois/internal/util"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// BingConfig is a set of options used by Bing.
type BingConfig struct {
	AspectRatio string
	Compact  bool
	Date string
	ImageColor string
	ImageSize string
	ImageType string
	License string
	People string
	Query    string
	SafeSearch string
}

// BingScraper is used to scrape data from bing search engine.
type BingScraper struct {
	Config *BingConfig
}

type BingResult struct {
	Height       int    `json:"-"`
	Width        int    `json:"-"`
	ReferenceURL string `json:"purl"`
	ThumbnailURL string `json:"turl"`
	Title        string `json:"t"`
	URL          string `json:"murl"`
}

func isValidCustomSize(customSize string) bool {
	match, _ := regexp.MatchString("[0-9]+_[0-9]+", customSize)

	return match
}

func (b BingScraper) makeFilterString() (string, error) {
	filters := make([]string, 0)

	if b.Config.AspectRatio == "square" {
		filters = append(filters, "filterui:aspect-square")
	} else if b.Config.AspectRatio == "wide" {
		filters = append(filters, "filterui:aspect-wide")
	} else if b.Config.AspectRatio == "tall" {
		filters = append(filters, "filterui:aspect-tall")
	} else if b.Config.AspectRatio != "" {
		return "", fmt.Errorf("--aspect-ratio: invalid value %s", b.Config.AspectRatio)
	}

	if b.Config.Date == "past-day" {
		filters = append(filters, "filterui:age-lt1440")
	} else if b.Config.Date == "past-week" {
		filters = append(filters, "filterui:age-lt10080")
	} else if b.Config.Date == "past-month" {
		filters = append(filters, "filterui:age-lt43200")
	} else if b.Config.Date == "past-year" {
		filters = append(filters, "filterui:age-lt525600")
	} else if b.Config.Date != "" {
		return "", fmt.Errorf("--date: invalid value %s", b.Config.Date)
	}

	if b.Config.ImageColor == "full-color" {
		filters = append(filters, "filterui:color2-color")
	} else if b.Config.ImageColor == "full-color" {
		filters = append(filters, "filterui:color2-color")
	} else if b.Config.ImageColor == "black-white" {
		filters = append(filters, "filterui:color2-bw")
	} else if b.Config.ImageColor == "red" {
		filters = append(filters, "filterui:color2-FGcls_RED")
	} else if b.Config.ImageColor == "orange" {
		filters = append(filters, "filterui:color2-FGcls_ORANGE")
	} else if b.Config.ImageColor == "yellow" {
		filters = append(filters, "filterui:color2-FGcls_YELLOW")
	} else if b.Config.ImageColor == "green" {
		filters = append(filters, "filterui:color2-FGcls_GREEN")
	} else if b.Config.ImageColor == "teal" {
		filters = append(filters, "filterui:color2-FGcls_TEAL")
	} else if b.Config.ImageColor == "blue" {
		filters = append(filters, "filterui:color2-FGcls_BLUE")
	} else if b.Config.ImageColor == "purple" {
		filters = append(filters, "filterui:color2-FGcls_PURPLE")
	} else if b.Config.ImageColor == "pink" {
		filters = append(filters, "filterui:color2-FGcls_PINK")
	} else if b.Config.ImageColor == "brown" {
		filters = append(filters, "filterui:color2-FGcls_BROWN")
	} else if b.Config.ImageColor == "black" {
		filters = append(filters, "filterui:color2-FGcls_BLACK")
	} else if b.Config.ImageColor == "gray" {
		filters = append(filters, "filterui:color2-FGcls_GRAY")
	} else if b.Config.ImageColor == "white" {
		filters = append(filters, "filterui:color2-FGcls_WHITE")
	} else if b.Config.ImageColor != "" {
		return "", fmt.Errorf("--image-color: invalid value %s", b.Config.ImageColor)
	}

	if b.Config.ImageSize == "small" {
		filters = append(filters, "filterui:imagesize-small")
	} else if b.Config.ImageSize == "medium" {
		filters = append(filters, "filterui:imagesize-medium")
	} else if b.Config.ImageSize == "large" {
		filters = append(filters, "filterui:imagesize-large")
	} else if b.Config.ImageSize == "extra-large" {
		filters = append(filters, "filterui:wallpaper")
	} else if b.Config.ImageSize != "" {
		if isValidCustomSize(b.Config.ImageSize) {
			filters = append(filters, fmt.Sprintf("filterui:imagesize-custom_%s", b.Config.ImageSize))
		} else {
			return "", fmt.Errorf("--image-size: invalid value %s", b.Config.ImageSize)
		}
	}

	if b.Config.ImageType == "photo" {
		filters = append(filters, "filterui:photo-photo")
	} else if b.Config.ImageType == "clip-art" {
		filters = append(filters, "filterui:photo-clipart")
	} else if b.Config.ImageType == "line-drawing" {
		filters = append(filters, "filterui:photo-linedrawing")
	} else if b.Config.ImageType == "animated" {
		filters = append(filters, "filterui:photo-animatedgif")
	} else if b.Config.ImageType == "transparent" {
		filters = append(filters, "filterui:photo-transparent")
	} else if b.Config.ImageType != "" {
		return "", fmt.Errorf("--image-type: invalid value %s", b.Config.ImageType)
	}

	if b.Config.License == "creative-commons" {
		filters = append(filters, "filterui:licenseType-Any")
	} else if b.Config.License == "public-domain" {
		filters = append(filters, "filterui:license-L1")
	} else if b.Config.License == "free-share-use" {
		filters = append(filters, "filterui:license-license-L2_L3_L4_L5_L6_L7")
	} else if b.Config.License == "free-share-use-commercially" {
		filters = append(filters, "filterui:license-L2_L3_L4")
	} else if b.Config.License == "free-modify-share-use" {
		filters = append(filters, "filterui:license-L2_L3_L5_L6")
	} else if b.Config.License == "free-modify-share-use-commercially" {
		filters = append(filters, "filterui:license-L2_L3")
	} else if b.Config.License != "" {
		return "", fmt.Errorf("--license: invalid value %s", b.Config.License)
	}

	if b.Config.People == "face" {
		filters = append(filters, "filterui:face-face")
	} else if b.Config.People == "head-shoulder" {
		filters = append(filters, "filterui:face-portrait")
	} else if b.Config.People != "" {
		return "", fmt.Errorf("--people: invalid value %s", b.Config.People)
	}

	return strings.Join(filters, "+"), nil
}

// Scrape is the entrypoint.
func (b BingScraper) Scrape() ([]interface{}, error) {
	filter, err := b.makeFilterString()

	if err != nil {
		return nil, err
	}

	var safeSearchOption string

	if b.Config.SafeSearch == "off" {
		safeSearchOption = "OFF"
	} else if b.Config.SafeSearch == "on" {
		safeSearchOption = "STRICT"
	} else if b.Config.SafeSearch == "moderate" {
		safeSearchOption = "DEMOTE"
	} else if b.Config.SafeSearch != "" {
		return nil, fmt.Errorf("--safe-search: invalid value %s", b.Config.SafeSearch)
	}

	paramFirst := -150
	headers := map[string]string{
		"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"cookie": fmt.Sprintf("SRCHHPGUSR=ADLT=%s", safeSearchOption),
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"upgrade-insecure-requests": "1",
	}
	params := map[string]string{
		"count": "150",
		"first": strconv.Itoa(paramFirst),
		"q": b.Config.Query,
		"qft": filter,
		"relp": "150",
	}

	hasMore := true
	items := make([]interface{}, 0)
	itemsURLCache := make(map[string]bool)

	for hasMore {
		newCount := 0
		paramFirst += 150
		params["first"] = strconv.Itoa(paramFirst)
		page, err := util.DownloadWebpage("https://www.bing.com/images/async", http.StatusOK, &headers, &params)

		if err != nil {
			return nil, err
		}

		matches, err := util.SearchRegexMultiple("iusc\"[\\s\\n]+style=\"[^\"]+\"[\\s\\n]+(?:gif=\"[^\"]+\"[\\s\\n]+)?m=\"([^\"]+)\"", string(page), "links", false)

		if err != nil {
			return nil, err
		}

		widths, err := util.SearchRegexMultiple("nowrap\">\\s*([0-9]+)\\s*x\\s*[0-9]+\\s*&#183;\\s*[a-zA-Z]+</span>", string(page), "widths", false)

		if err != nil {
			return nil, err
		}

		heights, err := util.SearchRegexMultiple("nowrap\">\\s*[0-9]+\\s*x\\s*([0-9]+)\\s*&#183;\\s*[a-zA-Z]+</span>", string(page), "heights", false)

		if err != nil {
			return nil, err
		}

		if len(matches) != len(widths) && len(widths) != len(heights) {
			return nil, fmt.Errorf("matche count != width count != height count")
		}

		for i, match := range matches {
			br := &BingResult{}
			err = json.Unmarshal([]byte(strings.ReplaceAll(match, "&quot;", "\"")), br)

			if err != nil {
				return nil, fmt.Errorf("cannot json.Unmarshal on match %s", strings.ReplaceAll(match, "&quot;", "\""))
			}

			br.Height, err = strconv.Atoi(heights[i])

			if err != nil {
				return nil, fmt.Errorf("cannot convert height (%s) to int", heights[i])
			}

			br.Width, err = strconv.Atoi(widths[i])

			if err != nil {
				return nil, fmt.Errorf("cannot convert width (%s) to int", widths[i])
			}

			if !itemsURLCache[br.URL] {
				items = append(items, br)
				itemsURLCache[br.URL] = true
				newCount++
			}
		}

		hasMore = newCount > 0
	}

	return items, nil
}
