package fetcher

import (
	"net/http"
	"time"
	"fmt"
)

type Fetcher struct {
	client *http.Client
}

func New(timeout time.Duration) *Fetcher {
	return &Fetcher {
		client: &http.Client{
			Timeout: timeout,
		},
	}
}
//pattern retrieves the raw HTTP response body from URL.

func (f *Fetcher) Fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "GoScrape/1.0")
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w",err)
	}
	// HTTPS code outside the 200 -299 range is treated as an application error.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close() // the caller must close the body when done
		return nil, fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}
	return resp, nil
}