package out

import (
	"testing"

	"github.com/mzbaulhaque/gois/internal/util/testutil"
)

func TestPrintImageInfo(t *testing.T) {
	type imageInfoItem struct {
		Height       int
		Width        int
		ReferenceURL string
		ThumbnailURL string
		Title        string
		URL          string
	}

	var item = imageInfoItem{
		Height:       1080,
		Width:        1920,
		ReferenceURL: "https://example.com",
		ThumbnailURL: "https://example.com",
		Title:        "https://example.com",
		URL:          "https://example.com",
	}

	type inputCase struct {
		items   []interface{}
		compact bool
	}

	var testCases = []testutil.TestCase{
		{
			In: inputCase{
				items: []interface{}{
					item,
					item,
				},
				compact: true,
			},
		},
		{
			In: inputCase{
				items: []interface{}{
					item,
					item,
				},
				compact: false,
			},
		},
		{
			In: inputCase{
				items: []interface{}{
					item,
				},
				compact: true,
			},
		},
		{
			In: inputCase{
				items: []interface{}{
					item,
				},
				compact: false,
			},
		},
	}

	for _, test := range testCases {
		err := PrintImageInfo(test.In.(inputCase).items, test.In.(inputCase).compact)
		testutil.CheckErr(t, err)
	}
}
