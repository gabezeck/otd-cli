package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F55385"}
	text      = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"}

	// Header
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(highlight).
		Padding(1, 2).
		MarginBottom(1)

	headerTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	headerDateStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#F55385")).
		Padding(0, 1)

	// Sections
	sectionHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(special).
		MarginTop(1).
		MarginBottom(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(subtle)

	// Timeline / Events
	yearStyle = lipgloss.NewStyle().
		Foreground(warning).
		Bold(true).
		Width(6).
		Align(lipgloss.Right).
		MarginRight(1)

	timelineDotStyle = lipgloss.NewStyle().
		Foreground(subtle).
		SetString("│") // Vertical line for timeline

	descStyle = lipgloss.NewStyle().
		Foreground(text).
		MarginLeft(1)

	// Birthdays
	birthdayContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		Padding(0, 1).
		MarginTop(1)

	nameStyle = lipgloss.NewStyle().
		Foreground(special).
		Bold(true)

	bYearStyle = lipgloss.NewStyle().
		Foreground(subtle).
		MarginLeft(1)
)
