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
		activeForegroundBold.Render(SumTasksTimes(m.tasks, time.Time{}).Round(time.Second).String()) +
		separator +
		secondaryForeground.Render("today: ") +
		activeForegroundBold.Render(SumTasksTimes(m.tasks, todayAtMidnight()).Round(time.Second).String())
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

func effectiveEndAt(endAt time.Time, now time.Time) time.Time {
	if endAt.IsZero() {
		return now
	}
	return endAt
}

func taskDuration(t model.Task, now time.Time) time.Duration {
	var d time.Duration
	switch {
	case !t.EndAt.IsZero():
		d = t.EndAt.Sub(t.StartAt) - t.PausedFor
	case !t.PausedAt.IsZero():
		d = t.PausedAt.Sub(t.StartAt) - t.PausedFor
	default:
		d = now.Sub(t.StartAt) - t.PausedFor
	}
	if d < 0 {
		return 0
	}
	return d
}

func SumTasksTimes(tasks []model.Task, since time.Time) time.Duration {
	d := time.Duration(0)
	now := time.Now()
	for _, t := range tasks {
		if t.StartAt.Before(since) {
			continue
		}
		d += taskDuration(t, now)
	}
	return d
}

func todayAtMidnight() time.Time {
	return time.Now().Truncate(time.Hour * 24)
}
