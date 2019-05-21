package snippet

import (
	"errors"
	"regexp"

	"github.com/syvanpera/gossip/meta"
	"github.com/syvanpera/gossip/ui"
)

var (
	ErrSomething = errors.New("some error")
)

type Service interface {
	AddBookmark(content, description, tags string) (Snippet, error)
	AddCommand(content, description, tags string) (Snippet, error)
	AddCode(content, description, tags string) (Snippet, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddBookmark(content, description, tags string) (Snippet, error) {
	data := SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        COMMAND,
	}

	if content == "" {
		if content = ui.Prompt("URL", content); content == "" {
			return nil, ErrSomething
		}
	}

	if matched, _ := regexp.MatchString("^https?://*", content); !matched {
		content = "http://" + content
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
			return nil, ErrSomething
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

func (s *service) AddCommand(content, description, tags string) (Snippet, error) {
	data := SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        COMMAND,
	}

	if content == "" {
		if content = ui.Prompt("Command", content); content == "" {
			return nil, ErrSomething
		}
	}

	if description == "" {
		if description = ui.Prompt("Description", description); description == "" {
			return nil, ErrSomething
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

func (s *service) AddCode(content, description, tags string) (Snippet, error) {
	data := SnippetData{
		Content:     "",
		Description: "",
		Tags:        tags,
		Language:    "",
		Type:        CODE,
	}

	if description == "" {
		if description = ui.Prompt("Description", description); description == "" {
			return nil, ErrSomething
		}
	}

	if content = ui.Editor(""); content == "" {
		return nil, ErrSomething
	}

	language := ui.Prompt("Language", "")

	data.Content = content
	data.Description = description
	data.Tags = tags
	data.Language = language

	snippet := NewSnippet(data)
	if err := s.repository.Create(snippet); err != nil {
		return nil, err
	}

	return snippet, nil
}
