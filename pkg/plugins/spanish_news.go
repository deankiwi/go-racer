package plugins

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type SpanishNewsSource struct{}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

func NewSpanishNewsSource() *SpanishNewsSource {
	return &SpanishNewsSource{}
}

func (s *SpanishNewsSource) Name() string {
	return "Spanish News"
}

func (s *SpanishNewsSource) Description() string {
	return "Headlines from El País (Spanish)"
}

func (s *SpanishNewsSource) GetContent() (*Content, error) {
	resp, err := http.Get("https://elpais.com/rss/elpais/portada.xml")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()

	var rss RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("failed to decode feed: %w", err)
	}

	if len(rss.Channel.Items) == 0 {
		return nil, fmt.Errorf("no stories found")
	}

	rand.Seed(time.Now().UnixNano())
	item := rss.Channel.Items[rand.Intn(len(rss.Channel.Items))]

	// Clean up title (sometimes it has CDATA or other artifacts, but XML decoder handles CDATA)
	// We might want to remove " - EL PAÍS" suffix if it exists, or similar cleanup.
	// For now, raw title is probably fine.

	return &Content{
		Text:      item.Title,
		SourceURL: item.Link,
		Author:    "El País",
	}, nil
}
