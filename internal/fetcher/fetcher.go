package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Fetcher struct {
	client *http.Client
}

func New(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchSource handles both remote HTTP/HTTPS URLs and local file paths,
// returning an io.ReadCloser that the caller MUST close.
func (f *Fetcher) FetchSource(source string) (io.ReadCloser, string, error) {
	if isURL(source) {
		req, err := http.NewRequest("GET", source, nil)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("User-Agent", "GoScrape/1.0")

		resp, err := f.client.Do(req)
		if err != nil {
			return nil, "", fmt.Errorf("network error: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			return nil, "", fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
		}

		return resp.Body, source, nil
	}

	// Local file fallback
	file, err := os.Open(source)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open local file: %w", err)
	}

	return file, "file://" + source, nil
}

func isURL(source string) bool {
	u, err := url.Parse(source)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https" || strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://")
}
