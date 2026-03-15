package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const Version = "1.0.0"

var banner = `
  ___       ___
 |_ _|_ __ |_ _|
  | || '_ \ | |
  | || | | || |
 |___|_| |_|___|`

func PrintBanner() {
	gradient := lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	version := lipgloss.NewStyle().
		Foreground(Secondary).
		Italic(true).
		Render(fmt.Sprintf("  v%s", Version))

	tagline := MutedStyle.Render("  Blazing fast project initializer")

	fmt.Println(gradient.Render(banner) + version)
	fmt.Println(tagline)
	fmt.Println()
}
