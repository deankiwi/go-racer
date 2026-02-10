package game

import (
	"go-racer/pkg/config"
	"testing"
)

func TestTypingTest_Metrics(t *testing.T) {
	target := "Hello"
	game := NewTypingTest(target)
	game.Start()

	// Type 'H' correctly
	game.AddInput('H')
	// Type 'e' correctly
	game.AddInput('e')
	// Type 'x' instead of 'l' (mistake)
	game.AddInput('x')
	// Backspace
	game.Backspace()
	// Type 'l' correctly (correction)
	game.AddInput('l')
	// Type 'l' correctly
	game.AddInput('l')
	// Type 'o' correctly
	game.AddInput('o')

	game.Complete()

	stats := game.GetSessionStats()

	// 'H': 1 attempt, 0 mistakes
	if s, ok := stats["H"]; !ok || s.Attempts != 1 || s.Mistakes != 0 {
		t.Errorf("Stats for H incorrect: %+v", s)
	}

	// 'e': 1 attempt, 0 mistakes
	if s, ok := stats["e"]; !ok || s.Attempts != 1 || s.Mistakes != 0 {
		t.Errorf("Stats for e incorrect: %+v", s)
	}

	// 'l': 2 attempts.
	// First 'l' (index 2): typed 'x' first -> mistake.
	// Second 'l' (index 3): typed 'l' first -> correct.
	if s, ok := stats["l"]; !ok || s.Attempts != 2 || s.Mistakes != 1 {
		t.Errorf("Stats for l incorrect: %+v", s)
	}

	// 'o': 1 attempt, 0 mistakes
	if s, ok := stats["o"]; !ok || s.Attempts != 1 || s.Mistakes != 0 {
		t.Errorf("Stats for o incorrect: %+v", s)
	}
}

func TestApplyFilters(t *testing.T) {
	cfg := &config.Config{
		IncludeNumbers:          true,
		IncludePunctuation:      true,
		IncludeCapitalLetters:   true,
		IncludeNonStandardChars: true,
	}

	tests := []struct {
		name     string
		input    string
		cfgMod   func(*config.Config)
		expected string
	}{
		{
			name:     "Default (All Included)",
			input:    "Hello 123!",
			cfgMod:   func(c *config.Config) {},
			expected: "Hello 123!",
		},
		{
			name:     "No Numbers",
			input:    "Hello 123!",
			cfgMod:   func(c *config.Config) { c.IncludeNumbers = false },
			expected: "Hello !",
		},
		{
			name:     "No Punctuation",
			input:    "Hello, World!",
			cfgMod:   func(c *config.Config) { c.IncludePunctuation = false },
			expected: "Hello World",
		},
		{
			name:     "No Capital Letters",
			input:    "Hello World",
			cfgMod:   func(c *config.Config) { c.IncludeCapitalLetters = false },
			expected: "hello world",
		},
		{
			name:     "No Non-Standard Chars",
			input:    "Hello â",
			cfgMod:   func(c *config.Config) { c.IncludeNonStandardChars = false },
			expected: "Hello",
		},
		{
			name:  "Combined Filters",
			input: "Hello, 123 â!",
			cfgMod: func(c *config.Config) {
				c.IncludeNumbers = false
				c.IncludePunctuation = false
				c.IncludeCapitalLetters = false
				c.IncludeNonStandardChars = false
			},
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Copy config to avoid side effects
			testCfg := *cfg
			tt.cfgMod(&testCfg)

			got := ApplyFilters(tt.input, &testCfg)
			// ApplyFilters also cleans up spaces, so we expect trimmed output with single spaces
			if got != tt.expected {
				t.Errorf("ApplyFilters() = %q, want %q", got, tt.expected)
			}
		})
	}
}
