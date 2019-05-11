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
func (r *Repository) New(s *Snippet) error {
	stmt, _ := r.db.Prepare("INSERT INTO snippets (snippet, description, tags, type, language) VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(s.Snippet, s.Description, nil, s.Type, nil)

	if err == nil {
		row := r.db.QueryRow("SELECT last_insert_rowid()")
		var ID int
		err = row.Scan(&ID)
		if err == nil {
			s.ID = ID
		}
	}

	return err
}

// Get returns a snippet with the given ID
func (r *Repository) Get(id int) *Snippet {
	row := r.db.QueryRow(`SELECT id, snippet, description, tags, type, language
FROM snippets WHERE id = $1`, id)

	var ID int
	var snippet, description, tags, _type, language sql.NullString
	err := row.Scan(&ID, &snippet, &description, &tags, &_type, &language)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
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

	return &s
}

// FindAll returns all snippets
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

// FindWithFilters returns snippets that match the given filters
func (r *Repository) FindWithFilters(filters Filters) []Snippet {
	var wheres []string
	fmt.Println(wheres)
	if filters.Type != "" {
		wheres = append(wheres, fmt.Sprintf("type = \"%s\"", filters.Type))
	}
	if filters.Language != "" {
		wheres = append(wheres, fmt.Sprintf("language = \"%s\"", filters.Language))
	}

	var qb strings.Builder
	qb.WriteString("SELECT id, snippet, description, tags, type, language FROM snippets")
	if len(wheres) > 0 {
		fmt.Fprintf(&qb, " WHERE %s", strings.Join(wheres, " AND "))
	}
	fmt.Println(qb.String())

	rows, err := r.db.Query(qb.String())
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

// FindWithTag returns snippets with given tag
func (r *Repository) FindWithTag(tag string) []Snippet {
	rows, err := r.db.Query(`SELECT id, snippet, description, tags, type, language
FROM snippets
WHERE tags like $1`, fmt.Sprintf("%%%s%%", tag))
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

// Del removes a snippet with the given ID
func (r *Repository) Del(id int) error {
	_, err := r.db.Exec("DELETE FROM snippets WHERE id = $1", id)

	return err
}
