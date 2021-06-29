package google

import (
	"testing"
	"time"

	"github.com/mzbaulhaque/gois/internal/util/request"
	"github.com/mzbaulhaque/gois/internal/util/testutil"
)

func TestGoogle_Scrape(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestGoogle_Scrape() panicked")
		}
	}()

	var optionTests = []*Config{
		{
			Query:    "cats",
			Explicit: true,
			GIF:      false,
			Gray:     false,
			Height:   0,
			Width:    0,
			Compact:  false,
		},
		{
			Query:    "cats",
			Explicit: false,
			GIF:      true,
			Gray:     false,
			Height:   0,
			Width:    0,
			Compact:  false,
		},
		{
			Query:    "cats",
			Explicit: false,
			GIF:      false,
			Gray:     true,
			Height:   0,
			Width:    0,
			Compact:  false,
		},
		{
			Query:    "cats",
			Explicit: false,
			GIF:      false,
			Gray:     false,
			Height:   1080,
			Width:    1920,
			Compact:  false,
		},
		{
			Query:    "cats",
			Explicit: true,
			GIF:      false,
			Gray:     false,
			Height:   1080,
			Width:    1920,
			Compact:  true,
		},
	}

	client := request.NewHTTPClient()

	for _, test := range optionTests {
		googleScraper := Google{
			Client: client,
			Config: test,
		}
		items, err := googleScraper.Scrape()

		testutil.CheckErr(t, err)

		if len(items) < 1 {
			t.Errorf("0 items scraped from google")
		}

		time.Sleep(500 * time.Millisecond)
	}
}
