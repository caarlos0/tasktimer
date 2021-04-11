package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// light palette: https://colorhunt.co/palette/201882
// dark palette:  https://colorhunt.co/palette/273948
var (
	primaryColor = lipgloss.AdaptiveColor{
		Light: "#1a1a2e",
		Dark:  "#f7f3e9",
	}
	secondaryColor = lipgloss.AdaptiveColor{
		Light: "#16213e",
		Dark:  "#a3d2ca",
	}
	activeColor = lipgloss.AdaptiveColor{
		Light: "#16213e",
		Dark:  "#5eaaa8",
	}
	errorColor = lipgloss.AdaptiveColor{
		Light: "#e94560",
		Dark:  "#f05945",
	}

	secondaryForeground   = lipgloss.NewStyle().Foreground(secondaryColor)
	primaryForegroundBold = lipgloss.NewStyle().Bold(true).Foreground(primaryColor)
	activeForeground      = lipgloss.NewStyle().Bold(true).Foreground(activeColor)
	activeForegroundBold  = lipgloss.NewStyle().Bold(true).Foreground(activeColor)
	errorFaintForeground  = lipgloss.NewStyle().Foreground(errorColor).Faint(true)
	errorForegroundPadded = lipgloss.NewStyle().Padding(4).Foreground(errorColor)
	separator             = secondaryForeground.Render(" • ")
)

const (
	iconDone    = "●"
	iconOngoing = "○"
)
