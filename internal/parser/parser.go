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
	var images []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					title = n.FirstChild.Data
				}	
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						resolved := resolveURL(baseURL, attr.Val)
						if resolved != "" {
							links = append(links, resolved)
						}
					}
				}
			case "img":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resolved := resolveURL(baseURL, attr.Val)
						if resolved != "" {
							images = append(images, reolvec)
						}
					}
				}
			
			}
		}
	}
}