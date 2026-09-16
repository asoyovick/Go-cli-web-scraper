package main

import (
	"flag"
	"time"
	"log"
	"fmt"
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
	target := args[0]
	f := fetcher.New(*timeout)
	s := Scraper.New(f)

	fmt.Printf("Scraping target: %...\n", target)
	data, err := s.Scrape(target)
}