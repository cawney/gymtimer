package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gymtimer/internal/audio"
	"gymtimer/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Determine assets path
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	assetsDir := filepath.Join(filepath.Dir(execPath), "assets")

	// Also check current directory
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		assetsDir = "assets"
	}

	beepPath := filepath.Join(assetsDir, "beep.wav")

	// Generate beep sound if it doesn't exist
	if _, err := os.Stat(beepPath); os.IsNotExist(err) {
		if err := os.MkdirAll(assetsDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create assets directory: %v\n", err)
		}
		if err := audio.GenerateBeepWAV(beepPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not generate beep sound: %v\n", err)
		}
	}

	// Create audio player
	audioPlayer := audio.New(beepPath)

	// Create the app model
	model := ui.New(audioPlayer)

	// Create and run the Bubbletea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
