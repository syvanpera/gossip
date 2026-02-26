// internal/ui/spinner.go
package ui

import (
	"fmt"

	"github.com/syvanpera/gossip/internal/fetcher"
	"github.com/syvanpera/gossip/internal/storage"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// successMsg and errMsg are internal messages passed to the Update loop
type successMsg struct {
	title string
}
type errMsg struct {
	err error
}

// AddModel holds the state for our Bubble Tea program
type AddModel struct {
	spinner  spinner.Model
	url      string
	tags     []string
	filepath string
	err      error
	done     bool
	title    string
}

// NewAddModel initializes the spinner model
func NewAddModel(url string, tags []string, filepath string) AddModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = StyleRunning

	return AddModel{
		spinner:  s,
		url:      url,
		tags:     tags,
		filepath: filepath,
	}
}

// Init starts the spinner ticking and kicks off the background fetch task
func (m AddModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchAndSave())
}

// Update handles incoming events (ticks, background task completion, keyboard interrupts)
func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case successMsg:
		m.done = true
		m.title = msg.title
		return m, tea.Quit

	case errMsg:
		m.err = msg.err
		m.done = true
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the current state to the terminal
func (m AddModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s %s\n", StyleFailed.Render("✗ Failed:"), m.err.Error())
	}
	if m.done {
		return fmt.Sprintf("%s Added %s\n", StyleSuccess.Render("✓ Success!"), StyleTitle.Render(m.title))
	}

	// While running, show the spinner and the "Working..." text
	return fmt.Sprintf("%s %s\n", m.spinner.View(), StyleRunning.Render("Working... Fetching metadata"))
}

// fetchAndSave is a Bubble Tea Command that runs asynchronously
func (m AddModel) fetchAndSave() tea.Cmd {
	return func() tea.Msg {
		// 1. Get the correct fetcher based on the URL
		f := fetcher.GetFetcher(m.url)
		meta, err := f.Fetch(m.url)
		if err != nil {
			// Fallback: If fetch fails entirely
			meta = &fetcher.PageMeta{Title: m.url, Description: ""}
		}

		title := meta.Title
		if title == "" {
			title = m.url
		}

		// 2. Decide which tags to use
		finalTags := m.tags
		if len(finalTags) == 0 && len(meta.Tags) > 0 {
			finalTags = meta.Tags // Fallback to fetched tags if none were provided
		}

		// 3. Create and save the bookmark using finalTags
		bm := storage.NewBookmark(m.url, title, finalTags, meta.Description)
		if err := storage.Add(m.filepath, bm); err != nil {
			return errMsg{err}
		}

		// 4. Return success
		return successMsg{title: title}
	}
}
