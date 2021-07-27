package util

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	RegisterMockHTTPClient()

	os.Exit(m.Run())
}

func TestDownloadWebpage(t *testing.T) {
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

func TestSetCookie(t *testing.T) {
	SetCookie(
		"dummy-cookie",
		"dummy-value",
		"/",
		".example.com",
		"https://www.example.com",
	)

	// test download after setting cookie
	url := "https://www.example.com"
	_, err := DownloadWebpage(url, http.StatusOK, nil, nil)

	if err != nil {
		t.Errorf("cannot download webpage after setting cookie: %v", err)
	}
}
