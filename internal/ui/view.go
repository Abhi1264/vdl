package ui

import (
	"fmt"
	"strings"

	"github.com/Abhi1264/vidforge/internal/downloader"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	profileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	activeProfileStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	doneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	helpStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 2).
			Margin(1)
)

func renderBar(p float64) string {
	width := 30
	filled := int(p * float64(width))
	if filled > width {
		filled = width
	}
	return barFilled.Render(strings.Repeat("█", filled)) +
		barEmpty.Render(strings.Repeat("░", width-filled))
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("VidForge — Universal Video Downloader\n\n"))

	if m.showHelp {
		b.WriteString(helpStyle.Render(helpText))
		b.WriteString("\n\n")
		return b.String()
	}

	if m.showProfiles {
		return m.renderProfiles()
	}

	if m.showDownloadPath {
		return m.renderDownloadPath()
	}

	b.WriteString("URL:\n")
	b.WriteString(m.input.View() + "\n\n")

	if m.profile != nil {
		profileText := fmt.Sprintf("Profile: %s (%s)", m.profile.Name, m.profile.Quality)
		b.WriteString(activeProfileStyle.Render(profileText))
		b.WriteString("\n")
		b.WriteString(profileStyle.Render(m.profile.Description))
		b.WriteString("\n\n")
	}

	sponsorText := "SponsorBlock: "
	if m.sponsorBlock {
		sponsorText += doneStyle.Render("ON")
	} else {
		sponsorText += statusStyle.Render("OFF")
	}
	b.WriteString(sponsorText)
	b.WriteString("\n\n")

	downloadPathText := fmt.Sprintf("Download Location: %s", m.downloadPath)
	b.WriteString(profileStyle.Render(downloadPathText))
	b.WriteString("\n\n")

	if len(m.jobs) > 0 {
		b.WriteString("Downloads:\n")
		jobIDs := m.getSortedIDs()

		for _, id := range jobIDs {
			job := m.jobs[id]
			isSelected := id == m.selectedID

			var line strings.Builder
			if isSelected {
				line.WriteString(selectedStyle.Render("▶"))
				line.WriteString(" ")
			} else {
				line.WriteString("  ")
			}

			bar := renderBar(job.progress)
			line.WriteString(bar)
			line.WriteString(fmt.Sprintf("  %.0f%%  ", job.progress*100))

			var statusStr string
			switch {
			case job.done && job.status == "done":
				statusStr = doneStyle.Render("✓ done")
			case strings.HasPrefix(job.status, "error"):
				statusStr = errorStyle.Render("✗ error")
			case job.status == "cancelled":
				statusStr = statusStyle.Render("⊘ cancelled")
			default:
				statusStr = statusStyle.Render(job.status)
			}
			line.WriteString(statusStr)
			line.WriteString("  ")

			url := job.url
			if len(url) > 50 {
				url = url[:47] + "..."
			}
			line.WriteString(url)

			if isSelected {
				b.WriteString(selectedStyle.Render(line.String()))
			} else {
				b.WriteString(line.String())
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("\n[Enter] download  [↑↓] select  [p] pause/cancel  [s] SponsorBlock  [f] profile  [d] download location  [?] help  [q] quit\n")
	return b.String()
}

func (m Model) renderProfiles() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Select Download Profile\n\n"))

	profiles := downloader.GetProfiles()
	for i, p := range profiles {
		var style lipgloss.Style
		if i == m.profileIndex {
			style = activeProfileStyle
		} else {
			style = profileStyle
		}

		line := fmt.Sprintf("[%d] %s\n", i+1, p.Description)
		b.WriteString(style.Render(line))
	}

	b.WriteString("\n[1-6] select  [↑↓] navigate  [Enter] confirm  [f] back  [q] quit\n")
	return b.String()
}

func (m Model) renderDownloadPath() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Set Download Location\n\n"))
	b.WriteString("Enter the path where downloads should be saved:\n")
	b.WriteString(m.downloadPathInput.View() + "\n\n")
	b.WriteString(profileStyle.Render(fmt.Sprintf("Current: %s", m.downloadPath)))
	b.WriteString("\n\n")
	b.WriteString("\n[Enter] save  [d] back  [q] quit\n")
	return b.String()
}
