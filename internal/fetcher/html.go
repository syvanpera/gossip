// internal/fetcher/html.go
package fetcher

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// HTMLFetcher is the default generic web scraper
type HTMLFetcher struct{}

// Fetch extracts the title and description using standard HTML parsing
func (f *HTMLFetcher) Fetch(urlStr string) (*PageMeta, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid URL request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; BookmarkCLI/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received bad status code: %d", resp.StatusCode)
	}

	meta := &PageMeta{}
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			meta.Title = strings.TrimSpace(meta.Title)
			meta.Description = strings.TrimSpace(meta.Description)
			return meta, nil

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			if token.Data == "title" {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					meta.Title = tokenizer.Token().Data
				}
				continue
			}

			if token.Data == "meta" {
				var isDesc bool
				var content string
				for _, attr := range token.Attr {
					if attr.Key == "name" && (attr.Val == "description" || attr.Val == "og:description") {
						isDesc = true
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if isDesc && content != "" && meta.Description == "" {
					meta.Description = content
				}
			}
		}
	}
}
