package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".go-racer.json"

type CharMetric struct {
	Attempts int `json:"attempts"`
	Mistakes int `json:"mistakes"`
}

type GameResult struct {
	WPM       float64 `json:"wpm"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp int64   `json:"timestamp"`
}

type Config struct {
	LastPlugin              string                `json:"last_plugin"`
	Metrics                 map[string]CharMetric `json:"metrics"`
	History                 []GameResult          `json:"history"`
	IncludeNumbers          bool                  `json:"include_numbers"`
	IncludePunctuation      bool                  `json:"include_punctuation"`
	IncludeCapitalLetters   bool                  `json:"include_capital_letters"`
	IncludeNonStandardChars bool                  `json:"include_non_standard_chars"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{
			LastPlugin:              "hn",
			IncludeNumbers:          true,
			IncludePunctuation:      true,
			IncludeCapitalLetters:   true,
			IncludeNonStandardChars: true,
		}, nil // Default
	}
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
