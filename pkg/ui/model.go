package ui

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
	uni "unicode"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/reflow/wordwrap"

	"go-racer/pkg/config"
	"go-racer/pkg/game"
	"go-racer/pkg/plugins"
)

type Model struct {
	Game              *game.TypingTest
	Plugin            plugins.ContentSource
	CurrentPluginName string
	Err               error
	IsLoading         bool
	Spinner           spinner.Model
	Quitting          bool
	Config            *config.Config
	ShowMetrics       bool
	ShowSettings      bool
	ShowTrend         bool
	CurrentContent    *plugins.Content
	width             int
	height            int
}

func InitialModel(plugin plugins.ContentSource, pluginName string, cfg *config.Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		Plugin:            plugin,
		CurrentPluginName: pluginName,
		IsLoading:         true,
		Spinner:           s,
		Config:            cfg,
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.Quitting = true
			return m, tea.Quit
		}

		if m.IsLoading {
			return m, nil
		}

		if m.Game.IsComplete {
			if m.ShowSettings {
				switch msg.String() {
				case "esc", ",":
					m.ShowSettings = false
				case "n":
					m.Config.IncludeNumbers = !m.Config.IncludeNumbers
					_ = config.Save(m.Config)
				case "p":
					m.Config.IncludePunctuation = !m.Config.IncludePunctuation
					_ = config.Save(m.Config)
				case "c":
					m.Config.IncludeCapitalLetters = !m.Config.IncludeCapitalLetters
					_ = config.Save(m.Config)
				case "s":
					m.Config.IncludeNonStandardChars = !m.Config.IncludeNonStandardChars
					_ = config.Save(m.Config)
				}
				return m, nil
			}

			if msg.String() == "q" || msg.Type == tea.KeyEsc {
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
			if msg.String() == "," {
				m.ShowSettings = !m.ShowSettings
				return m, nil
			}
			if msg.String() == "m" {
				m.ShowMetrics = !m.ShowMetrics
				return m, nil
			}

			if m.ShowMetrics {
				if msg.String() == "esc" {
					m.ShowMetrics = false
					return m, nil
				}
				return m, nil
			}

			if msg.String() == "t" {
				m.ShowTrend = !m.ShowTrend
				return m, nil
			}

			if m.ShowTrend {
				if msg.String() == "esc" {
					m.ShowTrend = false
					return m, nil
				}
				return m, nil
			}

			if msg.String() == "p" {
				// Switch plugin
				nextPlugin := "github"
				if m.CurrentPluginName == "github" {
					nextPlugin = "hn"
				}

				p, err := plugins.GetPlugin(nextPlugin)
				if err != nil {
					m.Err = err
					return m, nil
				}

				m.Plugin = p
				m.CurrentPluginName = nextPlugin
				m.IsLoading = true

				// Save config
				m.Config.LastPlugin = nextPlugin
				_ = config.Save(m.Config)

				return m, tea.Batch(
					m.Spinner.Tick,
					m.loadContent,
				)
			}

			if msg.Type == tea.KeyEnter {
				if m.CurrentContent != nil && m.CurrentContent.SourceURL != "" {
					// Open URL
					return m, tea.ExecProcess(exec.Command("open", m.CurrentContent.SourceURL), func(err error) tea.Msg {
						if err != nil {
							return errorMsg{err}
						}
						return nil
					})
				}
			}

			return m, nil
		}

		// Game logic input handling
		switch msg.Type {
		case tea.KeyEsc:
			m.Game.Complete()
			m.Game.CalculateStats()
			m.saveMetrics()
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
			m.saveMetrics()
		}

	case contentMsg:
		m.IsLoading = false
		m.Game = game.NewTypingTest(msg.content.Text)
		m.CurrentContent = msg.content
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
		if m.ShowSettings {
			return m.renderSettings()
		}
		if m.ShowMetrics {
			return m.renderMetrics()
		}
		if m.ShowTrend {
			return m.renderTrend()
		}
		return m.renderResults()
	}

	return m.renderGame()
}

func (m Model) renderGame() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render("Go Racer - " + m.Plugin.Name()))
	s.WriteString("\n\n")

	// Render text with highlighting
	var textBuilder strings.Builder
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

		textBuilder.WriteString(style.Render(string(char)))
	}

	// Apply word wrap
	width := m.width - 4 // Account for some padding
	if width < 20 {
		// Fallback if width is not yet set or too small
		if m.width == 0 {
			width = 60 // Default width if unknown
		} else {
			width = 20
		}
	}
	s.WriteString(wordwrap.String(textBuilder.String(), width))

	s.WriteString("\n\n")
	s.WriteString(UntypedStyle.Render("Start typing... Press Esc to finish, Ctrl+C to quit"))

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
	var textBuilder strings.Builder
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
		textBuilder.WriteString(style.Render(string(char)))
	}

	// Apply word wrap
	width := m.width - 8 // Account for border (2) + padding (4) + extra safety (2)
	if width < 20 {
		if m.width == 0 {
			width = 60
		} else {
			width = 20
		}
	}
	s.WriteString(wordwrap.String(textBuilder.String(), width))
	s.WriteString("\n\n")

	content := fmt.Sprintf(
		"WPM:      %.2f\n"+
			"Accuracy: %.2f%%\n"+
			"Time:     %.2fs\n\n"+
			"Press 'r' to retry, 'q' to quit\n"+
			"Press 'm' to view metrics\n"+
			"Press ',' for settings\n"+
			"Press 'p' to switch plugin (Current: %s)\n"+
			"Press 't' to view trend",
		wpm, accuracy, duration.Seconds(), m.Plugin.Name(),
	)

	if m.CurrentContent != nil && m.CurrentContent.SourceURL != "" {
		content += "\nPress 'Enter' to open source"
	}

	s.WriteString(content)

	return ResultsStyle.Render(s.String())
}

func (m Model) renderSettings() string {
	var s strings.Builder
	s.WriteString(ResultsStyle.Render("Settings"))
	s.WriteString("\n\n")

	checkbox := func(label string, checked bool, key string) string {
		check := "[ ]"
		if checked {
			check = "[x]"
		}
		return fmt.Sprintf("%s %-25s (%s)\n", check, label, key)
	}

	s.WriteString(checkbox("Include Numbers", m.Config.IncludeNumbers, "n"))
	s.WriteString(checkbox("Include Punctuation", m.Config.IncludePunctuation, "p"))
	s.WriteString(checkbox("Include Capital Letters", m.Config.IncludeCapitalLetters, "c"))
	s.WriteString(checkbox("Include Non-Standard", m.Config.IncludeNonStandardChars, "s"))

	s.WriteString("\nPress ',' or 'Esc' to return\n")

	return ResultsStyle.Render(s.String())
}

// Messages
type contentMsg struct {
	content *plugins.Content
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
	content.Text = game.ApplyFilters(content.Text, m.Config)
	return contentMsg{content}
}

func (m *Model) saveMetrics() {
	if m.Config.Metrics == nil {
		m.Config.Metrics = make(map[string]config.CharMetric)
	}

	sessionStats := m.Game.GetSessionStats()
	for char, stat := range sessionStats {
		existing := m.Config.Metrics[char]
		existing.Attempts += stat.Attempts
		existing.Mistakes += stat.Mistakes
		existing.Mistakes += stat.Mistakes
		m.Config.Metrics[char] = existing
	}

	// Save history
	result := config.GameResult{
		WPM:       m.Game.WPM(),
		Accuracy:  m.Game.Accuracy(),
		Timestamp: time.Now().Unix(),
	}
	m.Config.History = append(m.Config.History, result)

	_ = config.Save(m.Config)
}

func (m Model) renderMetrics() string {
	var s strings.Builder
	s.WriteString(ResultsStyle.Render("Character Metrics (Worst Accuracy First)"))
	s.WriteString("\n\n")

	type charStat struct {
		Char     string
		Attempts int
		Mistakes int
		Accuracy float64
	}

	// Map to aggregate stats by lowercase character
	aggregatedStats := make(map[string]struct{ Attempts, Mistakes int })

	for char, metric := range m.Config.Metrics {
		runes := []rune(char)
		if len(runes) != 1 {
			continue
		}
		r := runes[0]

		// Filter: Only allow ASCII letters (A-Z, a-z)
		// if r > uni.MaxASCII || !uni.IsLetter(r) {
		// 	continue
		// }

		lowerChar := string(uni.ToLower(r))
		s := aggregatedStats[lowerChar]
		s.Attempts += metric.Attempts
		s.Mistakes += metric.Mistakes
		aggregatedStats[lowerChar] = s
	}

	var stats []charStat
	for char, s := range aggregatedStats {
		accuracy := 0.0
		if s.Attempts > 0 {
			accuracy = (float64(s.Attempts-s.Mistakes) / float64(s.Attempts)) * 100
		}
		stats = append(stats, charStat{
			Char:     strings.ToUpper(char), // Display as Uppercase
			Attempts: s.Attempts,
			Mistakes: s.Mistakes,
			Accuracy: accuracy,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Accuracy != stats[j].Accuracy {
			return stats[i].Accuracy < stats[j].Accuracy
		}
		if stats[i].Attempts != stats[j].Attempts {
			return stats[i].Attempts > stats[j].Attempts
		}
		return stats[i].Char < stats[j].Char
	})

	// Header
	s.WriteString(fmt.Sprintf("%-5s | %-10s | %-10s | %s\n", "Char", "Accuracy", "Mistakes", "Attempts"))
	s.WriteString(strings.Repeat("-", 45) + "\n")

	// Limit to top 20 or fit screen? Let's show top 15 for now
	count := 0
	for _, stat := range stats {
		if count >= 15 {
			break
		}
		displayChar := stat.Char
		if displayChar == " " {
			displayChar = "SPC"
		}
		s.WriteString(fmt.Sprintf("%-5s | %-9.1f%% | %-10d | %d\n", displayChar, stat.Accuracy, stat.Mistakes, stat.Attempts))
		count++
	}

	s.WriteString("\nPress 'm' or 'Esc' to return\n")

	return ResultsStyle.Render(s.String())
}

func (m Model) renderTrend() string {
	var s strings.Builder
	s.WriteString(ResultsStyle.Render("WPM Trend (Last 20 Games)"))
	s.WriteString("\n\n")

	history := m.Config.History
	if len(history) == 0 {
		s.WriteString("No games played yet.")
		return ResultsStyle.Render(s.String())
	}

	// Limit to last 20 games
	startIdx := 0
	if len(history) > 20 {
		startIdx = len(history) - 20
	}
	recentHistory := history[startIdx:]

	if len(recentHistory) < 2 {
		s.WriteString(fmt.Sprintf("Not enough data to show trend (Played: %d)\n", len(history)))
		s.WriteString("\nPress 't' or 'Esc' to return\n")
		return ResultsStyle.Render(s.String())
	}

	// Find min and max WPM for scaling
	minWPM := recentHistory[0].WPM
	maxWPM := recentHistory[0].WPM
	for _, res := range recentHistory {
		if res.WPM < minWPM {
			minWPM = res.WPM
		}
		if res.WPM > maxWPM {
			maxWPM = res.WPM
		}
	}

	// Add some padding to Y axis
	minWPM = minWPM * 0.9
	maxWPM = maxWPM * 1.1
	if minWPM < 0 {
		minWPM = 0
	}

	graphHeight := 10
	graphWidth := 60
	if m.width > 20 {
		graphWidth = m.width - 20
	}

	// Create grid
	grid := make([][]string, graphHeight)
	for i := range grid {
		grid[i] = make([]string, graphWidth)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Plot points
	for i, res := range recentHistory {
		// Map x (index) to grid width
		// spread points evenly
		x := int(float64(i) / float64(len(recentHistory)-1) * float64(graphWidth-1))

		// Map y (WPM) to grid height
		// Y is inverted in grid (0 is top)
		normalizedWPM := (res.WPM - minWPM) / (maxWPM - minWPM)
		y := int((1.0 - normalizedWPM) * float64(graphHeight-1))

		if y >= 0 && y < graphHeight && x >= 0 && x < graphWidth {
			grid[y][x] = "•"
		}
	}

	// Render grid with Y axis labels
	for i, row := range grid {
		// Calculate WPM for this row
		rowWPM := maxWPM - (float64(i)/float64(graphHeight-1))*(maxWPM-minWPM)
		s.WriteString(fmt.Sprintf("%6.1f |", rowWPM))
		for _, cell := range row {
			s.WriteString(cell)
		}
		s.WriteString("\n")
	}

	// X axis
	s.WriteString("       " + strings.Repeat("-", graphWidth) + "\n")
	s.WriteString(fmt.Sprintf("       %-*s%s\n", graphWidth/2, "Oldest", "Newest"))

	s.WriteString("\nPress 't' or 'Esc' to return\n")

	return ResultsStyle.Render(s.String())
}
