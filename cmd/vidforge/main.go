package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Abhi1264/vidforge/internal/bootstrap"
	"github.com/Abhi1264/vidforge/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if os.Getenv("CI") == "true" {
		fmt.Println("VidForge CI mode: startup OK")
		return
	}

	if err := bootstrap.Ensure("yt-dlp"); err != nil {
		log.Printf("Warning: yt-dlp installation failed: %v", err)
	}
	if err := bootstrap.Ensure("ffmpeg"); err != nil {
		log.Printf("Warning: ffmpeg installation failed: %v", err)
	}

	p := tea.NewProgram(
		ui.NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
