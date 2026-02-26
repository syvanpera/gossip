// internal/storage/storage.go
package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Bookmark represents a single saved URL and its metadata
type Bookmark struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Comment   string   `json:"comment,omitempty"`
	DateAdded string   `json:"date_added"`
}

// NewBookmark is a helper to create a bookmark with an auto-generated ID and timestamp
func NewBookmark(url, title string, tags []string, comment string) Bookmark {
	return Bookmark{
		ID:        generateID(),
		URL:       url,
		Title:     title,
		Tags:      tags,
		Comment:   comment,
		DateAdded: time.Now().Format(time.RFC3339), // Standard ISO 8601 format
	}
}

// Load reads all bookmarks from the JSON file
func Load(filepath string) ([]Bookmark, error) {
	// If the file doesn't exist yet, just return an empty slice
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return []Bookmark{}, nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read storage file: %w", err)
	}

	// Handle completely empty files to prevent unmarshal errors
	if len(data) == 0 {
		return []Bookmark{}, nil
	}

	var bookmarks []Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return nil, fmt.Errorf("could not parse JSON data: %w", err)
	}

	return bookmarks, nil
}

// Save writes the list of bookmarks back to the JSON file
func Save(filepath string, bookmarks []Bookmark) error {
	// MarshalIndent makes the JSON human-readable and easy to edit by hand
	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal bookmarks: %w", err)
	}

	// 0644 gives read/write permissions to the user, and read-only to others
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("could not write to storage file: %w", err)
	}

	return nil
}

// GetByID retrieves a single bookmark by its ID
func GetByID(filepath string, id string) (Bookmark, error) {
	bookmarks, err := Load(filepath)
	if err != nil {
		return Bookmark{}, err
	}

	for _, b := range bookmarks {
		if b.ID == id {
			return b, nil
		}
	}

	return Bookmark{}, fmt.Errorf("bookmark with ID '%s' not found", id)
}

// Add appends a single new bookmark to the storage file
func Add(filepath string, b Bookmark) error {
	bookmarks, err := Load(filepath)
	if err != nil {
		return err
	}

	bookmarks = append(bookmarks, b)
	return Save(filepath, bookmarks)
}

// Search returns a list of bookmarks that match the given query
func Search(filepath string, query string) ([]Bookmark, error) {
	bookmarks, err := Load(filepath)
	if err != nil {
		return nil, err
	}

	var results []Bookmark
	query = strings.ToLower(query)

	for _, b := range bookmarks {
		// Check the main fields for a substring match
		if strings.Contains(strings.ToLower(b.URL), query) ||
			strings.Contains(strings.ToLower(b.Title), query) ||
			strings.Contains(strings.ToLower(b.Comment), query) {
			results = append(results, b)
			continue // Found a match, move to the next bookmark
		}

		// Check the tags
		for _, tag := range b.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, b)
				break // Found a tag match, stop checking other tags for this bookmark
			}
		}
	}

	return results, nil
}

// Delete removes a bookmark by its ID
func Delete(filepath string, id string) error {
	bookmarks, err := Load(filepath)
	if err != nil {
		return err
	}

	found := false
	for i, b := range bookmarks {
		if b.ID == id {
			// Remove the item by appending the slices before and after it
			bookmarks = append(bookmarks[:i], bookmarks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("bookmark with ID '%s' not found", id)
	}

	return Save(filepath, bookmarks)
}

// Update replaces an existing bookmark with an updated version
func Update(filepath string, updatedBookmark Bookmark) error {
	bookmarks, err := Load(filepath)
	if err != nil {
		return err
	}

	found := false
	for i, b := range bookmarks {
		if b.ID == updatedBookmark.ID {
			bookmarks[i] = updatedBookmark
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("bookmark with ID '%s' not found", updatedBookmark.ID)
	}

	return Save(filepath, bookmarks)
}

// generateID creates a short, random 8-character hex string for the ID
func generateID() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback if random generation fails
		return fmt.Sprintf("%x", time.Now().Unix())
	}
	return hex.EncodeToString(bytes)
}
