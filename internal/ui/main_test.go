package ui

import (
	"io"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/caarlos0/tasktimer/internal/store"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbletea/teatest"
	"github.com/dgraph-io/badger/v3"
)

func TestApp(t *testing.T) {
	log.SetOutput(io.Discard)
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
		p.Send(tea.WindowSizeMsg{
			Width:  100,
			Height: 20,
		})
		teatest.TypeText(p, "new task")
		p.Send(tea.KeyMsg{Type: tea.KeyEnter})
		time.Sleep(time.Second)
		p.Send(tea.KeyMsg{Type: tea.KeyEsc})
		p.Send(tea.KeyMsg{Type: tea.KeyEnter})
		teatest.TypeText(p, "another task")
		p.Send(tea.KeyMsg{Type: tea.KeyEnter})
		time.Sleep(time.Millisecond * 100)
		p.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	}, func(out []byte) {
		teatest.RequireEqualOutput(t, out)
	})

	tasks, err := store.GetTaskList(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	task := tasks[1]
	if d := task.EndAt.Sub(task.StartAt); d != time.Second {
		t.Fatalf("expected 1 second task, got %v", d)
	}
}
