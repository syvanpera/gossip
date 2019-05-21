package snippet

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/util"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteRepository struct {
	db *sqlx.DB
}

func NewSQLiteRepository() Repository {
	return &sqliteRepository{db: openDB(viper.GetString("database"))}
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

func (r *sqliteRepository) Create(s Snippet) error {
	sd := s.Data()
	query := `
		INSERT INTO snippets (content, description, tags, type, language)
		VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, sd.Content, sd.Description, sd.Tags, sd.Type, sd.Language)
	if err != nil {
		return err
	}

	ID, err := result.LastInsertId()

	sd.ID = ID
	if err != nil {
		return err
	}

	return nil
}

func (r *sqliteRepository) Update(s Snippet) error {
	sd := s.Data()
	query := `
		UPDATE snippets SET
		content = ?, description = ?, tags = ?, language = ?
		WHERE id = ?`
	if _, err := r.db.Exec(query, sd.Content, sd.Description, sd.Tags, sd.Language, sd.ID); err != nil {
		return err
	}

	return nil
}

func (r *sqliteRepository) Get(id int) (Snippet, error) {
	var data SnippetData
	err := r.db.Get(&data, "SELECT * FROM snippets WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return NewSnippet(data), nil
}

func (r *sqliteRepository) FindWithFilters(filters Filters) ([]Snippet, error) {
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
		return nil, err
	}

	var ss []Snippet
	for rows.Next() {
		var data SnippetData
		if err := rows.StructScan(&data); err != nil {
			return nil, err
		}
		ss = append(ss, NewSnippet(data))
	}

	return ss, nil
}

func (r *sqliteRepository) Del(id int) error {
	_, err := r.db.Exec("DELETE FROM snippets WHERE id = $1", id)

	return err
}
