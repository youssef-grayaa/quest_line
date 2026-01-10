package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"quest_line/domain"
)

// View represents the current screen
type View int

const (
	ViewProjectSelection View = iota
	ViewDashboard
	ViewProjectList
	ViewQuestDetail
	ViewCreateProject
	ViewEditProject
	ViewCreateQuest
	ViewEditQuest
	ViewCreateTask
	ViewEditTask
)

// ProjectItem represents a project in the list
type ProjectItem struct {
	project *domain.Project
}

// FilterValue returns the value to filter by
func (p ProjectItem) FilterValue() string {
	return p.project.Name
}

// Title returns the title
func (p ProjectItem) Title() string {
	return p.project.Name
}

// Description returns the description
func (p ProjectItem) Description() string {
	return fmt.Sprintf("%.1f%% complete - %d quests", p.project.Progress, len(p.project.Quests))
}

// QuestItem represents a quest in the list
type QuestItem struct {
	quest *domain.Quest
}

// FilterValue returns the value to filter by
func (q QuestItem) FilterValue() string {
	return q.quest.Title
}

// Title returns the title
func (q QuestItem) Title() string {
	return q.quest.Title
}

// Description returns the description
func (q QuestItem) Description() string {
	return fmt.Sprintf("%s\n%.1f%% complete - %d tasks", q.quest.Description, q.quest.Progress, len(q.quest.Tasks))
}

// TaskItem represents a task in the list
type TaskItem struct {
	task *domain.Task
}

// FilterValue returns the value to filter by
func (t TaskItem) FilterValue() string {
	return t.task.Description
}

// Title returns the title
func (t TaskItem) Title() string {
	status := "[ ]"
	if t.task.Done {
		status = "[✓]"
	}
	return fmt.Sprintf("%s %s", status, t.task.Description)
}

// Description returns the description
func (t TaskItem) Description() string {
	return ""
}

// listKeyMap defines keys for the list
type listKeyMap struct {
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

// delegateKeyMap defines keys for list items
type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "toggle"),
		),
		remove: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "delete"),
		),
	}
}

// newItemDelegate creates a custom delegate for list items
func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// Item-specific actions can be handled here
		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{keys.choose, keys.remove}
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{{keys.choose, keys.remove}}
	}

	return d
}

// newTaskDelegate creates a custom delegate for task items
func newTaskDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// Task-specific actions handled in model Update
		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{keys.choose} // Only choose for tasks (toggle via space)
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{{keys.choose}}
	}

	return d
}

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	// styling for title
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#4f4686")).
			Padding(0, 1)
)
