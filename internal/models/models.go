package models

import 
// pagedata represents the extracted information from a single webpage
type PageData struct {
	URl		string `json:"url"`
	Title	string 	`json:"title"`
	Links	[]string	`json:"links"`
	Images	[]string	`json:"images"`
	StatusCode	int		`json:"status_code"`

}

// config holds runtime configuration passed from command-line flags.
type Config struct {
	StartURL	string
	Timeout		time.Duration
}