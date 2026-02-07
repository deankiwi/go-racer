package main

import (
	"flag"
	"fmt"
	"os"

	"go-racer/pkg/plugins"
	"go-racer/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	pluginName := flag.String("plugin", "hn", "Plugin source to use (hn, github)")
	flag.Parse()

	plugin, err := plugins.GetPlugin(*pluginName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Available plugins:", plugins.ListPlugins())
		os.Exit(1)
	}

	p := tea.NewProgram(ui.InitialModel(plugin))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
