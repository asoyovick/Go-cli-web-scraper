package scraper

import( 
	"go-cli-web-scraper/internal/fetcher"
	"go-cli-web-scraper/internal/models"
	"go-cli-web-scraper/internal/parser"
	"fmt"
)


type Scraper struct {
	fetcher *fetcher.Fetcher
}

func New(f *fetcher.Fetcher) *Scraper {
	return &Scraper{fetcher: f}
}

func (s *Scraper) Scrape(source string) (*models.PageData, error) {
	stream, resolvedBase, err := s.fetcher.FetchSource(source)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	title, links, images, err := parser.ExtractData(stream, resolvedBase)
	if err != nil { 
		return nil, fmt.Errorf ("parssing failed: %w", err)
	}
	return &models.PageData{
		URL:	source,
		Title:	title,
		Links:	links,
		Images:	images,
	}, nil
}