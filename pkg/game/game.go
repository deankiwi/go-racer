package game

import (
	"time"
)

// TypingTest represents the state of a typing session
type TypingTest struct {
	TargetText   string
	UserInput    string
	StartTime    time.Time
	EndTime      time.Time
	IsComplete   bool
	IsStarted    bool
	Errors       int
	CorrectChars int
}

// NewTypingTest creates a new typing test with the given target text
func NewTypingTest(text string) *TypingTest {
	return &TypingTest{
		TargetText: text,
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

// Accuracy calculates the percentage of correct characters
func (t *TypingTest) Accuracy() float64 {
	if len(t.UserInput) == 0 {
		return 100
	}
	return float64(t.CorrectChars) / float64(len(t.UserInput)) * 100
}
