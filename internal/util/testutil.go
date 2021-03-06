package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

// TestCase represents a generic test case.
type TestCase struct {
	In  interface{}
	Out interface{}
}

// CheckErr can be used to report generic test errors.
func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

// CheckCmdOutput can be used to match output of a command with a target Regexp.
func CheckCmdOutput(t *testing.T, output []byte, matchWith *regexp.Regexp) {
	if !matchWith.Match(output) {
		t.Errorf("command output not matching")
	}
}

var (
	ResponseBingImageAsync   []byte
	ResponseBingImagesSearch []byte
	ResponseFlickrSearch     []byte
	ResponseFlickrSearchAPI  []byte
	ResponseGoogleSearch     []byte
	ResponseYahooSearch      []byte
	ResponseYandexSearch     []byte
	ResponseYandexSearchRaw  []byte
)

// MockClient is the mock client.
type MockClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

// RegisterMockHTTPClient replaces real http client with a mock client.
func RegisterMockHTTPClient() {
	Client = &MockClient{
		MockDo: func(req *http.Request) (*http.Response, error) {
			url := req.URL.String()

			if strings.Contains(url, "test-do-error") {
				return nil, fmt.Errorf("dummy error")
			}

			var body io.ReadCloser

			switch {
			case strings.HasPrefix(url, "https://www.bing.com/images/async"):
				body = ioutil.NopCloser(bytes.NewReader(ResponseBingImageAsync))
			case strings.HasPrefix(url, "https://www.bing.com/images/search"):
				body = ioutil.NopCloser(bytes.NewReader(ResponseBingImagesSearch))
			case strings.HasPrefix(url, "https://www.flickr.com/search"):
				body = ioutil.NopCloser(bytes.NewReader(ResponseFlickrSearch))
			case strings.HasPrefix(url, "https://api.flickr.com/services/rest"):
				body = ioutil.NopCloser(bytes.NewReader(ResponseFlickrSearchAPI))
			case strings.HasPrefix(url, "https://www.google.com/search"):
				body = ioutil.NopCloser(bytes.NewReader(ResponseGoogleSearch))
			case strings.HasPrefix(url, "https://images.search.yahoo.com/search/images"):
				body = ioutil.NopCloser(bytes.NewReader(ResponseYahooSearch))
			case strings.HasPrefix(url, "https://yandex.com/images/search"):
				if _, exists := req.URL.Query()["format"]; exists {
					body = ioutil.NopCloser(bytes.NewReader(ResponseYandexSearch))
				} else {
					body = ioutil.NopCloser(bytes.NewReader(ResponseYandexSearchRaw))
				}
			default:
				body = ioutil.NopCloser(bytes.NewReader([]byte("dummy response")))
			}

			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		},
	}
}
