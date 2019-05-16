package meta

import (
	"fmt"
	"regexp"
)

type Github struct{}

func (Github) CanHandle(url string) bool {
	matched, _ := regexp.MatchString("^https?://github.com*", url)
	return matched
}

func (Github) Extract(url string) *MetaData {
	fmt.Println("Extract from github")
	return nil
}
