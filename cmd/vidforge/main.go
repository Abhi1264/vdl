package main

import (
	"log"

	"github.com/Abhi1264/vidforge/internal/bootstrap"
	"github.com/Abhi1264/vidforge/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	if err := bootstrap.Ensure("yt-dlp"); err != nil {
		log.Printf("Warning: yt-dlp installation failed: %v", err)
	}
	if err := bootstrap.Ensure("ffmpeg"); err != nil {
		log.Printf("Warning: ffmpeg installation failed: %v", err)
	}
}

func main() {
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
