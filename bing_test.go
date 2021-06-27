package main

import (
	"testing"
	"time"
)

func TestBing_Scrape(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestBing_Scrape() panicked")
		}
	}()

	var optionTests = []Options{
		{
			query:    "cats",
			safe:     true,
			gif:      false,
			gray:     false,
			height:   0,
			width:    0,
			testMode: true,
		},
		{
			query:    "cats",
			safe:     false,
			gif:      true,
			gray:     false,
			height:   0,
			width:    0,
			testMode: true,
		},
		{
			query:    "cats",
			safe:     false,
			gif:      false,
			gray:     true,
			height:   0,
			width:    0,
			testMode: true,
		},
		{
			query:    "cats",
			safe:     false,
			gif:      false,
			gray:     false,
			height:   1080,
			width:    1920,
			testMode: true,
		},
	}

	client := MakeHTTPClient()

	for _, test := range optionTests {
		bing := Bing{
			client: client,
			opts:   &test,
		}
		items := bing.Scrape()

		if len(items) < 1 {
			t.Errorf("0 items scraped from bing")
		}

		time.Sleep(500 * time.Millisecond)
	}
}
