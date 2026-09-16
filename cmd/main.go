package main

import (
	"flag"
	"time"
	"log"
	"fmt"
	urlpkg "net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL")
		return
	}

	url := os.Args[1]

	_, err := urlpkg.ParseRequestURI(url)
	if err != nil {
		fmt.Println("Invalid URL")
		return
	}

	fmt.Println("Valid URL:", url)
}
