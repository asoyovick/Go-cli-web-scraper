package main

import (
	"flag"
	"time"
	"log"
	"fmt"
	"go-cli-web-scraper/internal/scraper"
	"go-cli-web-scraper/internal/fetcher"
)

func main() {
	// Define and configure command-linr flags
	// Default to 10 second timeout if the flag is ommited
	timeout := flag.Duration("timeout", 10*time.Second, "HTTP timeout duration")
	// Parse the flag arguments from os.Args[1:]. must be called before accessing flags
	flag.Parse()
	// Grab the remaining non-flag command-line arguments
	args := flag.Args()
	// Ensure atleast one argument (the target URL or file path) was porvided
	if len(args) <1 {
		log.Fatal("Usage: go run cmd/min.go [options] <URL or local file path>\nExample: go run cmd/main.go https://example.com")
	}
	//Extract the Primary target (URL or local path) from the validated arguments
	target := args[0]
	// Initialize the HTTP fetcher using the duration configured from the cli flag
	f := fetcher.New(*timeout)
	//inject the fetcher into the scraper service ingine
	s := scraper.New(f)

	fmt.Printf("Scraping target: %...\n", target)
	//Execute the core scraping and parsing processing pipline
	data, err := s.Scrape(target)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	//output the extracted webpage metadata directly to the terminal.
	fmt.Printf("\nTitle: %s\n", data.Title)
	fmt.Printf("Links Found(%d):\n", len(data.Links))
	for _, link := range data.Links {
		fmt.Printf(" - %s\n", link)
	}
	fmt.Printf("Images Found (%d:\n)", len(data.Images))
	for _, img := range data.Images {
		fmt.Printf(" -%s\n", img)
	}
}