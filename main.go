package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HNStory represents a Hacker News story
type HNStory struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
}

// HNItem represents a single HN item
type HNItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
	Type  string `json:"type"`
}

// TypingTest represents the typing test state
type TypingTest struct {
	story       string
	userInput   string
	startTime   time.Time
	endTime     time.Time
	isComplete  bool
	errors      int
	totalChars  int
	correctChars int
}

// Model represents the application state
type Model struct {
	typingTest *TypingTest
	status     string
	err        error
}

// InitialModel returns the initial model
func InitialModel() Model {
	return Model{
		typingTest: &TypingTest{},
		status:     "Loading story...",
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return loadRandomStory()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.typingTest.isComplete {
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil
		}

		if m.typingTest.startTime.IsZero() {
			m.typingTest.startTime = time.Now()
		}

		switch msg.Type {
		case tea.KeyEnter:
			// End the test
			m.typingTest.endTime = time.Now()
			m.typingTest.isComplete = true
			m.calculateResults()
			return m, nil
		case tea.KeyBackspace:
			if len(m.typingTest.userInput) > 0 {
				m.typingTest.userInput = m.typingTest.userInput[:len(m.typingTest.userInput)-1]
			}
		case tea.KeyRunes:
			m.typingTest.userInput += string(msg.Runes)
		}
		return m, nil

	case storyLoadedMsg:
		m.typingTest.story = msg.story
		m.status = "Ready to start typing! Press any key to begin..."
		return m, nil

	case errorMsg:
		m.err = msg.err
		m.status = "Error loading story"
		return m, nil
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Render(fmt.Sprintf("Error: %v\nPress q to quit", m.err))
	}

	if m.typingTest.story == "" {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00")).
			Render(m.status)
	}

	if m.typingTest.isComplete {
		return m.renderResults()
	}

	return m.renderTypingInterface()
}

func (m Model) renderTypingInterface() string {
	story := m.typingTest.story
	userInput := m.typingTest.userInput

	// Create a visual representation of the typing
	var display strings.Builder
	display.WriteString(lipgloss.NewStyle().Bold(true).Render("Type the following story:\n\n"))

	// Show the story with user input highlighted
	for i, char := range story {
		if i < len(userInput) {
			if i < len(userInput) && userInput[i] == byte(char) {
				// Correct character
				display.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("#00ff00")).
					Render(string(char)))
			} else {
				// Incorrect character
				display.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("#ff0000")).
					Render(string(char)))
			}
		} else {
			// Not typed yet
			display.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render(string(char)))
		}
	}

	display.WriteString("\n\n")
	display.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffff00")).
		Render("Press Enter when finished\n"))

	return display.String()
}

func (m Model) renderResults() string {
	typingTest := m.typingTest
	duration := typingTest.endTime.Sub(typingTest.startTime)
	wpm := float64(typingTest.correctChars) / 5.0 / duration.Minutes()
	accuracy := float64(typingTest.correctChars) / float64(typingTest.totalChars) * 100

	var results strings.Builder
	results.WriteString(lipgloss.NewStyle().Bold(true).Render("Typing Test Results\n\n"))
	results.WriteString(fmt.Sprintf("Words per minute: %.1f\n", wpm))
	results.WriteString(fmt.Sprintf("Accuracy: %.1f%%\n", accuracy))
	results.WriteString(fmt.Sprintf("Time: %.2f seconds\n", duration.Seconds()))
	results.WriteString(fmt.Sprintf("Characters typed: %d\n", len(typingTest.userInput)))
	results.WriteString(fmt.Sprintf("Correct characters: %d\n", typingTest.correctChars))
	results.WriteString(fmt.Sprintf("Errors: %d\n", typingTest.errors))
	results.WriteString("\nPress 'q' to quit\n")

	return results.String()
}

func (m *Model) calculateResults() {
	typingTest := m.typingTest
	story := typingTest.story
	userInput := typingTest.userInput

	typingTest.totalChars = len(story)
	typingTest.correctChars = 0
	typingTest.errors = 0

	for i, char := range story {
		if i < len(userInput) {
			if userInput[i] == byte(char) {
				typingTest.correctChars++
			} else {
				typingTest.errors++
			}
		} else {
			typingTest.errors++
		}
	}

	// Add errors for extra characters typed
	if len(userInput) > len(story) {
		typingTest.errors += len(userInput) - len(story)
	}
}

// Messages
type storyLoadedMsg struct {
	story string
}

type errorMsg struct {
	err error
}

// Commands
func loadRandomStory() tea.Cmd {
	return func() tea.Msg {
		story, err := fetchRandomStory()
		if err != nil {
			return errorMsg{err: err}
		}
		return storyLoadedMsg{story: story}
	}
}

func fetchRandomStory() (string, error) {
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

	var story HNItem
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return "", err
	}

	// Use title as the text to type
	if story.Title == "" {
		return "", fmt.Errorf("no title found for story")
	}

	return story.Title, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
} 