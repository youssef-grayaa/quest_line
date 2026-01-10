package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"quest_line/domain"
	"strings"
)

// DashboardModel displays today's active quests
type DashboardModel struct {
	projects           []domain.Project
	selectedProjectIdx int
	keymap             KeyMap
	selectedIdx        int
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel(projects []domain.Project, selectedProjectIdx int, keymap KeyMap) DashboardModel {
	return DashboardModel{
		projects:           projects,
		selectedProjectIdx: selectedProjectIdx,
		keymap:             keymap,
		selectedIdx:        0,
	}
}

// Update handles messages for the dashboard
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		activeQuests := domain.DailyPlanner(m.projects)
		if k == "k" || k == "up" {
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		} else if k == "j" || k == "down" {
			if m.selectedIdx < len(activeQuests)-1 {
				m.selectedIdx++
			}
		}
	}
	return m, nil
}

// View renders the dashboard
func (m DashboardModel) View() string {
	var b strings.Builder

	projectName := ""
	if m.selectedProjectIdx >= 0 && m.selectedProjectIdx < len(m.projects) {
		projectName = m.projects[m.selectedProjectIdx].Name + " - "
	}
	b.WriteString(titleStyle.Render("Dashboard - " + projectName + "Today's Active Quests"))
	b.WriteString("\n\n")

	var activeQuests []domain.Quest
	if m.selectedProjectIdx >= 0 {
		activeQuests = domain.DailyPlannerForProject(m.projects, m.selectedProjectIdx)
	} else {
		activeQuests = domain.DailyPlanner(m.projects)
	}

	if m.selectedIdx >= len(activeQuests) {
		m.selectedIdx = 0
	}

	if len(activeQuests) == 0 {
		b.WriteString("No active quests!\n")
		b.WriteString("Create a quest to get started.\n\n")
	} else {
		for i, quest := range activeQuests {
			line := fmt.Sprintf("%s: %.1f%% complete", quest.Title, quest.Progress)
			if i == m.selectedIdx {
				line = selectedStyle.Render(line)
			}
			b.WriteString(line + "\n")
			if quest.Deadline != nil {
				b.WriteString(fmt.Sprintf("  Due: %s\n", quest.Deadline.Format("2006-01-02")))
			}
			b.WriteString(fmt.Sprintf("  Priority: %d\n\n", quest.Priority))
		}
	}

	return b.String()
}

// SelectedIndex returns the currently selected index
func (m DashboardModel) SelectedIndex() int {
	return m.selectedIdx
}

// ProjectListModel displays all projects
type ProjectListModel struct {
	projects    []domain.Project
	selectedIdx int
	keymap      KeyMap
}

// NewProjectListModel creates a new project list model
func NewProjectListModel(projects []domain.Project, keymap KeyMap) ProjectListModel {
	return ProjectListModel{
		projects:    projects,
		selectedIdx: 0,
		keymap:      keymap,
	}
}

// Update handles messages for the project list
func (m ProjectListModel) Update(msg tea.Msg) (ProjectListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "k" || k == "up" {
			m.MoveUp()
		} else if k == "j" || k == "down" {
			m.MoveDown()
		}
	case DataChangedMsg:
		m.projects = msg.Projects
		if m.selectedIdx >= len(m.projects) {
			m.selectedIdx = 0
		}
	}
	return m, nil
}

// View renders the project list
func (m ProjectListModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Projects"))
	b.WriteString("\n\n")

	if len(m.projects) == 0 {
		b.WriteString("No projects. Press 'c' to create one.\n\n")
	} else {
		for i, project := range m.projects {
			line := fmt.Sprintf("%2d. %s (%d quests)", i+1, project.Name, len(project.Quests))
			if i == m.selectedIdx {
				line = selectedStyle.Render(line)
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// ProjectSelectionModel displays available projects for selection at startup
type ProjectSelectionModel struct {
	projects    []domain.Project
	selectedIdx int
	keymap      KeyMap
}

// NewProjectSelectionModel creates a new project selection model
func NewProjectSelectionModel(projects []domain.Project, keymap KeyMap) ProjectSelectionModel {
	return ProjectSelectionModel{
		projects:    projects,
		selectedIdx: 0,
		keymap:      keymap,
	}
}

// Update handles messages for the project selection
func (m ProjectSelectionModel) Update(msg tea.Msg) (ProjectSelectionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "k" || k == "up" {
			m.MoveUp()
		} else if k == "j" || k == "down" {
			m.MoveDown()
		}
	}
	return m, nil
}

// View renders the project selection
func (m ProjectSelectionModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Select Project to Work On"))
	b.WriteString("\n\n")

	if len(m.projects) == 0 {
		b.WriteString("No projects available.\n")
	} else {
		for i, project := range m.projects {
			line := fmt.Sprintf("%2d. %s (%d quests)", i+1, project.Name, len(project.Quests))
			if i == m.selectedIdx {
				line = selectedStyle.Render(line)
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// MoveUp moves selection up
func (m *ProjectSelectionModel) MoveUp() {
	if m.selectedIdx > 0 {
		m.selectedIdx--
	}
}

// MoveDown moves selection down
func (m *ProjectSelectionModel) MoveDown() {
	if m.selectedIdx < len(m.projects)-1 {
		m.selectedIdx++
	}
}

// SelectedIndex returns the currently selected index
func (m ProjectSelectionModel) SelectedIndex() int {
	return m.selectedIdx
}

// MoveUp moves selection up
func (m *ProjectListModel) MoveUp() {
	if m.selectedIdx > 0 {
		m.selectedIdx--
	}
}

// MoveDown moves selection down
func (m *ProjectListModel) MoveDown() {
	if m.selectedIdx < len(m.projects)-1 {
		m.selectedIdx++
	}
}

// SelectedProject returns the currently selected project
func (m ProjectListModel) SelectedProject() *domain.Project {
	if m.selectedIdx >= 0 && m.selectedIdx < len(m.projects) {
		return &m.projects[m.selectedIdx]
	}
	return nil
}

// SelectedIndex returns the currently selected index
func (m ProjectListModel) SelectedIndex() int {
	return m.selectedIdx
}

// QuestDetailModel displays a single quest with its tasks
type QuestDetailModel struct {
	projects         []domain.Project
	selectedProjIdx  int
	selectedQuestIdx int
	selectedTaskIdx  int
	keymap           KeyMap
	taskList         list.Model
}

// NewQuestDetailModel creates a new quest detail model
func NewQuestDetailModel(projects []domain.Project, projIdx int, questIdx int, keymap KeyMap) QuestDetailModel {
	var items []list.Item
	if projIdx >= 0 && questIdx >= 0 && projIdx < len(projects) && questIdx < len(projects[projIdx].Quests) {
		for i := range projects[projIdx].Quests[questIdx].Tasks {
			items = append(items, TaskItem{task: &projects[projIdx].Quests[questIdx].Tasks[i]})
		}
	}

	taskDelegate := newTaskDelegate(newDelegateKeyMap())
	taskList := list.New(items, taskDelegate, 0, 0)
	taskList.Title = "Tasks"
	taskList.Styles.Title = titleStyle

	// Set default size in case WindowSizeMsg is not received
	taskList.SetSize(80, 15)

	taskList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			newListKeyMap().toggleTitleBar,
			newListKeyMap().toggleStatusBar,
			newListKeyMap().togglePagination,
			newListKeyMap().toggleHelpMenu,
		}
	}

	return QuestDetailModel{
		projects:         projects,
		selectedProjIdx:  projIdx,
		selectedQuestIdx: questIdx,
		selectedTaskIdx:  0,
		keymap:           keymap,
		taskList:         taskList,
	}
}

// Update handles messages for the quest detail
func (m QuestDetailModel) Update(msg tea.Msg) (QuestDetailModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.taskList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Toggle selected task
			if item, ok := m.taskList.SelectedItem().(TaskItem); ok && item.task != nil {
				item.task.Done = !item.task.Done
				// Update progress
				if m.selectedProjIdx >= 0 && m.selectedQuestIdx >= 0 &&
					m.selectedProjIdx < len(m.projects) &&
					m.selectedQuestIdx < len(m.projects[m.selectedProjIdx].Quests) {
					quest := &m.projects[m.selectedProjIdx].Quests[m.selectedQuestIdx]
					quest.CalculateProgress()
					m.projects[m.selectedProjIdx].CalculateProgress()
					// Send save command
					cmd = func() tea.Msg {
						err := domain.SaveProjects(m.projects)
						return SaveCompleteMsg{Err: err}
					}
				}
				// Refresh task list items
				var items []list.Item
				if m.selectedProjIdx >= 0 && m.selectedQuestIdx >= 0 &&
					m.selectedProjIdx < len(m.projects) &&
					m.selectedQuestIdx < len(m.projects[m.selectedProjIdx].Quests) {
					for i := range m.projects[m.selectedProjIdx].Quests[m.selectedQuestIdx].Tasks {
						items = append(items, TaskItem{task: &m.projects[m.selectedProjIdx].Quests[m.selectedQuestIdx].Tasks[i]})
					}
				}
				m.taskList.SetItems(items)
			}
		default:
			// Delegate to list
			m.taskList, cmd = m.taskList.Update(msg)
		}
	case DataChangedMsg:
		m.projects = msg.Projects
		// Keep indices valid
		if m.selectedProjIdx >= len(m.projects) {
			m.selectedProjIdx = -1
			m.selectedQuestIdx = -1
		} else if m.selectedQuestIdx >= len(m.projects[m.selectedProjIdx].Quests) {
			m.selectedQuestIdx = -1
		}
		// Update items
		var items []list.Item
		if m.selectedProjIdx >= 0 && m.selectedQuestIdx >= 0 &&
			m.selectedProjIdx < len(m.projects) &&
			m.selectedQuestIdx < len(m.projects[m.selectedProjIdx].Quests) {
			for i := range m.projects[m.selectedProjIdx].Quests[m.selectedQuestIdx].Tasks {
				items = append(items, TaskItem{task: &m.projects[m.selectedProjIdx].Quests[m.selectedQuestIdx].Tasks[i]})
			}
		}
		m.taskList.SetItems(items)
	}
	m.selectedTaskIdx = m.taskList.Index()
	return m, cmd
}

// View renders the quest detail
func (m QuestDetailModel) View() string {
	if m.selectedProjIdx < 0 || m.selectedQuestIdx < 0 ||
		m.selectedProjIdx >= len(m.projects) ||
		m.selectedQuestIdx >= len(m.projects[m.selectedProjIdx].Quests) {
		return "No quest selected."
	}

	quest := &m.projects[m.selectedProjIdx].Quests[m.selectedQuestIdx]

	var b strings.Builder

	b.WriteString(titleStyle.Render("Quest: " + quest.Title))
	b.WriteString("\n")
	b.WriteString(quest.Description)
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Progress: %.1f%% | Priority: %d | Status: %s\n",
		quest.Progress,
		quest.Priority,
		quest.State.String()))

	if quest.Deadline != nil {
		b.WriteString(fmt.Sprintf("Deadline: %s\n", quest.Deadline.Format("2006-01-02")))
	}

	b.WriteString("\n\n")

	b.WriteString(m.taskList.View())

	return b.String()
}

// SelectedTaskIndex returns the selected task index
func (m QuestDetailModel) SelectedTaskIndex() int {
	return m.taskList.Index()
}

// Styling for tasks
var (
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#d90368")).
			Foreground(lipgloss.Color("#FFFDF5"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
)
