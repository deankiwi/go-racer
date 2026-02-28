package plugins

import (
	"encoding/json"
	"fmt"
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
	return "Types out trending repositories and their descriptions"
}

type GitHubSearchResponse struct {
	Items []struct {
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		HTMLURL     string `json:"html_url"`
	} `json:"items"`
}

func (g *GitHubSource) GetContent() (*Content, error) {
	// Fetch trending repositories created in the last 7 days
	// We use the search API to get the top 50 most starred repositories
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=created:>%s&sort=stars&order=desc&per_page=50", sevenDaysAgo)

	// Create a custom client with a short timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// GitHub API requires a User-Agent header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-racer-terminal-app")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var searchResp GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	if len(searchResp.Items) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	// Pick a random repository from the top 50
	rand.Seed(time.Now().UnixNano())
	repo := searchResp.Items[rand.Intn(len(searchResp.Items))]

	// Handle cases where description might be missing
	desc := repo.Description
	if desc == "" {
		desc = "No description provided."
	}

	// Format the text to type out
	text := fmt.Sprintf("%s\n\n%s", repo.FullName, desc)

	return &Content{
		Text:      text,
		SourceURL: repo.HTMLURL,
	}, nil
}
