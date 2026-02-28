package plugins

import (
	"math/rand"
	"time"
)

type RandomSource struct{}

func NewRandomSource() *RandomSource {
	return &RandomSource{}
}

func (r *RandomSource) Name() string {
	return "Random"
}

func (r *RandomSource) Description() string {
	return "Selects a random plugin to get content from"
}

func (r *RandomSource) GetContent() (*Content, error) {
	// Available plugins, intentionally omitting "random" to avoid infinite recursion
	available := []string{"hn", "github", "spanish-news"}

	rand.Seed(time.Now().UnixNano())
	pluginName := available[rand.Intn(len(available))]

	source, err := GetPlugin(pluginName)
	if err != nil {
		return nil, err
	}

	return source.GetContent()
}
