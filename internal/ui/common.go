package ui

import "github.com/muesli/termenv"

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func faint(s string) string {
	return termenv.String(s).Faint().String()
}

func bold(s string) string {
	return termenv.String(s).Bold().String()
}

const (
	iconDone    = "●"
	iconOngoing = "○"
)
