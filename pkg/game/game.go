package game

import (
	"time"
)

// TypingTest represents the state of a typing session
type TypingTest struct {
	TargetText     string
	UserInput      string
	StartTime      time.Time
	EndTime        time.Time
	IsComplete     bool
	IsStarted      bool
	Errors         int
	CorrectChars   int
	InitialMistake map[int]bool // Tracks indices where the first attempt was incorrect
}

// NewTypingTest creates a new typing test with the given target text
func NewTypingTest(text string) *TypingTest {
	return &TypingTest{
		TargetText:     text,
		InitialMistake: make(map[int]bool),
	}
}

// Start begins the typing test if not already started
func (t *TypingTest) Start() {
	if !t.IsStarted {
		t.StartTime = time.Now()
		t.IsStarted = true
	}
}

// AddInput appends a character to the user input
func (t *TypingTest) AddInput(r rune) {
	if t.IsComplete {
		return
	}
	t.Start()

	index := len(t.UserInput)
	// Track initial mistake if this is the first attempt at this index
	if _, attempted := t.InitialMistake[index]; !attempted && index < len(t.TargetText) {
		if byte(r) != t.TargetText[index] {
			t.InitialMistake[index] = true
		} else {
			// Mark as attempted but correct (false)
			t.InitialMistake[index] = false
		}
	}

	t.UserInput += string(r)

	// Check for completion
	if len(t.UserInput) >= len(t.TargetText) {
		t.Complete()
	}
}

// Backspace removes the last character from user input
func (t *TypingTest) Backspace() {
	if t.IsComplete || len(t.UserInput) == 0 {
		return
	}
	t.UserInput = t.UserInput[:len(t.UserInput)-1]
}

// BackspaceWord removes the last word from user input
func (t *TypingTest) BackspaceWord() {
	if t.IsComplete || len(t.UserInput) == 0 {
		return
	}

	// Convert to runes for safe handling
	runes := []rune(t.UserInput)
	if len(runes) == 0 {
		return
	}

	// 1. Remove trailing spaces
	for len(runes) > 0 && runes[len(runes)-1] == ' ' {
		runes = runes[:len(runes)-1]
	}

	// 2. Remove characters until space or start
	for len(runes) > 0 && runes[len(runes)-1] != ' ' {
		runes = runes[:len(runes)-1]
	}

	t.UserInput = string(runes)
}

// Complete finishes the test and calculates final stats
func (t *TypingTest) Complete() {
	if !t.IsStarted {
		return
	}
	t.EndTime = time.Now()
	t.IsComplete = true
	t.CalculateStats()
}

// CalculateStats updates the error and correct character counts
func (t *TypingTest) CalculateStats() {
	t.CorrectChars = 0
	t.Errors = 0

	for i, char := range t.TargetText {
		if i < len(t.UserInput) {
			if t.UserInput[i] == byte(char) {
				t.CorrectChars++
			} else {
				t.Errors++
			}
		} else {
			// Count untyped characters as errors only if test is complete?
			// For now, let's just count typed errors.
		}
	}

	// key point: errors should also account for extra characters typed if any (though we capped it above)
	if len(t.UserInput) > len(t.TargetText) {
		t.Errors += len(t.UserInput) - len(t.TargetText)
	}
}

// WPM calculates words per minute
func (t *TypingTest) WPM() float64 {
	var duration time.Duration
	if t.IsComplete {
		duration = t.EndTime.Sub(t.StartTime)
	} else if t.IsStarted {
		duration = time.Since(t.StartTime)
	} else {
		return 0
	}

	if duration.Minutes() == 0 {
		return 0
	}

	// Standard WPM calculation: (characters / 5) / minutes
	return (float64(len(t.UserInput)) / 5.0) / duration.Minutes()
}

// Accuracy calculates the percentage of characters correct on the first try
func (t *TypingTest) Accuracy() float64 {
	if len(t.UserInput) == 0 {
		return 100
	}

	totalMistakes := 0
	for _, mistake := range t.InitialMistake {
		if mistake {
			totalMistakes++
		}
	}

	// Based on total characters typed so far (up to length of target) or just total target length?
	// User said "got it correct first time". Usually means compared to total text.
	// Let's use the length of the text we've attempted so far.
	// Actually, if we've typed N chars, we have N entries in InitialMistake (true or false).
	// So len(InitialMistake) is the number of attempted indices.

	if len(t.InitialMistake) == 0 {
		return 100
	}

	return float64(len(t.InitialMistake)-totalMistakes) / float64(len(t.InitialMistake)) * 100
}

// GetSessionStats calculates character-level metrics for the current session
func (t *TypingTest) GetSessionStats() map[string]struct{ Attempts, Mistakes int } {
	stats := make(map[string]struct{ Attempts, Mistakes int })

	for i, mistyped := range t.InitialMistake {
		if i >= len(t.TargetText) {
			continue
		}
		char := string(t.TargetText[i])
		s := stats[char]
		s.Attempts++
		if mistyped {
			s.Mistakes++
		}
		stats[char] = s
	}
	return stats
}
