package cmd

import (
	"io"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
	gap "github.com/muesli/go-app-paths"
)

func paths(project string) (string, string, error) {
	var home = gap.NewScope(gap.User, "tasktimer")

	logfile, err := home.LogPath(project + ".log")
	if err != nil {
		return "", "", err
	}

	dbfile, err := home.DataPath(project + ".db")
	if err != nil {
		return "", "", err
	}

	return logfile, dbfile, nil
}

func setup(project string) (*badger.DB, io.Closer, error) {
	logfile, dbfile, err := paths(project)
	if err != nil {
		return nil, nil, err
	}

	if err := os.MkdirAll(filepath.Dir(logfile), 0754); err != nil {
		return nil, nil, err
	}

	f, err := tea.LogToFile(logfile, "tasktimer")
	if err != nil {
		return nil, nil, err
	}

	// TODO: maybe sync writes?
	var options = badger.DefaultOptions(dbfile).
		WithLogger(badgerStdLoggerAdapter{}).
		WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(options)
	return db, f, err
}
