package util

import (
	"net/http"
	"testing"
)

func TestDownloadWebpage(t *testing.T) {
	RegisterMockHTTPClient()

	// test download
	url := "https://www.example.com"
	_, err := DownloadWebpage(url, http.StatusOK, nil, nil)

	if err != nil {
		t.Errorf("cannot download webpage: %v", err)
	}

	// test download with custom headers
	_, err = DownloadWebpage(url, http.StatusOK, map[string]string{
		"x-token": "dummy-token",
	}, nil)

	if err != nil {
		t.Errorf("cannot download webpage with custom headers: %v", err)
	}

	// test download with custom parameters
	_, err = DownloadWebpage(url, http.StatusOK, nil, map[string]string{
		"q": "dummy-query",
	})

	if err != nil {
		t.Errorf("cannot download webpage with custom parameters: %v", err)
	}

	// test download with different expected status, should return error
	_, err = DownloadWebpage(url, http.StatusCreated, nil, nil)

	if err == nil {
		t.Errorf("should return error because of different status code")
	}

	// test download, should return error
	_, err = DownloadWebpage("https://www.test-do-error.com", http.StatusOK, nil, nil)

	if err == nil {
		t.Errorf("should return error after Do method execution")
	}
}
