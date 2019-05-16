package meta

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/google/go-github/v25/github"
)

type Github struct{}

func (Github) Extract(url string) (*MetaData, error) {
	if matched, _ := regexp.MatchString("^https?://github.com*", url); !matched {
		return nil, ErrMetaNotSupported
	}

	fmt.Print("Fetching metadata from github... ")

	re := regexp.MustCompile("^https?://github.com/([^/]*)/([^/]*)")
	result := re.FindStringSubmatch(url)
	if result == nil {
		fmt.Println("Failed")
		return nil, ErrMetaExtraction
	}

	if len(result) < 3 {
		fmt.Println("Failed")
		return nil, ErrMetaExtraction
	}

	user := result[1]
	repo := result[2]

	client := github.NewClient(nil)

	repos, _, err := client.Repositories.Get(context.Background(), user, repo)
	if err != nil {
		fmt.Println("Failed")
		return nil, ErrMetaExtraction
	}

	meta := MetaData{
		Description: strings.TrimSpace(html.UnescapeString(fmt.Sprintf("%s: %s", *repos.FullName, *repos.Description))),
		Tags:        fmt.Sprintf("%s,%s", "github", strings.ToLower(*repos.Language)),
	}

	fmt.Println("Done")
	return &meta, nil
}
