package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DownloadWebpage downloads a page and returns its contents as byte array
func DownloadWebpage(url string, expectedStatus int, headers *map[string]string, params *map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Add(k, v)
		}
	}

	if params != nil {
		q := req.URL.Query()

		for k, v := range *params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}

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

// DownloadJson downloads a webpage and parses content into target
func DownloadJson(url string, expectedStatus int, headers *map[string]string, params *map[string]string, target interface{}) error {
	content, err := DownloadWebpage(url, expectedStatus, headers, params)

	if err != nil {
		return err
	}

	err = json.Unmarshal(content, target)

	if err != nil {
		return err
	}

	return nil
}
