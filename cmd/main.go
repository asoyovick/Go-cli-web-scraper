package main

import (
	"flag"
	"fmt"
<<<<<<< HEAD
	"go-cli-web-scraper/internal/fetcher"
	"go-cli-web-scraper/internal/scraper"
	"log"
	"time"
)

func main() {
	// Define and configure command-linr flags
	// Default to 10 second timeout if the flag is ommited
	timeout := flag.Duration("timeout", 30*time.Second, "HTTP timeout duration")
	// Parse the flag arguments from os.Args[1:]. must be called before accessing flags
	flag.Parse()
	// Grab the remaining non-flag command-line arguments
	args := flag.Args()
	// Ensure atleast one argument (the target URL or file path) was porvided
	if len(args) < 1 {
		log.Fatal("Usage: go run cmd/min.go [options] <URL or local file path>\nExample: go run cmd/main.go https://example.com")
=======
	urlpkg "net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL")
		return
>>>>>>> 697ee6768f4d0b944ac062e2f6a3b70c703d6879
	}

<<<<<<< HEAD
	fmt.Printf("Scraping target: %s...\n", target)
	//Execute the core scraping and parsing processing pipline
	data, err := s.Scrape(target)
=======
	url := os.Args[1]

	_, err := urlpkg.ParseRequestURI(url)
>>>>>>> 697ee6768f4d0b944ac062e2f6a3b70c703d6879
	if err != nil {
		fmt.Println("Invalid URL")
		return
	}

<<<<<<< HEAD
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
=======
	fmt.Println("Valid URL:", url)
>>>>>>> 697ee6768f4d0b944ac062e2f6a3b70c703d6879
}
