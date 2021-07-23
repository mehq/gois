package util

import (
	"net/http"
	"testing"
)

func TestDownloadWebpage(t *testing.T) {
	RegisterMockHTTPClient()

	_, err := DownloadWebpage("https://www.example.com", http.StatusOK, nil, nil)

	if err != nil {
		t.Errorf("cannot download webpage: %v", err)
	}
}
