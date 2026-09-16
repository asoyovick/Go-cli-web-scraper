package scraper

import "go-cli-web-scraper/internal/fetcher"


type Scraper struct {
	fetcher *fetcher.Fetcher
}

func New(f *fetcher.Fetcher) *Scraper {
	return &Scraper{fetcher: f}
}