package util

import (
	"testing"
)

func TestSearchRegex(t *testing.T) {
	var testCases = []TestCase{
		{
			In: map[string]string{
				"expr":   "",
				"target": "",
			},
			Out: "",
		},
		{
			In: map[string]string{
				"expr":   "[", // invalid pattern
				"target": "",
			},
			Out: "",
		},
		{
			In: map[string]string{
				"expr":   "w([^d]+)d",
				"target": "hello world",
			},
			Out: "orl",
		},
	}

	for _, test := range testCases {
		input := test.In.(map[string]string)
		expectedOutput := test.Out.(string)
		output, err := SearchRegex(input["expr"], input["target"], "")

		if expectedOutput == "" && err == nil {
			t.Errorf("should return an error")
		} else if expectedOutput != "" && expectedOutput != output {
			t.Errorf("expected '%s', got '%s'", expectedOutput, output)
		}
	}
}
