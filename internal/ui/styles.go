// internal/ui/styles.go
package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Status Colors
	StyleRunning = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)  // Cornflower Blue
	StyleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)  // Green
	StyleFailed  = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red

	// List Separators
	StyleID      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))                // Dark Gray
	StyleTitle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true)     // Light Gray
	StyleURL     = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Underline(true) // Blue
	StyleTags    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Italic(true)   // Orange
	StyleComment = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))                // Medium Gray
)
