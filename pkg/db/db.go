package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/syvanpera/gossip/pkg/util"
)

func GetConnection(path string) (*sql.DB, error) {
	util.EnsureDir(path)

	return sql.Open("sqlite3", path)
}
