module github.com/caarlos0/tasktimer

go 1.16

require (
	github.com/charmbracelet/bubbles v0.8.0
	github.com/charmbracelet/bubbletea v0.14.1
	github.com/charmbracelet/glamour v0.3.0
	github.com/charmbracelet/lipgloss v0.3.0
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/mattn/go-isatty v0.0.13
	github.com/muesli/go-app-paths v0.2.1
	github.com/spf13/cobra v1.2.1
)

replace github.com/charmbracelet/bubbles => ../../charm/bubbles-internal
