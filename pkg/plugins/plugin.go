package plugins

// Content represents the data returned by a plugin
type Content struct {
	Text      string
	SourceURL string // Optional URL
	Author    string // Optional
}

// ContentSource defines the interface for data sources that provide text to type.
type ContentSource interface {
	// Name returns the display name of the plugin
	Name() string
	// GetContent returns text for the user to type and optional metadata
	GetContent() (*Content, error)
	// Description returns a brief description of what the plugin provides
	Description() string
}
