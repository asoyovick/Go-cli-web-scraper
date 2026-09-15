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
}