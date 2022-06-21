package ui

import (
	"io"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbletea/teatest"
	"github.com/dgraph-io/badger/v3"
)

func TestApp(t *testing.T) {
	options := badger.DefaultOptions(filepath.Join(t.TempDir(), "db")).
		WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(options)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	m := Init(db, "test")
	teatest.TestModel(t, m, func(p teatest.Program, in io.Writer) {
		teatest.TypeText(p, "new task")
		p.Send(tea.KeyMsg{Type: tea.KeyEnter})
		time.Sleep(time.Second)
		p.Send(tea.KeyMsg{Type: tea.KeyEsc})
		p.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	}, func(out []byte) {
		teatest.RequireEqualOutput(t, out)
	})
}
