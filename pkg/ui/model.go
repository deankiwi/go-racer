package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"go-racer/pkg/game"
	"go-racer/pkg/plugins"
)

type Model struct {
	Game      *game.TypingTest
	Plugin    plugins.ContentSource
	Err       error
	IsLoading bool
	Spinner   spinner.Model
	Quitting  bool
}

func InitialModel(plugin plugins.ContentSource) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		Plugin:    plugin,
		IsLoading: true,
		Spinner:   s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.loadContent,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.Quitting = true
			return m, tea.Quit
		}

		if m.IsLoading {
			return m, nil
		}

		if m.Game.IsComplete {
			if msg.String() == "q" {
				m.Quitting = true
				return m, tea.Quit
			}
			if msg.String() == "r" {
				m.IsLoading = true
				return m, tea.Batch(
					m.Spinner.Tick,
					m.loadContent,
				)
			}
			return m, nil
		}

		// Game logic input handling
		switch msg.Type {
		case tea.KeyBackspace:
			if msg.Alt {
				m.Game.BackspaceWord()
			} else {
				m.Game.Backspace()
			}
		case tea.KeyCtrlW:
			m.Game.BackspaceWord()
		case tea.KeyRunes:
			// Handle space as a rune
			m.Game.AddInput(msg.Runes[0])
		case tea.KeySpace:
			m.Game.AddInput(' ')
		}

		// Verify completion after input
		if len(m.Game.UserInput) >= len(m.Game.TargetText) {
			m.Game.Complete()
			// Stats need to be calculated one last time to be sure
			m.Game.CalculateStats()
		}

	case contentMsg:
		m.IsLoading = false
		m.Game = game.NewTypingTest(msg.content)
		m.Game.Start() // Start timer immediately on load? Or wait for first keypress?
		// Let's modify game to start on first input in a future iteration if needed.
		// For now, let's just start a timer but only count "active" time?
		// Actually, Game.Start() sets StartTime. We should probably reset StartTime on first input.
		// But for simplicity, let's just let it be.
		m.Game.StartTime = time.Now()
		return m, nil

	case errorMsg:
		m.Err = msg.err
		m.IsLoading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit", m.Err)
	}

	if m.IsLoading {
		return fmt.Sprintf("\n %s Loading content from %s...\n\n", m.Spinner.View(), m.Plugin.Name())
	}

	if m.Game == nil {
		return "Initializing..."
	}

	if m.Game.IsComplete {
		return m.renderResults()
	}

	return m.renderGame()
}

func (m Model) renderGame() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render("Go Racer - " + m.Plugin.Name()))
	s.WriteString("\n\n")

	// Render text with highlighting
	for i, char := range m.Game.TargetText {
		var style lipgloss.Style

		if i < len(m.Game.UserInput) {
			if byte(char) == m.Game.UserInput[i] {
				style = CorrectStyle
			} else {
				style = ErrorStyle
			}
		} else {
			style = UntypedStyle
		}

		// Underline current character
		if i == len(m.Game.UserInput) {
			style = style.Copy().Underline(true)
		}

		s.WriteString(style.Render(string(char)))
	}

	s.WriteString("\n\n")
	s.WriteString(UntypedStyle.Render("Start typing... Press Ctrl+C to quit"))

	return s.String()
}

func (m Model) renderResults() string {
	duration := m.Game.EndTime.Sub(m.Game.StartTime)
	wpm := (float64(len(m.Game.UserInput)) / 5.0) / duration.Minutes()
	accuracy := m.Game.Accuracy()

	var s strings.Builder
	s.WriteString(ResultsStyle.Render("Results"))
	s.WriteString("\n\n")

	// Render the text with historical accuracy colors
	for i, char := range m.Game.TargetText {
		var style lipgloss.Style
		if mistyped, attempted := m.Game.InitialMistake[i]; attempted {
			if mistyped {
				style = ErrorStyle
			} else {
				style = CorrectStyle
			}
		} else {
			style = UntypedStyle
		}
		s.WriteString(style.Render(string(char)))
	}
	s.WriteString("\n\n")

	content := fmt.Sprintf(
		"WPM:      %.2f\n"+
			"Accuracy: %.2f%%\n"+
			"Time:     %.2fs\n\n"+
			"Press 'r' to retry, 'q' to quit",
		wpm, accuracy, duration.Seconds(),
	)

	s.WriteString(content)

	return ResultsStyle.Render(s.String())
}

// Messages
type contentMsg struct {
	content string
}

type errorMsg struct {
	err error
}

// Commands
func (m Model) loadContent() tea.Msg {
	content, err := m.Plugin.GetContent()
	if err != nil {
		return errorMsg{err}
	}
	return contentMsg{content}
}
