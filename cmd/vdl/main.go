package main

import (
	"github.com/Abhi1264/vdl/internal/bootstrap"
	"github.com/Abhi1264/vdl/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"log"
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
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
