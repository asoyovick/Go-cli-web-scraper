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
func TestFetch_Remote(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>hi</body></html>"))
	}))
	defer srv.Close()

	f := New(5 * time.Second)
	rc, resolved, err := f.FetchSource(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rc.Close()

	if resolved != srv.URL {
		t.Errorf("resolved = %q, want %q", resolved, srv.URL)
	}

	body, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if string(body) != "<html><body>hi</body></html>" {
		t.Errorf("body = %q, unexpected content", string(body))
	}
}

func TestFetch_UserAgent(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := New(5 * time.Second)
	rc, _, err := f.FetchSource(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rc.Close()

	if gotUA != "GoScrape/1.0" {
		t.Errorf("User-Agent = %q, want %q", gotUA, "GoScrape/1.0")
	}
}

func TestFetch_Local(t *testing.T) {
	statusCodes := []int{400, 401, 403, 404, 500, 502, 503}

	for _, code := range statusCodes {
		code := code
		t.Run(http.StatusText(code), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))
			defer srv.Close()

			f := New(5 * time.Second)
			rc, _, err := f.FetchSource(srv.URL)
			if err == nil {
				rc.Close()
				t.Fatalf("expected error for status %d, got nil", code)
			}
			if rc != nil {
				t.Errorf("expected nil ReadCloser on error, got non-nil")
			}
		})
	}
}

func TestFetch_StatusBounds(t *testing.T) {
	tests := []struct {
		code    int
		wantErr bool
	}{
		{200, false},
		{299, false},
		{300, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(http.StatusText(tt.code), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.code)
			}))
			defer srv.Close()

			f := New(5 * time.Second)
			rc, _, err := f.FetchSource(srv.URL)
			if tt.wantErr && err == nil {
				rc.Close()
				t.Fatalf("status %d: expected error, got nil", tt.code)
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("status %d: unexpected error: %v", tt.code, err)
				}
				rc.Close()
			}
		})
	}
}
