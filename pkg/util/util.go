package util

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

func ConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to resolve user config directory.")
	}

	return dir
}

func DataPath() string {
	dir, err := UserDataDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to resolve user data directory.")
	}

	return dir
}

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

func Contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}

	return false
}

func PrintError(err error) {
	colors := viper.GetBool("config.color") != viper.GetBool("color")
	au := aurora.NewAurora(colors)

	fmt.Println(fmt.Sprintf("%s %s", au.Bold(au.BrightRed("[ERROR]")), err))
}

func UserDataDir() (string, error) {
	var dir string

	// TODO Figure out the correct path for different environments
	switch runtime.GOOS {
	// case "windows":
	// 	dir = Getenv("LocalAppData")
	// 	if dir == "" {
	// 		return "", errors.New("%LocalAppData% is not defined")
	// 	}

	// case "darwin", "ios":
	// 	dir = Getenv("HOME")
	// 	if dir == "" {
	// 		return "", errors.New("$HOME is not defined")
	// 	}
	// 	dir += "/Library/Caches"

	// case "plan9":
	// 	dir = Getenv("home")
	// 	if dir == "" {
	// 		return "", errors.New("$home is not defined")
	// 	}
	// 	dir += "/lib/cache"

	default: // Unix
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_DATA_HOME nor $HOME are defined")
			}
			dir += "/.local/share"
		}
	}

	return dir, nil
}

func NormalizeTags(tags string) string {
	tagsSlice := strings.Split(tags, ",")
	unique := map[string]bool{}

	for v := range tagsSlice {
		t := strings.ToLower(strings.TrimSpace(tagsSlice[v]))
		if t != "" {
			unique[t] = true
		}
	}

	result := []string{}
	for key := range unique {
		result = append(result, key)
	}

	return strings.Join(result, ",")
}
