package snippet

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/util"

	// bring in the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Repository is the main accessor for snippets in the DB
type Repository struct {
	db *sqlx.DB
}

// NewRepository returns a new snippet resository
func NewRepository() *Repository {
	return &Repository{db: openDB(viper.GetString("database"))}
}

func openDB(file string) *sqlx.DB {
	util.EnsureDir(file)
	db, _ := sqlx.Open("sqlite3", file)

	schema := `
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER PRIMARY KEY,
			content TEXT,
			description TEXT,
			tags TEXT,
			type TEXT,
			language TEXT
		)`

	db.MustExec(schema)

	return db
}

// Add a new snippet to the database
func (r *Repository) Add(s Snippet) {
	sd := s.Data()
	query := `
		INSERT INTO snippets (content, description, tags, type, language)
		VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, sd.Content, sd.Description, sd.Tags, sd.Type, sd.Language)
	if err != nil {
		panic(err)
	}

	ID, err := result.LastInsertId()
	sd.ID = ID
	if err != nil {
		panic(err)
	}
}

// Save an existing snippet to the database
func (r *Repository) Save(s Snippet) {
	sd := s.Data()
	query := `
		UPDATE snippets SET
		content = ?, description = ?, tags = ?, language = ?
		WHERE id = ?`
	if _, err := r.db.Exec(query, sd.Content, sd.Description, sd.Tags, sd.Language, sd.ID); err != nil {
		panic(err)
	}
}

// Get returns a snippet with the given ID
func (r *Repository) Get(id int) Snippet {
	var data SnippetData
	err := r.db.Get(&data, "SELECT * FROM snippets WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	return New(data)
}

// FindWithFilters returns snippets that match the given filters
func (r *Repository) FindWithFilters(filters Filters) []Snippet {
	var wheres []string
	if filters.Type != "" {
		wheres = append(wheres, fmt.Sprintf("type = \"%s\"", filters.Type))
	}
	if filters.Language != "" {
		wheres = append(wheres, fmt.Sprintf("language = \"%s\"", filters.Language))
	}
	if filters.Tags != "" {
		tags := strings.Split(filters.Tags, ",")
		for _, tag := range tags {
			wheres = append(wheres, fmt.Sprintf("',' || tags || ',' like '%%,%s,%%'", tag))
		}
	}

	var sb strings.Builder
	sb.WriteString("SELECT * FROM snippets")
	if len(wheres) > 0 {
		fmt.Fprintf(&sb, " WHERE %s", strings.Join(wheres, " AND "))
	}

	rows, err := r.db.Queryx(sb.String())
	if err != nil {
		panic(err)
	}

	var ss []Snippet
	for rows.Next() {
		var data SnippetData
		if err := rows.StructScan(&data); err != nil {
			panic(err)
		}
		ss = append(ss, New(data))
	}

	return ss
}

// Del removes a snippet with the given ID
func (r *Repository) Del(id int) {
	r.db.MustExec("DELETE FROM snippets WHERE id = $1", id)
}
