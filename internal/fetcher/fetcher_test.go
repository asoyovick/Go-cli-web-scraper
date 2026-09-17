package fetcher

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   bool
	}{
		{"http url", "http://example.com", true},
		{"https url", "https://example.com/page", true},
		{"https url with query", "https://example.com/page?q=1", true},
		{"ftp scheme", "ftp://example.com/file", false},
		{"absolute unix path", "/home/user/file.html", false},
		{"relative path", "index.html", false},
		{"relative path with dots", "../pages/index.html", false},
		{"empty string", "", false},
		{"bare domain no scheme", "example.com", false},
		{"windows path", `C:\Users\test\file.html`, false},
		{"mailto scheme", "mailto:test@example.com", false},
		{"scheme-prefixed but garbage host", "http://", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isURL(tt.source)
			if got != tt.want {
				t.Errorf("isURL(%q) = %v, want %v", tt.source, got, tt.want)
			}
		})
	}
}

func TestNew_Timeout(t *testing.T) {
	f := New(5 * time.Second)
	if f == nil {
		t.Fatal("New returned nil")
	}
	if f.client == nil {
		t.Fatal("New did not initialize an http.Client")
	}
	if f.client.Timeout != 5*time.Second {
		t.Errorf("client.Timeout = %v, want %v", f.client.Timeout, 5*time.Second)
	}
}
