package scraper

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-cli-web-scraper/internal/fetcher"
)

func TestScrape_Remote(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Test Page</title></head><body>
			<a href="/one">One</a>
			<img src="/logo.png">
		</body></html>`))
	}))
	defer srv.Close()

	s := New(fetcher.New(5 * time.Second))
	data, err := s.Scrape(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.URL != srv.URL {
		t.Errorf("data.URL = %q, want %q (original source, not resolved)", data.URL, srv.URL)
	}
	if data.Title != "Test Page" {
		t.Errorf("data.Title = %q, want %q", data.Title, "Test Page")
	}
	wantLink := srv.URL + "/one"
	if len(data.Links) != 1 || data.Links[0] != wantLink {
		t.Errorf("data.Links = %v, want [%s]", data.Links, wantLink)
	}
	wantImage := srv.URL + "/logo.png"
	if len(data.Images) != 1 || data.Images[0] != wantImage {
		t.Errorf("data.Images = %v, want [%s]", data.Images, wantImage)
	}
}

func TestScrape_Local(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "page.html")
	content := `<html><head><title>Local Page</title></head><body><a href="/x">x</a></body></html>`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	s := New(fetcher.New(5 * time.Second))
	data, err := s.Scrape(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.URL != path {
		t.Errorf("data.URL = %q, want %q", data.URL, path)
	}
	if data.Title != "Local Page" {
		t.Errorf("data.Title = %q, want %q", data.Title, "Local Page")
	}
	if len(data.Links) != 1 {
		t.Fatalf("data.Links = %v, want exactly one entry", data.Links)
	}
}
func TestScrape_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := New(fetcher.New(5 * time.Second))
	data, err := s.Scrape(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Title != "" || data.Links != nil || data.Images != nil {
		t.Errorf("expected empty PageData fields, got %+v", data)
	}
}