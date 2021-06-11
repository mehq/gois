package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	client  *http.Client
	headers []*Header
)

func makeFilterString(opts *Options) string {
	filters := make([]string, 0)

	if opts.gif {
		filters = append(filters, "filterui:photo-animatedgif")
	}

	if opts.gray {
		filters = append(filters, "filterui:color2-bw")
	}

	return strings.Join(filters, "+")
}

func defang(opts *Options) {
	params := &url.Values{}
	params.Add("q", opts.query)

	res, err := client.Do(MakeRequest("GET", "https://www.bing.com/images/search", params, headers))

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		panic(err)
	}

	guid, exists := doc.Find("input#GUID").Attr("value")

	if !exists {
		panic("guid not found")
	}

	ru, exists := doc.Find("input#ru").Attr("value")

	if !exists {
		panic("ru not found")
	}

	params.Del("q")
	params.Add("pref_sbmt", "1")
	params.Add("adlt_set", "off")
	params.Add("adlt_confirm", "1")
	params.Add("GUID", guid)
	params.Add("is_child", "0")
	params.Add("ru", ru)

	_, err = client.Do(MakeRequest("GET", "https://www.bing.com/settings.aspx", params, headers))

	if err != nil {
		panic(err)
	}
}

func ScrapeBing(opts *Options) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	headers := make([]*Header, 3)
	headers[0] = &Header{
		Name:  "accept-language",
		Value: "en-US,en;q=0.9,pl;q=0.8,fr;q=0.7,bn;q=0.6",
	}
	headers[1] = &Header{
		Name:  "user-agent",
		Value: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36",
	}
	headers[2] = &Header{
		Name:  "upgrade-insecure-requests",
		Value: "1",
	}

	client = &http.Client{
		Jar: jar,
	}

	if opts.explicit {
		defang(opts)
	}

	params := &url.Values{}
	params.Add("q", opts.query)
	params.Add("first", "0")
	params.Add("count", "150")
	params.Add("relp", "150")
	params.Add("qft", makeFilterString(opts))

	hasMore := true

	for hasMore {
		newCount := 0
		res, err := client.Do(MakeRequest("GET", "https://www.bing.com/images/async", params, headers))

		if err != nil {
			panic(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			panic(err)
		}

		doc.Find("a.iusc").Each(func(_ int, element *goquery.Selection) {
			rawInfo, exists := element.Attr("m")

			if !exists {
				panic("Does not exist")
			}

			fmt.Println(rawInfo)
			newCount++
		})

		err = res.Body.Close()

		if err != nil {
			panic(err)
		}

		hasMore = false
	}
}
