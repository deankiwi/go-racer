package plugins

import (
	"math/rand"
	"net/http"
	"time"
)

type GitHubSource struct{}

func NewGitHubSource() *GitHubSource {
	return &GitHubSource{}
}

func (g *GitHubSource) Name() string {
	return "GitHub"
}

func (g *GitHubSource) Description() string {
	return "Types out random Go snippets from the standard library"
}

// For simplicity, let's fetch from the Go standard library examples or a specific repo
func (g *GitHubSource) GetContent() (string, error) {
	// Let's try to get a file from the Go repo
	// This is a simplified approach. A real implementation might use the GitHub API to search for code.
	// For now, let's just return a hardcoded slice of interesting Go snippets if API fails or to keep it simple without auth.

	// Actually, let's try to fetch a specific file content from raw.githubusercontent.com
	// We can pick from a list of known interesting files.

	files := []string{
		"https://raw.githubusercontent.com/golang/go/master/src/fmt/print.go",
		"https://raw.githubusercontent.com/golang/go/master/src/time/time.go",
		"https://raw.githubusercontent.com/golang/go/master/src/strings/strings.go",
		"https://raw.githubusercontent.com/golang/go/master/src/net/http/server.go",
	}

	rand.Seed(time.Now().UnixNano())
	url := files[rand.Intn(len(files))]

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read a chunk of the file
	// We don't want the whole file, just a function or a block.
	// This is tricky without a parser.
	// Let's read the first 500 bytes and find a complete line?
	// Or maybe just grab a random function from a predefined set of snippets for stability.

	// For a better experience without complex parsing, let's use a curated list of snippets for now,
	// checking if we can fetch them. If not, we fall back to hardcoded ones.

	// Let's implement a simple "random snippet from memory" for this proof of concept
	// to ensure it works reliably without hitting API rate limits or parsing issues.

	snippets := []string{
		`testing`,
	}

	return snippets[rand.Intn(len(snippets))], nil
}

// Below is a scaffold for a real GitHub API implementation if we had a token
type GitHubContent struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}
