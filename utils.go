package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Header is a http header item
type Header struct {
	Name  string
	Value string
}

// Options is collection of option used by scrapers to find specific
// data.
type Options struct {
	query  string
	safe   bool
	gif    bool
	gray   bool
	height int
	width  int
}

var defaultHeaders = []*Header{
	{
		Name:  "accept-language",
		Value: "en-US,en;q=0.9,pl;q=0.8,fr;q=0.7,bn;q=0.6",
	},
	{
		Name:  "user-agent",
		Value: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36",
	},
	{
		Name:  "upgrade-insecure-requests",
		Value: "1",
	},
}

// Atoi is equivalent to strconv.Atoi except it returns 0 on error.
func Atoi(rawValue string) int {
	value, err := strconv.Atoi(rawValue)

	if err != nil {
		return 0
	}

	return value
}

// Download is a utility function that downloads file into local machine at outFilePath.
func Download(client *http.Client, fileURL string, outFilePath string) (bool, int64) {
	resp, err := client.Get(fileURL)

	if err != nil {
		return false, 0
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	out, err := os.Create(outFilePath)

	if err != nil {
		return false, 0
	}

	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	bytesWritten, err := io.Copy(out, resp.Body)

	if err != nil {
		return false, 0
	}

	return true, bytesWritten
}

// Itoa is equivalent to strconv.Itoa
func Itoa(rawValue int) string {
	return strconv.Itoa(rawValue)
}

// MakeHTTPClient returns a pointer to http.Client with a default cookiejar.Jar
func MakeHTTPClient() *http.Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}

	return &client
}

// MakeProgressBarOutput returns formatted information shown by ProgressBar
func MakeProgressBarOutput(downloadStartedAt *time.Time, bw int64, sc int, fc int, rc int) string {
	speed := 0.0
	if !downloadStartedAt.IsZero() {
		elapsed := time.Since(*downloadStartedAt).Seconds()
		speed = (float64(bw) / (1024.0 * 1024.0)) / elapsed
	}
	return fmt.Sprintf("Downloaded: %4d | Failed: %4d | Total: %4d | %7.3fMbps", sc, fc, rc, speed)
}

// MakeRequest is a utility function to make a new http.Request instance with
// passed information (params, url, headers etc).
func MakeRequest(method, url string, params *url.Values, headers []*Header) *http.Request {
	if headers == nil {
		headers = defaultHeaders
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	for _, v := range headers {
		req.Header.Set(v.Name, v.Value)
	}

	return req
}
