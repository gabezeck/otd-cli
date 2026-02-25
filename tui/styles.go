package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	subtle  = lipgloss.AdaptiveColor{Light: "#C9C1B5", Dark: "#4B4741"}
	accent  = lipgloss.AdaptiveColor{Light: "#B96A3A", Dark: "#D18A5E"}
	muted   = lipgloss.AdaptiveColor{Light: "#6B6258", Dark: "#B7AFA4"}
	text    = lipgloss.AdaptiveColor{Light: "#2A2520", Dark: "#E6DED2"}
	soft    = lipgloss.AdaptiveColor{Light: "#8E8073", Dark: "#8A8378"}
	accent2 = lipgloss.AdaptiveColor{Light: "#7C8A6B", Dark: "#9AA787"}

	// Header
	headerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accent)

	headerDateStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(soft)

	headerSummaryStyle = lipgloss.NewStyle().
				Foreground(muted)

	dividerStyle = lipgloss.NewStyle().
			Foreground(subtle)

	// Sections
	sectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accent2).
				MarginTop(1).
				MarginBottom(0)

	// Timeline / Events
	yearStyle = lipgloss.NewStyle().
			Foreground(accent).
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
	nameStyle = lipgloss.NewStyle().
			Foreground(accent2).
			Bold(true)

	bYearStyle = lipgloss.NewStyle().
			Foreground(muted).
			MarginLeft(1)
)
