package request

import (
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	if NewHTTPClient() == nil {
		t.Errorf("NewHTTPClient returned nil, expected a valid *http.Client")
	}
}

func TestNewRequest(t *testing.T) {
	if NewRequest("GET", "https://example.com", nil) == nil {
		t.Errorf("NewRequest returned nil, expected a valid *http.Request")
	}
}
