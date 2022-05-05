package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/pkg/util"
)

func tableSchema() string {
	return `
CREATE TABLE IF NOT EXISTS bookmarks(
    id INTEGER NOT NULL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    flags INTEGER NOT NULL DEFAULT 0,

    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now'))
)`
}

func GetConnection(path string) *sql.DB {
	firstRun := false
	if !util.PathExists(viper.GetString("database.path")) {
		util.EnsurePath(path)
		firstRun = true
	}

	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println("Database connection failed")
		os.Exit(1)
	}

	if firstRun {
		conn.Exec(tableSchema())
	}

	return conn
}
