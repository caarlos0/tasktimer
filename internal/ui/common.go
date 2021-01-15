package ui

import "github.com/muesli/termenv"

var (
	primary   = termenv.ColorProfile().Color("205")
	secondary = termenv.ColorProfile().Color("#89F0CB")
	gray      = termenv.ColorProfile().Color("#626262")
	midGray   = termenv.ColorProfile().Color("#4A4A4A")
	red       = termenv.ColorProfile().Color("#ED567A")
)

const (
	iconDone    = "●"
	iconOngoing = "○"
)

var separator = midGrayForeground(" • ")

func redForeground(s string) string {
	return termenv.String(s).Foreground(red).String()
}

func redFaintForeground(s string) string {
	return termenv.String(s).Foreground(red).Faint().String()
}

func boldPrimaryForeground(s string) string {
	return termenv.String(s).Foreground(primary).Bold().String()
}

func boldSecondaryForeground(s string) string {
	return termenv.String(s).Foreground(secondary).Bold().String()
}

func secondaryForeground(s string) string {
	return termenv.String(s).Foreground(secondary).String()
}

func grayForeground(s string) string {
	return termenv.String(s).Foreground(gray).String()
}

func midGrayForeground(s string) string {
	return termenv.String(s).Foreground(midGray).String()
}

func faint(s string) string {
	return termenv.String(s).Faint().String()
}

func bold(s string) string {
	return termenv.String(s).Bold().String()
}
