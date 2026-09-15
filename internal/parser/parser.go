package parser

import (
	"io"
	"net/url"
	"golang.org/x/net/html"
	"fmt"
)
// sets up a function to accept the io reader and baseline url string
func ExtractData(r io.Reader, rawBaseURL string)(string, []string, []string, error) {
	baseURL, _ := url.Parse(rawBaseURL)
	// takes raw text urls and conerts to structured url.URL object.
	doc, err :=html.Parse(r)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	var title string
	var links []string
	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
				case "title"
			}
		}
	}
}