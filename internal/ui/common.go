package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// light palette: https://colorhunt.co/palette/201882
// dark palette:  https://colorhunt.co/palette/273948
var (
	defaultStyles = list.NewDefaultItemStyles()

	activeColor    = defaultStyles.SelectedTitle.GetForeground()
	secondaryColor = defaultStyles.NormalTitle.GetForeground()

	errorColor = lipgloss.AdaptiveColor{
		Light: "#e94560",
		Dark:  "#f05945",
	}

	secondaryForeground   = lipgloss.NewStyle().Foreground(secondaryColor)
	boldStyle             = lipgloss.NewStyle().Bold(true)
	activeForegroundBold  = lipgloss.NewStyle().Bold(true).Foreground(activeColor)
	errorFaintForeground  = lipgloss.NewStyle().Foreground(errorColor).Faint(true)
	errorForegroundPadded = lipgloss.NewStyle().Padding(4).Foreground(errorColor)
	separator             = secondaryForeground.Render(" • ")
	listStyle             = lipgloss.NewStyle().Margin(6, 2, 0, 2)
)
