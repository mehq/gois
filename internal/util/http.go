package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var defaultHeaders = map[string]string{
	"accept":          "*/*",
	"accept-language": "en-US,en;q=0.9,bn;q=0.8,zh-CN;q=0.7,zh;q=0.6",
	"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/91.0.4472.164 Safari/537.36",
	"upgrade-insecure-requests": "1",
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	jar    *cookiejar.Jar
	Client HTTPClient
)

func init() {
	jar, _ = cookiejar.New(nil)
	Client = &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}
}

func SetCookie(name, value, path, domain, forUrl string) {
	u, _ := url.Parse(forUrl)
	jar.SetCookies(u, []*http.Cookie{
		{
			Name:   name,
			Value:  value,
			Path:   path,
			Domain: domain,
		},
	})
}

// DownloadWebpage downloads a webpage and returns content as byte array if successful. An error is returned
// otherwise.
func DownloadWebpage(
	pageUrl string,
	expectedStatus int,
	headers map[string]string,
	params map[string]string,
) ([]byte, error) {
	req, _ := http.NewRequest("GET", pageUrl, nil)

	for k, v := range defaultHeaders {
		req.Header.Add(k, v)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()

	for k, v := range params {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	res, err := Client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != expectedStatus {
		return nil, fmt.Errorf("got status %d, expected %d", res.StatusCode, expectedStatus)
	}

	content, _ := ioutil.ReadAll(res.Body)

	return content, nil
}
