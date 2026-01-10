package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"quest_line/domain"
)

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update screen models with new size if needed
		// For now, assume fixed size
	case tea.KeyMsg:
		// Handle form view
		if m.inForm {
			return m.handleFormInput(msg)
		}
		// Handle delete confirmation
		if m.pendingDelete {
			return m.handleDeleteConfirmation(msg)
		}
		// Global quit handler
		if key.Matches(msg, m.keymap.Quit) {
			m.saveProjects()
			return m, tea.Quit
		}
		// Help toggle
		if key.Matches(msg, m.keymap.Help) {
			m.help.ToggleHelp()
			return m, nil
		}
		// Handle view-specific keys
		return m.handleViewSpecificInput(msg)
	case SaveCompleteMsg:
		if msg.Err != nil {
			m.errorMsg = msg.Err.Error()
		} else {
			m.errorMsg = ""
		}
		// Update screen models with data changed
		m.updateScreenModels()
		return m, nil
	case DataChangedMsg:
		m.projects = msg.Projects
		m.updateScreenModels()
		return m, nil
	}
	return m, nil
}

func (m *RootModel) handleFormInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.cancelForm()
		return m, nil
	}

	if msg.String() == "enter" {
		if m.form.focusIdx < len(m.form.inputs)-1 {
			m.form.nextInput()
			return m, nil
		} else {
			if err := m.form.Validate(); err != nil {
				m.form.errorMsg = err.Error()
				return m, nil
			}
			m.form.errorMsg = ""
			m.submitForm()
			return m, m.saveProjectsCmd()
		}
	}

	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)

	return m, cmd
}

func (m *RootModel) handleDeleteConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.confirmDelete()
		return m, m.saveProjectsCmd()
	case "n", "N":
		m.pendingDelete = false
		return m, nil
	}
	return m, nil
}

func (m *RootModel) handleViewSpecificInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentView {
	case ViewProjectSelection:
		if msg.String() == "enter" {
			selectedIdx := m.projectSelection.SelectedIndex()
			if selectedIdx >= 0 && selectedIdx < len(m.projects) {
				m.selectedProjectIdx = selectedIdx
				m.navigateTo(ViewDashboard)
			}
			return m, nil
		}
		// Delegate to projectSelection
		m.projectSelection, cmd = m.projectSelection.Update(msg)
		return m, cmd
	case ViewDashboard:
		switch {
		case key.Matches(msg, m.keymap.Create):
			if m.selectedProjectIdx >= 0 {
				m.startCreateQuest()
			} else {
				m.startCreateProject()
			}
			return m, nil
		case key.Matches(msg, m.keymap.Edit):
			selectedIdx := m.dashboard.SelectedIndex()
			var activeQuests []domain.Quest
			if m.selectedProjectIdx >= 0 {
				activeQuests = domain.DailyPlannerForProject(m.projects, m.selectedProjectIdx)
			} else {
				activeQuests = domain.DailyPlanner(m.projects)
			}
			if selectedIdx >= 0 && selectedIdx < len(activeQuests) {
				questID := activeQuests[selectedIdx].ID
				pIdx, qIdx := domain.FindQuestIndices(m.projects, questID)
				if pIdx >= 0 && qIdx >= 0 {
					m.selectedProjectIdx = pIdx
					m.selectedQuestIdx = qIdx
					m.startEditQuest()
				}
			}
			return m, nil
		case key.Matches(msg, m.keymap.Delete):
			selectedIdx := m.dashboard.SelectedIndex()
			var activeQuests []domain.Quest
			if m.selectedProjectIdx >= 0 {
				activeQuests = domain.DailyPlannerForProject(m.projects, m.selectedProjectIdx)
			} else {
				activeQuests = domain.DailyPlanner(m.projects)
			}
			if selectedIdx >= 0 && selectedIdx < len(activeQuests) {
				questID := activeQuests[selectedIdx].ID
				pIdx, qIdx := domain.FindQuestIndices(m.projects, questID)
				if pIdx >= 0 && qIdx >= 0 {
					// Confirm delete
					m.pendingDelete = true
					m.deleteType = "quest"
					m.deleteIndices = [3]int{pIdx, qIdx, -1}
				}
			}
			return m, nil
		}
		if msg.String() == "enter" {
			selectedIdx := m.dashboard.SelectedIndex()
			activeQuests := domain.DailyPlanner(m.projects)
			if selectedIdx >= 0 && selectedIdx < len(activeQuests) {
				questID := activeQuests[selectedIdx].ID
				pIdx, qIdx := domain.FindQuestIndices(m.projects, questID)
				if pIdx >= 0 && qIdx >= 0 {
					m.selectedProjectIdx = pIdx
					m.selectedQuestIdx = qIdx
					m.updateTaskList()
					m.navigateTo(ViewQuestDetail)
				}
			}
			return m, nil
		}
		// Delegate to dashboard
		var cmd tea.Cmd
		m.dashboard, cmd = m.dashboard.Update(msg)
		return m, cmd
	case ViewProjectList:
		switch {
		case key.Matches(msg, m.keymap.Create):
			selectedIdx := m.projectList.SelectedIndex()
			if selectedIdx >= 0 && selectedIdx < len(m.projects) {
				m.selectedProjectIdx = selectedIdx
				m.startCreateQuest()
			} else {
				m.startCreateProject()
			}
			return m, nil
		case key.Matches(msg, m.keymap.Edit):
			m.startEditProject()
			return m, nil
		case key.Matches(msg, m.keymap.Delete):
			m.startDeleteProject()
			return m, nil
		case key.Matches(msg, m.keymap.Dashboard):
			m.navigateTo(ViewDashboard)
			return m, nil

		}
		// Handle selection
		if msg.String() == "enter" {
			selectedIdx := m.projectList.SelectedIndex()
			if selectedIdx >= 0 && selectedIdx < len(m.projects) {
				m.selectedProjectIdx = selectedIdx
				project := m.projects[selectedIdx]
				if len(project.Quests) > 0 {
					m.selectedQuestIdx = 0
					m.updateTaskList()
					m.navigateTo(ViewQuestDetail)
				} else {
					m.navigateTo(ViewDashboard)
				}
			}
			return m, nil
		}
		// Delegate to projectList
		var cmd tea.Cmd
		m.projectList, cmd = m.projectList.Update(msg)
		return m, cmd
	case ViewQuestDetail:
		switch {
		case key.Matches(msg, m.keymap.Toggle):
			m.toggleTask()
			return m, m.saveProjectsCmd()
		case key.Matches(msg, m.keymap.Create):
			m.startCreateTask()
			return m, nil
		case key.Matches(msg, m.keymap.Edit):
			m.startEditTask()
			return m, nil
		case key.Matches(msg, m.keymap.Delete):
			m.startDeleteTask()
			return m, nil
		case key.Matches(msg, m.keymap.Dashboard):
			m.navigateTo(ViewDashboard)
			return m, nil
		}
		m.taskList, cmd = m.taskList.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *RootModel) updateScreenModels() {
	m.dashboard = NewDashboardModel(m.projects, m.selectedProjectIdx, m.keymap)
	m.projectList = NewProjectListModel(m.projects, m.keymap)
	if m.selectedQuestIdx >= 0 {
		m.taskList = NewQuestDetailModel(m.projects, m.selectedProjectIdx, m.selectedQuestIdx, m.keymap)
	} else {
		m.taskList = NewQuestDetailModel(m.projects, -1, -1, m.keymap)
	}
}

func (m *RootModel) navigateTo(view View) {
	m.currentView = view
	if view == ViewQuestDetail {
		m.updateTaskList()
	}
	m.updateScreenModels()
}

func (m *RootModel) updateTaskList() {
	if m.selectedQuestIdx >= 0 {
		m.taskList = NewQuestDetailModel(m.projects, m.selectedProjectIdx, m.selectedQuestIdx, m.keymap)
	}
}
