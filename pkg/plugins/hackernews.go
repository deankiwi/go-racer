package plugins

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type HNStory struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type HackerNewsSource struct{}

func NewHackerNewsSource() *HackerNewsSource {
	return &HackerNewsSource{}
}

func (h *HackerNewsSource) Name() string {
	return "Hacker News"
}

func (h *HackerNewsSource) Description() string {
	return "Types out the titles of top Hacker News stories"
}

func (h *HackerNewsSource) GetContent() (string, error) {
	// Fetch top stories
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		return "", err
	}

	if len(storyIDs) == 0 {
		return "", fmt.Errorf("no stories found")
	}

	// Get a random story ID from the top 50
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(min(50, len(storyIDs)))
	storyID := storyIDs[randomIndex]

	// Fetch the story details
	storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyID)
	resp, err = http.Get(storyURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var story HNStory
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return "", err
	}

	if story.Title == "" {
		return "", fmt.Errorf("story has no title")
	}

	return story.Title, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
