package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/tasktimer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
)

var project = flag.String("project", "default", "your current project name")

func main() {
	flag.Parse()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	var folder = filepath.Join(home, "tasktimer")
	if err := os.MkdirAll(folder, 0764); err != nil {
		log.Fatalln(err)
	}

	f, err := tea.LogToFile(filepath.Join(folder, *project+".log"), "")
	defer f.Close()

	db, err := badger.Open(badger.DefaultOptions(filepath.Join(folder, *project+".db")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var p = tea.NewProgram(ui.Init(db, *project))
	p.EnterAltScreen()
	defer p.ExitAltScreen()
	if err = p.Start(); err != nil {
		log.Fatalln(err)
	}
}
