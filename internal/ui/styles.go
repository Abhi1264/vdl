package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	barFilled = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	barEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
