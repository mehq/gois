package main

import (
	"net/http"
	"net/url"
)

type Header struct {
	Name string
	Value string
}

type Options struct {
	query string
	explicit bool
	gif bool
	gray bool
	height int
	width int
}

func MakeUrl(baseUrl string, params *url.Values) string {
	base, err := url.Parse(baseUrl)

	if err != nil {
		panic(err)
	}

	base.RawQuery = params.Encode()

	return base.String()
}

func MakeRequest(method, url string, params *url.Values, headers []*Header) *http.Request {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}

	req.URL.RawQuery = params.Encode()

	for _, v := range headers {
		req.Header.Set(v.Name, v.Value)
	}

	return req
}
