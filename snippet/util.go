package snippet

import (
	"database/sql"
	"os"
	"path/filepath"
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
