package util

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func EnsureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, err := os.Stat(dirName); err != nil {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func CenterStr(s string, w int) string {
	offset := (w + len(s)) / 2
	if (w+len(s))%2 != 0 {
		offset++
	}
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", offset, s))
}

func GetTerminalSize() (int, int) {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		width = 80
		height = 25
	}

	return width, height
}

func ReplaceRuneAtIndex(s string, r rune, i int) string {
	out := []rune(s)
	out[i] = r
	return string(out)
}

func StrPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] +
			input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

func ExtractTitleFromURL(url string) string {
	response, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	titleStartIndex := strings.Index(pageContent, "<title>")
	if titleStartIndex == -1 {
		return ""
	}
	titleStartIndex += 7

	titleEndIndex := strings.Index(pageContent, "</title>")
	if titleEndIndex == -1 {
		return ""
	}

	return pageContent[titleStartIndex:titleEndIndex]
}
