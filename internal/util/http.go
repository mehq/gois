package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

var defaultHeaders = map[string]string{
	"accept":                    "*/*",
	"accept-language":           "en-US,en;q=0.9,bn;q=0.8,zh-CN;q=0.7,zh;q=0.6",
	"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36",
	"upgrade-insecure-requests": "1",
}

var httpClient *http.Client = nil

func getHttpClient() *http.Client {
	if httpClient == nil {
		jar, _ := cookiejar.New(nil)

		httpClient = &http.Client{
			Jar: jar,
		}
	}

	return httpClient
}

// DownloadWebpage downloads a page and returns its contents as byte array
func DownloadWebpage(pageUrl string, expectedStatus int, headers map[string]string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", pageUrl, nil)

	if err != nil {
		return nil, err
	}

	for k, v := range defaultHeaders {
		req.Header.Add(k, v)
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	if params != nil {
		q := req.URL.Query()

		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	client := getHttpClient()

	defer func() {
		client.CloseIdleConnections()
	}()

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := res.Body.Close()

		if err != nil {
			panic(err)
		}
	}()

	if res.StatusCode != expectedStatus {
		return nil, fmt.Errorf("got status %d, expected %d", res.StatusCode, expectedStatus)
	}

	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return content, nil
}
