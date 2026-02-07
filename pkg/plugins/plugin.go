package plugins

// ContentSource defines the interface for data sources that provide text to type.
type ContentSource interface {
	// Name returns the display name of the plugin
	Name() string
	// GetContent returns a string of text for the user to type
	GetContent() (string, error)
	// Description returns a brief description of what the plugin provides
	Description() string
}
