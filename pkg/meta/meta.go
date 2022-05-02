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

var ErrMetaExtraction = errors.New("error while extracting metadata")
var ErrMetaNotSupported = errors.New("can't handle this")

var extractors []Extractor

type MetaData struct {
	Description string
	Tags        string
}

type Extractor interface {
	Extract(url string) (*MetaData, error)
}

type Generic struct{}

func (Generic) Extract(url string) (*MetaData, error) {
	fmt.Printf("Fetching metadata from %s... ", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed")
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	content := string(data)

	re := regexp.MustCompile("<title>(.*)</title>")
	result := re.FindStringSubmatch(content)
	if result == nil || len(result) < 2 {
		fmt.Println("Failed")
		return nil, errors.New("can't find <title> tag")
	}
	description := strings.TrimSpace(html.UnescapeString(result[1]))

	meta := MetaData{
		Description: description,
	}

	fmt.Println("Done")
	return &meta, nil
}

func Extract(url string) *MetaData {
	for _, e := range extractors {
		meta, err := e.Extract(url)
		if err != nil {
			continue
		}
		return meta
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
