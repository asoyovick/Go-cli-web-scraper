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
		return nil, fmt.Errorf("Network error: %w",err)
	}
}