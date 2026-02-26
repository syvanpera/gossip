// internal/fetcher/fetcher.go
package fetcher

import (
	"net/url"
	"strings"
)

// PageMeta holds the extracted data from a webpage
type PageMeta struct {
	Title       string
	Description string
	Tags        []string // Added to hold automatically extracted tags
}

// Fetcher defines the contract for all metadata extractors
type Fetcher interface {
	Fetch(urlStr string) (*PageMeta, error)
}

// GetFetcher analyzes the URL and returns the appropriate Fetcher implementation
func GetFetcher(urlStr string) Fetcher {
	parsed, err := url.Parse(urlStr)
	if err == nil {
		host := strings.ToLower(parsed.Host)
		if host == "github.com" || host == "www.github.com" {
			return &GitHubFetcher{}
		}
		// You can easily add more here later (e.g., twitter.com -> TwitterFetcher)
	}

	// Fallback to the standard HTML web scraper
	return &HTMLFetcher{}
}
