package main

import (
	"flag"
	"fmt"
	"go-cli-web-scraper/internal/fetcher"
	"go-cli-web-scraper/internal/scraper"
	"log"
	"time"
)

func main() {
	// Define and configure command-line flags
	// Default to 10 second timeout if the flag is omitted
	timeout := flag.Duration("timeout", 30*time.Second, "HTTP timeout duration")
	// Parse the flag arguments from os.Args[1:]. must be called before accessing flags
	flag.Parse()
	// Grab the remaining non-flag command-line arguments
	args := flag.Args()
	// Ensure at least one argument (the target URL or file path) was provided
	if len(args) < 1 {
		log.Fatal("Usage: go run cmd/main.go [options] <URL or local file path>\nExample: go run cmd/main.go https://example.com")
	}

	target := args[0]
	f := fetcher.New(*timeout)
	s := scraper.New(f)

	fmt.Printf("Scraping target: %s...\n", target)
	// Execute the core scraping and parsing processing pipeline
	data, err := s.Scrape(target)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Output the extracted webpage metadata directly to the terminal.
	fmt.Printf("\nTitle: %s\n", data.Title)
	fmt.Printf("Links Found(%d):\n", len(data.Links))
	for _, link := range data.Links {
		fmt.Printf(" - %s\n", link)
	}

	fmt.Printf("Images Found (%d):\n", len(data.Images))
	for _, img := range data.Images {
		fmt.Printf(" - %s\n", img)
	}
}
