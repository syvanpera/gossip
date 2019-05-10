package snippet

// Snippet contains the data for one snippet
type Snippet struct {
	ID          int      `json:"id"`
	Snippet     string   `json:"snippet"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
	Language    string   `json:"language,omitempty"`
}

func (s Snippet) String() string {
	switch s.Type {
	case "code", "cmd":
		return renderCode(s)
	default:
		return render(s)
	}
}
