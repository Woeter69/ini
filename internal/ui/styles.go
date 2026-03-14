package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Brand colors
	Primary   = lipgloss.Color("#7C3AED") // Vibrant purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Success   = lipgloss.Color("#10B981") // Emerald green
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Error     = lipgloss.Color("#EF4444") // Red
	Muted     = lipgloss.Color("#6B7280") // Gray
	White     = lipgloss.Color("#F9FAFB") // Near white

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// Box for final summary
	SummaryBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2).
			MarginTop(1)

	// Checkmark and cross
	CheckMark = SuccessStyle.Render("✓")
	CrossMark = ErrorStyle.Render("✗")
	Arrow     = lipgloss.NewStyle().Foreground(Secondary).Render("→")
)
