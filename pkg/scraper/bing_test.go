package scraper

import (
	"github.com/mzbaulhaque/gomage/internal"
	"testing"
	"time"
)

func TestBing_Scrape(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestBing_Scrape() panicked")
		}
	}()

	var optionTests = []internal.Options{
		{
			Query:    "cats",
			Safe:     true,
			Gif:      false,
			Gray:     false,
			Height:   0,
			Width:    0,
			TestMode: true,
		},
		{
			Query:    "cats",
			Safe:     false,
			Gif:      true,
			Gray:     false,
			Height:   0,
			Width:    0,
			TestMode: true,
		},
		{
			Query:    "cats",
			Safe:     false,
			Gif:      false,
			Gray:     true,
			Height:   0,
			Width:    0,
			TestMode: true,
		},
		{
			Query:    "cats",
			Safe:     false,
			Gif:      false,
			Gray:     false,
			Height:   1080,
			Width:    1920,
			TestMode: true,
		},
		{
			Query:    "cats",
			Safe:     false,
			Gif:      false,
			Gray:     false,
			Height:   0,
			Width:    0,
			TestMode: false,
		},
	}

	client := internal.MakeHTTPClient()

	for _, test := range optionTests {
		bing := Bing{
			Client: client,
			Opts:   &test,
		}
		items := bing.Scrape()

		if len(items) < 1 {
			t.Errorf("0 items scraped from bing")
		}

		time.Sleep(500 * time.Millisecond)
	}
}
