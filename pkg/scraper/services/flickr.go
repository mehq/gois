package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
)

// FlickrConfig is a set of options used by FlickrScraper to perform/filter/format search results.
type FlickrConfig struct {
	Compact     bool
	ImageColor  string
	ImageSize   string
	ImageType   string
	Orientation string
	Query       string
	SafeSearch  string
}

// FlickrScraper represents scraper for flickr image search.
type FlickrScraper struct {
	Config *FlickrConfig
}

// FlickrResult is a set of attributes that defines an image result.
type FlickrResult struct {
	Height       int    `json:"height_o"`
	Width        int    `json:"width_o"`
	Owner        string `json:"owner"`
	ID           string `json:"id"`
	ReferenceURL string `json:"-"`
	ThumbnailURL string `json:"url_t"`
	Title        string `json:"title"`
	URL          string `json:"url_o"`
}

type flickrResponsePhotos struct {
	Page  int `json:"page"`
	Pages int `json:"pages"`
	Photo []FlickrResult
}

type flickrResponse struct {
	Photos flickrResponsePhotos `json:"photos"`
	Stat   string               `json:"stat"`
}

func (f FlickrScraper) getFilters() (map[string]string, error) {
	filters := map[string]string{}

	switch f.Config.ImageColor {
	case "", params.ParamAll:
	case params.ColorRed:
		filters["color_codes"] = "0"
	case params.ColorDarkOrange:
		filters["color_codes"] = "1"
	case params.ColorOrange:
		filters["color_codes"] = "2"
	case params.ColorPalePink:
		filters["color_codes"] = "b"
	case params.ColorLemonYellow:
		filters["color_codes"] = "4"
	case params.ColorSchoolBusYellow:
		filters["color_codes"] = "3"
	case params.ColorGreen:
		filters["color_codes"] = "5"
	case params.ColorDarkLimeGreen:
		filters["color_codes"] = "6"
	case params.ColorCyan:
		filters["color_codes"] = "7"
	case params.ColorBlue:
		filters["color_codes"] = "8"
	case params.ColorViolet:
		filters["color_codes"] = "9"
	case params.ColorPink:
		filters["color_codes"] = "a"
	case params.ColorWhite:
		filters["color_codes"] = "c"
	case params.ColorGray:
		filters["color_codes"] = "d"
	case params.ColorBlack:
		filters["color_codes"] = "e"
	default:
		return nil, fmt.Errorf("--image-color: invalid value %s", f.Config.ImageColor)
	}

	switch f.Config.ImageSize {
	case "", params.ParamAll:
	case params.ImageSizeMedium:
		filters["dimension_search_mode"] = "min"
		filters["height"] = "640"
		filters["width"] = "640"
	case params.ImageSizeLarge:
		filters["dimension_search_mode"] = "min"
		filters["height"] = "1024"
		filters["width"] = "1024"
	default:
		return nil, fmt.Errorf("--image-size: invalid value %s", f.Config.ImageSize)
	}

	switch f.Config.ImageType {
	case "", params.ParamAll:
	case params.ColorBlackAndWhite:
		filters["styles"] = "blackandwhite"
	case params.ImageTypeShallowDepthOfField:
		filters["styles"] = "depthoffield"
	case params.ImageTypeMinimal:
		filters["styles"] = "minimalism"
	case params.ImageTypePatterns:
		filters["styles"] = "pattern"
	default:
		return nil, fmt.Errorf("--image-type: invalid value %s", f.Config.ImageType)
	}

	switch f.Config.Orientation {
	case "", params.ParamAll:
	case params.OrientationLandscape:
		filters["orientation"] = "landscape"
	case params.OrientationPortrait:
		filters["orientation"] = "portrait"
	case params.AspectRatioSquare:
		filters["orientation"] = "square"
	case params.AspectRatioPanoramic:
		filters["orientation"] = "panorama"
	default:
		return nil, fmt.Errorf("--orientation: invalid value %s", f.Config.Orientation)
	}

	return filters, nil
}

func (f FlickrScraper) getAPIKey() (string, error) {
	page, err := util.DownloadWebpage("https://www.flickr.com/search", http.StatusOK, nil, map[string]string{
		"text": f.Config.Query,
	})

	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	key, err := util.SearchRegex(
		"root.YUI_config.flickr.api.site_key\\s*=\\s*\"([a-z0-9]+)\"",
		string(page),
		"apiKey",
	)

	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return key, nil
}

// Scrape parses and returns the results from flickr image search if successful. An error is returned otherwise.
func (f FlickrScraper) Scrape() ([]interface{}, int, error) {
	filters, err := f.getFilters()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	if len(filters) > 0 {
		filters["advanced"] = "1"
	}

	apiKey, err := f.getAPIKey()

	if err != nil {
		return nil, 0, fmt.Errorf("%v", err)
	}

	paramPage := 0
	qParams := map[string]string{
		"sort":         "relevance",
		"content_type": "7",
		"extras": "count_comments,count_faves,count_views,date_taken,date_upload,description,icon_urls_deep," +
			"license,path_alias,perm_print,realname,url_t,url_o,visibility,visibility_source,o_dims",
		"per_page":       "500",
		"page":           strconv.Itoa(paramPage),
		"text":           f.Config.Query,
		"media":          "photos",
		"method":         "flickr.photos.search",
		"api_key":        apiKey,
		"format":         "json",
		"nojsoncallback": "1",
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
		paramPage += 1
		qParams["page"] = strconv.Itoa(paramPage)
		page, err := util.DownloadWebpage(
			"https://api.flickr.com/services/rest",
			http.StatusOK,
			nil,
			qParams,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("%v", err)
		}

		res := &flickrResponse{}

		err = json.Unmarshal(page, res)

		if err != nil {
			return nil, 0, fmt.Errorf("cannot json.Unmarshal on response")
		}

		if res.Stat != "ok" {
			return nil, 0, fmt.Errorf("response not ok, got %s", res.Stat)
		}

		for _, v := range res.Photos.Photo {
			if !itemsURLCache[v.URL] {
				v.ReferenceURL = fmt.Sprintf("https://www.flickr.com/photos/%s/%s", v.Owner, v.ID)
				items = append(items, v)
				itemsURLCache[v.URL] = true
				newCount++
			}
		}

		hasMore = res.Photos.Page < res.Photos.Pages
	}

	return items, pages, nil
}
