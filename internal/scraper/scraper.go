package scraper

import "go-cli-web-scraper/internal/fetcher"


type Scraper struct {
	fetcher *fetcher.Fetcher
}

func New(f *fetcher.Fetcher) *Scraper {
	return &Scraper{fetcher: f}
}

func (s *Scraper) Scrape(source string) (*models.PageData, error) {
	stream, resolvedBase, err := s.Fetcher.FetchSource(source)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	title, links, images err := parse.ExtractData(stream, resovedBase)
	if err != nil { 
		return nil, fmt.Errorf ("parssing failed: %w", err)
	}
}