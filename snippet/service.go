package snippet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/syvanpera/gossip/meta"
	"github.com/syvanpera/gossip/ui"
	"github.com/syvanpera/gossip/util"
)

var (
	ErrCanceled = errors.New("canceled")
	ErrNotFound = errors.New("snippet not found")
)

type Service interface {
	GetSnippet(id int) (Snippet, error)
	CreateBookmark(content, description, tags string) (Snippet, error)
	CreateCommand(content, description, tags string) (Snippet, error)
	CreateCode(content, description, tags string, language string) (Snippet, error)
	UpdateSnippet(id int, tags string) (Snippet, error)
	DeleteSnippet(id int, force bool) error
	FindSnippets(f Filters) ([]Snippet, error)
	ExecuteSnippet(id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) GetSnippet(id int) (Snippet, error) {
	return s.repository.Get(id)
}

func (s *service) CreateBookmark(content, description, tags string) (Snippet, error) {
	data := SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        BOOKMARK,
	}

	if content == "" {
		if content = ui.Prompt("URL", content); content == "" {
			return nil, ErrCanceled
		}
	}

	if matched, _ := regexp.MatchString("^https?://*", content); !matched {
		content = "https://" + content
	}

	if description == "" || tags == "" {
		if meta := meta.Extract(content); meta != nil {
			if description == "" {
				description = meta.Description
			}
			if tags == "" {
				tags = meta.Tags
			}
		}
	}

	if description == "" {
		if description = ui.Prompt("Description", description); description == "" {
			return nil, ErrCanceled
		}
	}

	data.Content = content
	data.Description = description
	data.Tags = tags

	snippet := NewSnippet(data)
	if err := s.repository.Create(snippet); err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *service) CreateCommand(content, description, tags string) (Snippet, error) {
	data := SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        COMMAND,
	}

	if content == "" {
		if content = ui.Prompt("Command", content); content == "" {
			return nil, ErrCanceled
		}
	}

	if description == "" {
		if description = ui.Prompt("Description", description); description == "" {
			return nil, ErrCanceled
		}
	}

	data.Content = content
	data.Description = description
	data.Tags = tags

	snippet := NewSnippet(data)
	if err := s.repository.Create(snippet); err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *service) CreateCode(content, description, tags string, language string) (Snippet, error) {
	data := SnippetData{
		Content:     "",
		Description: "",
		Tags:        "",
		Language:    "",
		Type:        CODE,
	}

	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		content = string(bytes)
	}

	if description == "" {
		if description = ui.Prompt("Description", description); description == "" {
			return nil, ErrCanceled
		}
	}

	if content == "" {
		if content = ui.Editor(""); content == "" {
			return nil, ErrCanceled
		}
	}

	if language == "" {
		language = ui.Prompt("Language", "")
	}

	data.Content = content
	data.Description = description
	if tags == "" {
		tags = language
	} else {
		tags = language + "," + tags
	}
	data.Tags = tags
	data.Language = language

	snippet := NewSnippet(data)
	if err := s.repository.Create(snippet); err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *service) UpdateSnippet(id int, tags string) (Snippet, error) {
	snippet, err := s.GetSnippet(id)
	if err != nil {
		return nil, err
	}

	if tags == "" {
		if err := snippet.Edit(); err != nil {
			return nil, err
		}
	}

	existingTags := strings.Split(snippet.Data().Tags, ",")
	newTags := strings.Split(tags, ",")

	for _, t := range newTags {
		tag := strings.ToLower(t)
		if !util.Contains(existingTags, tag) {
			existingTags = append(existingTags, tag)
		}
	}
	snippet.Data().Tags = strings.Trim(strings.Join(existingTags, ","), " ,")
	if err := s.repository.Update(snippet); err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *service) DeleteSnippet(id int, force bool) error {
	snippet, err := s.GetSnippet(id)
	if err != nil {
		return err
	}

	fmt.Println(snippet.Render())

	if !force && !ui.Confirm(fmt.Sprintf("Are you sure you want to delete this %s", snippet.String(false))) {
		return ErrCanceled
	}

	return s.repository.Delete(id)
}

func (s *service) FindSnippets(filters Filters) ([]Snippet, error) {
	snippets, err := s.repository.FindWithFilters(filters)
	if err != nil {
		return nil, err
	}
	return snippets, nil
}

func (s *service) ExecuteSnippet(id int) error {
	snippet, err := s.GetSnippet(id)
	if err != nil {
		return err
	}

	return snippet.Execute()
}
