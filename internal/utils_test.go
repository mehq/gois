package internal

import (
	"os"
	"testing"
	"time"
)

type Test struct {
	in  string
	out int
}

var AtoiTests = []Test{
	{
		in:  "123",
		out: 123,
	},
	{
		in:  "0",
		out: 0,
	},
	{
		in:  "-1",
		out: -1,
	},
}

func TestAtoi(t *testing.T) {
	for _, test := range AtoiTests {
		got := Atoi(test.in)
		if test.out != got {
			t.Errorf("Atoi invalid output %d for input %s, expected %d", got, test.in, test.out)
		}
	}
}

func TestDownload(t *testing.T) {
	client := MakeHTTPClient()
	success, bw := Download(client, "https://example.com", "testfile.html")

	if !success || bw < 1 {
		t.Errorf("Error encountered while downloading file.")
	}

	if _, err := os.Stat("testfile.html"); os.IsNotExist(err) {
		t.Errorf("Error encountered while downloading file. File not written to disk.")
	}
}

func TestItoa(t *testing.T) {
	for _, test := range AtoiTests {
		got := Itoa(test.out)
		if test.in != got {
			t.Errorf("Atoi invalid output %s for input %d, expected %s", got, test.out, test.in)
		}
	}
}

func TestMakeHTTPClient(t *testing.T) {
	client := MakeHTTPClient()

	if client == nil {
		t.Errorf("MakeHTTPClient returned nil, expected a valid http.Client")
	}
}

func TestMakeRequest(t *testing.T) {
	req := MakeRequest("GET", "https://example.com", nil, nil)

	if req == nil {
		t.Errorf("MakeRequest returned nil, expected a valid http.Request")
	}
}

func TestMakeProgressBarOutput(t *testing.T) {
	out := MakeProgressBarOutput(&time.Time{}, 1024*1024, 1, 1, 2)
	expected := "Downloaded:    1 | Failed:    1 | Total:    2 |   0.000Mbps"

	if out != expected {
		t.Errorf("MakeProgressBarOutput invalid output %s, expected %s", out, expected)
	}
}
