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
		log.Fatalf("Failed to ensure yt-dlp is available: %v\nPlease ensure yt-dlp is installed or the binary contains embedded dependencies.", err)
	}
	if err := bootstrap.Ensure("ffmpeg"); err != nil {
		log.Fatalf("Failed to ensure ffmpeg is available: %v\nPlease ensure ffmpeg is installed or the binary contains embedded dependencies.", err)
	}

	if _, err := bootstrap.GetCommandPath("yt-dlp"); err != nil {
		log.Fatalf("yt-dlp not found after installation attempt: %v", err)
	}
	if _, err := bootstrap.GetCommandPath("ffmpeg"); err != nil {
		log.Fatalf("ffmpeg not found after installation attempt: %v", err)
	}

	p := tea.NewProgram(
		ui.NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
