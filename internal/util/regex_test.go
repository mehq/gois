// +build !regex

package util

import (
	"testing"
)

func TestSearchRegex(t *testing.T) {
	var testCases = []TestCase{
		{
			In:  map[string]string{
				"expr": "w([^d]+)d",
				"target": "hello world",
				"name": "",
				"fatal": "false",
			},
			Out: "orl",
		},
	}

	for _, test := range testCases {
		input := test.In.(map[string]string)
		expectedOutput := test.Out.(string)
		fatal := true

		if input["fatal"] == "false" {
			fatal = false
		}

		output, err := SearchRegex(input["expr"], input["target"], input["name"], fatal)

		if expectedOutput == "" && err == nil {
			t.Errorf("should raise an error")
		} else if expectedOutput != "" && expectedOutput != output {
			t.Errorf("expected %s, got %s", expectedOutput, output)
		}
	}
}
