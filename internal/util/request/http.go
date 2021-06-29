package request

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type addCommonHeaderTransport struct {
	T http.RoundTripper
}

func (adt *addCommonHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("accept-language", "en-US,en;q=0.9,pl;q=0.8,fr;q=0.7,bn;q=0.6")
	req.Header.Add("user-agent", `Mozilla/5.0 (X11; Linux x86_64) `+
		`AppleWebKit/537.36 (KHTML, like Gecko)`+
		`Chrome/79.0.3945.117 Safari/537.36`)
	req.Header.Add("upgrade-insecure-requests", "1")

	res, err := adt.T.RoundTrip(req)

	if err != nil {
		return nil, fmt.Errorf("failed to execute http transaction")
	}

	return res, nil
}

// NewHTTPClient returns a pointer to http.Client with a default cookiejar.Jar.
func NewHTTPClient() *http.Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := http.Client{
		Jar:       jar,
		Timeout:   5 * time.Second,
		Transport: &addCommonHeaderTransport{http.DefaultTransport},
	}

	return &client
}

// NewRequest is a utility function to make a new http.Request instance with
// passed information (params, url).
func NewRequest(method, url string, params *url.Values) *http.Request {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	return req
}
