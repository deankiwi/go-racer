# Go Racer - Typing Speed Tester

A CLI tool written in Go that tests your typing speed using random stories from Hacker News.

## Features

- Fetches random stories from Hacker News API
- Real-time typing feedback with color-coded accuracy
- Calculates words per minute (WPM) and accuracy
- Beautiful terminal UI using Bubble Tea
- Tracks errors and correct characters

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone this repository
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

Build and run the application:

```bash
go run main.go
```

Or build it as a binary:

```bash
go build -o go-racer main.go
./go-racer
```

## How to Use

1. The application will automatically fetch a random story from Hacker News
2. Once loaded, you'll see the story title displayed
3. Start typing the text exactly as shown
4. Your typing will be highlighted in real-time:
   - Green: Correct characters
   - Red: Incorrect characters
   - Gray: Not yet typed
5. Press Enter when you're finished
6. View your results including WPM, accuracy, and time taken
7. Press 'q' to quit

## Controls

- Type normally to input characters
- Backspace to delete characters
- Enter to finish the test
- 'q' to quit (when viewing results)

## API

This tool uses the official Hacker News API:
- `https://hacker-news.firebaseio.com/v0/topstories.json` - Gets top story IDs
- `https://hacker-news.firebaseio.com/v0/item/{id}.json` - Gets individual story details

## Dependencies

- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/mattn/go-runewidth` - Character width calculation
- `github.com/rivo/uniseg` - Unicode segmentation

## License

MIT 