package main

import (
	"flag"
	"fmt"
	"os"

	"go-racer/pkg/config"
	"go-racer/pkg/plugins"
	"go-racer/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		// Ignore error, use default
		cfg = &config.Config{LastPlugin: "hn"}
	}

	pluginName := flag.String("plugin", cfg.LastPlugin, "Plugin source to use (hn, github)")
	flag.Parse()

	// Update config with the selected plugin (whether from flag or default)
	if *pluginName != cfg.LastPlugin {
		cfg.LastPlugin = *pluginName
		_ = config.Save(cfg)
	}

	plugin, err := plugins.GetPlugin(*pluginName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Available plugins:", plugins.ListPlugins())
		os.Exit(1)
	}

	p := tea.NewProgram(ui.InitialModel(plugin, *pluginName, cfg))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
