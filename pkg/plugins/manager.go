package plugins

import "fmt"

func GetPlugin(name string) (ContentSource, error) {
	switch name {
	case "hn":
		return NewHackerNewsSource(), nil
	case "github":
		return NewGitHubSource(), nil
	case "spanish-news":
		return NewSpanishNewsSource(), nil
	default:
		return nil, fmt.Errorf("unknown plugin: %s", name)
	}
}

func ListPlugins() []string {
	return []string{"hn", "github", "spanish-news"}
}
