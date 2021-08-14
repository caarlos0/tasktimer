package ui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/tasktimer/internal/model"
	"github.com/caarlos0/tasktimer/internal/store"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgraph-io/badger/v3"
)

func Init(db *badger.DB, project string) tea.Model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.Placeholder = "New task description..."
	input.Focus()
	input.CharLimit = 250
	input.Width = 50

	l := list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Task List"
	l.SetSpinner(spinner.MiniDot)
	l.DisableQuitKeybindings()

	return mainModel{
		list:    l,
		timer:   projectTimerModel{},
		db:      db,
		input:   input,
		project: project,
	}
}

type mainModel struct {
	input   textinput.Model
	list    list.Model
	timer   projectTimerModel
	db      *badger.DB
	project string
	err     error
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		m.list.StartSpinner(),
		askForTaskListUpdateCmd(),
		textinput.Blink,
	)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var newMsg tea.Msg

	switch msg := msg.(type) {
	case errMsg:
		log.Println("errMsg")
		m.err = msg.error
	case tea.WindowSizeMsg:
		log.Println("tea.WindowSizeMsg")
		top, right, bottom, left := listStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case updateTaskListMsg:
		log.Println("updateTaskListMsg")
		cmds = append(cmds, m.list.StartSpinner(), updateTaskListCmd(m.db))
	case taskListUpdatedMsg:
		log.Println("taskListUpdatedMsg")
		var items = make([]list.Item, 0, len(msg.tasks))
		for _, t := range msg.tasks {
			items = append(items, item{
				title: t.Title,
				start: t.StartAt,
				end:   t.EndAt,
			})
		}

		m.list.StopSpinner()
		m.list.SetItems(items)
		cmds = append(
			cmds,
			m.list.ResetSelected(),
			m.list.ResetFilter(),
			updateProjectTimerCmd(msg.tasks),
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			log.Println("tea.KeyMsg -> ctrl+c")
			return m, tea.Sequentially(closeTasksCmd(m.db), tea.Quit)
		case "esc":
			log.Println("tea.KeyMsg -> esc")
			if !m.list.SettingFilter() {
				log.Println("tea.KeyMsg -> esc -> !SettingFilter")
				if m.input.Focused() {
					log.Println("tea.KeyMsg -> esc -> !SettingFilter -> input.Focused")
					m.input.Blur()
					cmds = append(cmds, tea.Sequentially(
						closeTasksCmd(m.db),
						updateTaskListCmd(m.db)),
					)
				}
				newMsg = doNotPropagateMsg{}
			}
		case "enter":
			log.Println("tea.KeyMsg -> enter")
			if !m.list.SettingFilter() {
				log.Println("tea.KeyMsg -> enter -> !SettingFilter")
				if !m.input.Focused() {
					log.Println("tea.KeyMsg -> enter -> !SettingFilter -> input.Focused")
					m.input.Focus()
					cmds = append(cmds, textinput.Blink)
				} else {
					log.Println("tea.KeyMsg -> enter -> !SettingFilter -> !input.Focused")
					cmds = append(cmds, tea.Sequentially(
						closeTasksCmd(m.db),
						createTaskCmd(m.db, strings.TrimSpace(m.input.Value())),
					))
					m.input.SetValue("")
				}
			}
		default:
			log.Println("tea.KeyMsg -> default")
			if m.input.Focused() {
				log.Println("tea.KeyMsg -> default -> input.Focused")
				// only send key presses to input if it is focused
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
				newMsg = doNotPropagateMsg{}
			}
		}
	}

	if newMsg != nil {
		log.Println("tea.KeyMsg -> override original msg")
		msg = newMsg
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	if m.err != nil {
		return "\n" +
			errorFaintForeground.Render("Oops, something went wrong:") +
			"\n\n" +
			errorForegroundPadded.Render(m.err.Error()) +
			"\n\n" +
			errorFaintForeground.Render("Check the logs for more details...")
	}
	return secondaryForeground.Render("project: ") +
		activeForegroundBold.Render(m.project) +
		separator + m.timer.View() + "\n\n" +
		m.input.View() + "\n\n" +
		m.list.View() + "\n"
}

// msgs

type doNotPropagateMsg struct{}

type updateTaskListMsg struct{}

type taskListUpdatedMsg struct {
	tasks []model.Task
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// cmds

func closeTasksCmd(db *badger.DB) tea.Cmd {
	return func() tea.Msg {
		log.Println("closeTasksCmd")
		if err := store.CloseTasks(db); err != nil {
			return errMsg{err}
		}
		return nil
	}
}

func createTaskCmd(db *badger.DB, t string) tea.Cmd {
	return func() tea.Msg {
		log.Println("createTaskCmd")
		if err := store.CreateTask(db, t); err != nil {
			return errMsg{err}
		}
		return updateTaskListMsg{}
	}
}

func askForTaskListUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		log.Println("askForTaskListUpdateCmd")
		return updateTaskListMsg{}
	}
}

func updateTaskListCmd(db *badger.DB) tea.Cmd {
	return func() tea.Msg {
		log.Println("updateTaskListCmd")
		time.Sleep(time.Second)
		tasks, err := store.GetTaskList(db)
		if err != nil {
			return errMsg{err}
		}
		return taskListUpdatedMsg{tasks}
	}
}

// models

type item struct {
	title      string
	start, end time.Time
}

func (i item) Title() string {
	if i.end.IsZero() {
		return boldStyle.Render(i.title)
	}
	return i.title
}

func (i item) Description() string {
	end := time.Now()
	if !i.end.IsZero() {
		end = i.end
	}
	return relativeDays(i.start) + " - " + end.Sub(i.start).Round(time.Second).String()
}

func (i item) FilterValue() string { return i.title }

func relativeDays(t time.Time) string {
	days := time.Since(t).Round(time.Hour*24).Hours() / 24
	if days == 0 {
		return "today"
	}
	if days == 1 {
		return "yesterday"
	}
	return fmt.Sprintf("%.0f days ago", days)
}
