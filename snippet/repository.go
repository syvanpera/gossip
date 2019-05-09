package snippet

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"

	// bring in the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Repository is the main accessor for snippets in the DB
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new snippet resository
func NewRepository() *Repository {
	return &Repository{db: openDB(viper.GetString("database"))}
}

func openDB(file string) *sql.DB {
	EnsureDir(file)
	db, _ := sql.Open("sqlite3", file)
	stmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS snippets (
id INTEGER PRIMARY KEY,
snippet TEXT,
description TEXT,
tags TEXT,
type TEXT,
language TEXT
)`)
	_, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// New adds a snippet to the DB
func (r *Repository) New(snippet *Snippet) error {
	stmt, _ := r.db.Prepare("INSERT INTO snippets (snippet, description, tags, type, language) VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(snippet.Snippet, snippet.Description, nil, snippet.Type, nil)

	return err
}

// FindAll returns all snipets
func (r *Repository) FindAll() []Snippet {
	rows, err := r.db.Query(`SELECT id, snippet, description, tags, type, language FROM snippets`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ss []Snippet
	for rows.Next() {
		var ID int
		var snippet, description, tags, _type, language sql.NullString
		if err := rows.Scan(&ID, &snippet, &description, &tags, &_type, &language); err != nil {
			log.Fatal(err)
		}
		s := Snippet{
			ID:          ID,
			Snippet:     snippet.String,
			Description: description.String,
			Tags:        strings.Split(tags.String, ","),
			Type:        _type.String,
			Language:    language.String,
		}
		ss = append(ss, s)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ss
}

// Get returns a snippet with the given ID
func (r *Repository) Get(id int) *Snippet {
	row := r.db.QueryRow(`SELECT id, snippet, description, tags, type, language
FROM snippets WHERE id = $1`, id)

	var s Snippet
	var tags string
	err := row.Scan(&s.ID, &s.Snippet, &s.Description, &tags, &s.Type, &s.Language)
	if err == sql.ErrNoRows {
		fmt.Printf("Snippet with ID %d not found\n", id)
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	s.Tags = strings.Split(tags, ",")

	return &s
}
