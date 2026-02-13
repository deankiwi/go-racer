package ui

import (
	"fmt"
	"go-racer/pkg/config"
	"strings"
	"testing"
)

func TestRenderTrend(t *testing.T) {
	// Create a dummy config with history
	cfg := &config.Config{
		History: []config.GameResult{
			{WPM: 10, Accuracy: 100, Timestamp: 1},
			{WPM: 20, Accuracy: 100, Timestamp: 2},
			{WPM: 30, Accuracy: 100, Timestamp: 3},
			{WPM: 40, Accuracy: 100, Timestamp: 4},
			{WPM: 50, Accuracy: 100, Timestamp: 5},
		},
	}

	m := Model{
		Config: cfg,
		width:  100, // Sufficient width
	}

	// We can't access renderTrend directly if it is private, but in the same package (ui) we can.
	// Since the test file is package ui, we are good.
	output := m.renderTrend()

	// Verification
	if !strings.Contains(output, "WPM Trend") {
		t.Error("Output should contain title 'WPM Trend'")
	}

	// Check if graph content is roughly there
	if !strings.Contains(output, "Newest") {
		t.Error("Output should contain X axis label 'Newest'")
	}

	// Check for points
	if !strings.Contains(output, "•") {
		t.Error("Output should contain data points '•'")
	}

	fmt.Println(output)
}

func TestRenderTrend_Empty(t *testing.T) {
	cfg := &config.Config{
		History: []config.GameResult{},
	}

	m := Model{
		Config: cfg,
	}

	output := m.renderTrend()
	if !strings.Contains(output, "No games played yet") {
		t.Error("Should display 'No games played yet' for empty history")
	}
}

func TestRenderTrend_NotEnoughData(t *testing.T) {
	cfg := &config.Config{
		History: []config.GameResult{
			{WPM: 10, Accuracy: 100, Timestamp: 1},
		},
	}

	m := Model{
		Config: cfg,
	}

	output := m.renderTrend()
	if !strings.Contains(output, "Not enough data") {
		t.Error("Should display 'Not enough data' for single game history")
	}
}
