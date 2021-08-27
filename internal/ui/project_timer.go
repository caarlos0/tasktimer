package ui

import (
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

type projectTimerModel struct {
	tasks []model.Task
}

func (m projectTimerModel) Init() tea.Cmd {
	return nil
}

func (m projectTimerModel) Update(msg tea.Msg) (projectTimerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case projectTimerUpdateMsg:
		m.tasks = msg.tasks
	}
	return m, nil
}

func (m projectTimerModel) View() string {
	return secondaryForeground.Render("total: ") +
		activeForegroundBold.Render(sumTasksTimes(m.tasks, time.Time{}).Round(time.Second).String()) +
		separator +
		secondaryForeground.Render("today: ") +
		activeForegroundBold.Render(sumTasksTimes(m.tasks, time.Now().Truncate(time.Hour*24)).Round(time.Second).String())
}

// msgs and cmds

type projectTimerUpdateMsg struct {
	tasks []model.Task
}

func updateProjectTimerCmd(tasks []model.Task) tea.Cmd {
	return func() tea.Msg {
		return projectTimerUpdateMsg{tasks}
	}
}

func sumTasksTimes(tasks []model.Task, since time.Time) time.Duration {
	d := time.Duration(0)
	loc := time.Now().Location()
	for _, t := range tasks {
		if t.StartAt.In(loc).Before(since) {
			continue
		}
		z := t.EndAt.In(loc)
		if z.IsZero() {
			z = time.Now()
		}
		d += z.Sub(t.StartAt.In(loc))
	}
	return d
}
