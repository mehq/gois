package main

import (
	"testing"
)

var optionTests = []Options{
	{
		query:  "cats",
		safe:   true,
		gif:    false,
		gray:   false,
		height: 0,
		width:  0,
	},
	{
		query:  "cats",
		safe:   false,
		gif:    false,
		gray:   false,
		height: 0,
		width:  0,
	},
	{
		query:  "cats",
		safe:   true,
		gif:    true,
		gray:   false,
		height: 0,
		width:  0,
	},
	{
		query:  "cats",
		safe:   true,
		gif:    false,
		gray:   true,
		height: 0,
		width:  0,
	},
	{
		query:  "cats",
		safe:   true,
		gif:    false,
		gray:   false,
		height: 1080,
		width:  0,
	},
	{
		query:  "cats",
		safe:   true,
		gif:    false,
		gray:   false,
		height: 0,
		width:  1920,
	},
}

func TestBingScrape(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestBingScrape() panicked")
		}
	}()

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
	}
}
