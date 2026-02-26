// internal/fetcher/github.go
package fetcher

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v83/github"
)

// GitHubFetcher uses the official GitHub API to extract repository metadata
type GitHubFetcher struct{}

// Fetch parses the GitHub URL and queries the API for repo details
func (f *GitHubFetcher) Fetch(urlStr string) (*PageMeta, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Clean the path and split it to find the owner and repo
	// e.g., /charmbracelet/bubbletea -> ["charmbracelet", "bubbletea"]
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")

	// If it's a valid repo URL
	if len(parts) >= 2 {
		owner := parts[0]
		repo := parts[1]

		client := github.NewClient(nil) // Unauthenticated client works fine for public repos
		ctx := context.Background()

		repoMeta, _, err := client.Repositories.Get(ctx, owner, repo)
		if err != nil {
			// If the API fails (e.g., rate limit) or it's a private repo, fallback gracefully to HTML
			htmlF := &HTMLFetcher{}
			return htmlF.Fetch(urlStr)
		}

		// Safely extract the data
		desc := ""
		if repoMeta.Description != nil {
			desc = *repoMeta.Description
		}

		lang := "Unknown"
		if repoMeta.Language != nil {
			lang = *repoMeta.Language
		}

		stars := 0
		if repoMeta.StargazersCount != nil {
			stars = *repoMeta.StargazersCount
		}

		// Grab the topics directly from the repository metadata
		var tags []string
		if repoMeta.Topics != nil {
			tags = repoMeta.Topics
		}

		title := fmt.Sprintf("%s/%s", owner, repo)
		comment := fmt.Sprintf("%s (★ %d | %s)", desc, stars, lang)

		return &PageMeta{
			Title:       title,
			Description: comment,
			Tags:        tags, // Pass the topics back as tags
		}, nil
	}

	// If it's just a user profile or issue page, fallback to HTML
	htmlF := &HTMLFetcher{}
	return htmlF.Fetch(urlStr)
}
