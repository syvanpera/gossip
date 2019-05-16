package meta

import (
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var ErrMetaExtraction = errors.New("can't extract metadata")

var extractors []Extractor

type MetaData struct {
	Description string
	Tags        string
}

type Extractor interface {
	CanHandle(url string) bool
	Extract(url string) (*MetaData, error)
}

type Generic struct{}

func (Generic) CanHandle(url string) bool { return true }

func (Generic) Extract(url string) (*MetaData, error) {
	fmt.Printf("Fetching description from %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	content := string(data)

	re := regexp.MustCompile("<title>(.*)</title>")
	result := re.FindStringSubmatch(content)
	if result == nil || len(result) < 2 {
		return nil, errors.New("can't find <title> tag")
	}
	description := strings.TrimSpace(html.UnescapeString(result[1]))

	meta := MetaData{
		Description: description,
	}

	return &meta, nil
}

func Extract(url string) *MetaData {
	for _, e := range extractors {
		if e.CanHandle(url) {
			meta, _ := e.Extract(url)
			return meta
		}
	}

	return nil
}

func registerExtractor(e Extractor) {
	extractors = append(extractors, e)
}

func init() {
	registerExtractor(Github{})
	// generic must always be the last extractor, as that will handle anything
	registerExtractor(Generic{})
}
