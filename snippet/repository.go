package snippet

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/util"

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
	util.EnsureDir(file)
	db, _ := sql.Open("sqlite3", file)
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER PRIMARY KEY,
			content TEXT,
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

// Upsert either inserts a new snippet or updates an existing one
func (r *Repository) Upsert(s Snippet) error {
	sd := s.Data()
	stmt, _ := r.db.Prepare(`
		INSERT INTO snippets (content, description, tags, type, language)
		VALUES (?, ?, ?, ?, ?)`)
	_, err := stmt.Exec(sd.Content, sd.Description, strings.Join(sd.Tags, ","), sd.Type, sd.Language)

	if err == nil {
		row := r.db.QueryRow("SELECT last_insert_rowid()")
		var ID int
		err = row.Scan(&ID)
		if err == nil {
			sd.ID = ID
		}
	}

	return err
}

// Get returns a snippet with the given ID
func (r *Repository) Get(id int) Snippet {
	var ID int
	var content, description, tags, _type, language sql.NullString

	row := r.db.QueryRow(`
		SELECT id, content, description, tags, type, language
		FROM snippets WHERE id = $1`, id)
	err := row.Scan(&ID, &content, &description, &tags, &_type, &language)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	sd := SnippetData{
		ID:          ID,
		Content:     content.String,
		Description: description.String,
		Tags:        strings.Split(tags.String, ","),
		Type:        SnippetType(_type.String),
		Language:    language.String,
	}

	return New(sd)
}

// FindAll returns all snippets
func (r *Repository) FindAll() []Snippet {
	rows, err := r.db.Query(`
		SELECT id, content, description, tags, type, language
		FROM snippets`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ss []Snippet
	for rows.Next() {
		var ID int
		var content, description, tags, _type, language sql.NullString
		if err := rows.Scan(&ID, &content, &description, &tags, &_type, &language); err != nil {
			log.Fatal(err)
		}
		s := SnippetData{
			ID:          ID,
			Content:     content.String,
			Description: description.String,
			Tags:        strings.Split(tags.String, ","),
			Type:        SnippetType(_type.String),
			Language:    language.String,
		}
		ss = append(ss, New(s))
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

	var sb strings.Builder
	sb.WriteString("SELECT id, content, description, tags, type, language FROM snippets")
	if len(wheres) > 0 {
		fmt.Fprintf(&sb, " WHERE %s", strings.Join(wheres, " AND "))
	}
	fmt.Println(sb.String())

	rows, err := r.db.Query(sb.String())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ss []Snippet
	for rows.Next() {
		var ID int
		var content, description, tags, _type, language sql.NullString
		if err := rows.Scan(&ID, &content, &description, &tags, &_type, &language); err != nil {
			log.Fatal(err)
		}
		s := SnippetData{
			ID:          ID,
			Content:     content.String,
			Description: description.String,
			Tags:        strings.Split(tags.String, ","),
			Type:        SnippetType(_type.String),
			Language:    language.String,
		}
		ss = append(ss, New(s))
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ss
}

// FindWithTag returns snippets with given tag
func (r *Repository) FindWithTag(tag string) []Snippet {
	rows, err := r.db.Query(`SELECT id, content, description, tags, type, language
FROM snippets
WHERE tags like $1`, fmt.Sprintf("%%%s%%", tag))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ss []Snippet
	for rows.Next() {
		var ID int
		var content, description, tags, _type, language sql.NullString
		if err := rows.Scan(&ID, &content, &description, &tags, &_type, &language); err != nil {
			log.Fatal(err)
		}
		s := SnippetData{
			ID:          ID,
			Content:     content.String,
			Description: description.String,
			Tags:        strings.Split(tags.String, ","),
			Type:        SnippetType(_type.String),
			Language:    language.String,
		}
		ss = append(ss, New(s))
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
