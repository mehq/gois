// +build test

package util

import (
	"io/ioutil"
)

func init() {
	content, err := ioutil.ReadFile("testdata/bing_images_async")

	if err != nil {
		panic(err)
	}

	ResponseBingImageAsync = content

	ResponseBingImagesSearch, err = ioutil.ReadFile("testdata/bing_images_search")

	if err != nil {
		panic(err)
	}

	ResponseBingSettings, err = ioutil.ReadFile("testdata/bing_settings")

	if err != nil {
		panic(err)
	}

	ResponseFlickrSearch, err = ioutil.ReadFile("testdata/flickr_search")

	if err != nil {
		panic(err)
	}

	ResponseFlickrSearchAPI, err = ioutil.ReadFile("testdata/flickr_search_api")

	if err != nil {
		panic(err)
	}

	ResponseGoogleSearch, err = ioutil.ReadFile("testdata/google_search")

	if err != nil {
		panic(err)
	}

	RegisterMockHTTPClient()
}
