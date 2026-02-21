# AI Agents Guide

This document outlines the usage and context for AI agents (LLMs) working on this repository (`go-racer`). 

> **IMPORTANT:** This document should be kept up to date as the project evolves, new practices are adopted, or new agentic workflows are introduced.

## Project Structure

- `cmd/go-racer/`: The main entry point for the CLI application.
- `pkg/`: Contains the core logic organized by domain:
  - `config/`: Handles loading and saving user configurations and metrics.
  - `game/`: The core typing test mechanics and state.
  - `plugins/`: Different text sources for the typing test (e.g., Hacker News, Wikipedia).
  - `ui/`: Terminal UI components built with bubbletea and lipgloss.
- `doc/`: Documentation files.
- `AGENTS.md`: This guide.
- `TODO.md`: Tracks upcoming features, bug fixes, and general tasks.

## Running the Application

To run the project locally from the root of the repository during development:
```bash
go run cmd/go-racer/main.go
```

To build and install the application globally on your system:
```bash
go install ./cmd/go-racer
```
This will allow you to run the `go-racer` command from anywhere, provided your `$GOPATH/bin` is in your system's PATH.

## User Configuration

The user configuration file is stored in the home directory and named `.go-racer.json`.
- Path expansion: `~/.go-racer.json`
- It is loaded via `pkg/config/config.go`.
- This file stores the user's plugin preferences, typing history, and feature toggles (e.g., whether to include numbers or punctuation).

---
*Note: Any time the project's development workflow changes significantly, please update this file to reflect the new state so future agents have accurate context.*
