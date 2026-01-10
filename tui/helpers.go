package tui

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"quest_line/domain"
)

// Form lifecycle
func (m *RootModel) startCreateProject() {
	m.form = NewProjectForm("Create Project", "")
	m.currentView = ViewCreateProject
	m.inForm = true
	m.editingIdx = -1
}

func (m *RootModel) startEditProject() {
	if m.selectedProjectIdx >= 0 && m.selectedProjectIdx < len(m.projects) {
		p := &m.projects[m.selectedProjectIdx]
		m.form = NewProjectForm("Edit Project", p.Name)
		m.currentView = ViewEditProject
		m.inForm = true
		m.editingIdx = m.selectedProjectIdx
	}
}

func (m *RootModel) startCreateQuest() {
	if m.selectedProjectIdx >= 0 {
		m.form = NewQuestForm("Create Quest", nil)
		m.currentView = ViewCreateQuest
		m.inForm = true
		m.editingIdx = -1
	}
}

func (m *RootModel) startEditQuest() {
	if m.selectedProjectIdx >= 0 && m.selectedQuestIdx >= 0 &&
		m.selectedProjectIdx < len(m.projects) &&
		m.selectedQuestIdx < len(m.projects[m.selectedProjectIdx].Quests) {
		q := &m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx]
		m.form = NewQuestForm("Edit Quest", q)
		m.currentView = ViewEditQuest
		m.inForm = true
		m.editingIdx = m.selectedQuestIdx
	}
}

func (m *RootModel) startCreateTask() {
	if m.selectedQuestIdx >= 0 {
		m.form = NewTaskForm("Create Task", "")
		m.currentView = ViewCreateTask
		m.inForm = true
		m.editingIdx = -1
	}
}

func (m *RootModel) startEditTask() {
	if m.selectedProjectIdx >= 0 && m.selectedQuestIdx >= 0 &&
		m.selectedProjectIdx < len(m.projects) &&
		m.selectedQuestIdx < len(m.projects[m.selectedProjectIdx].Quests) &&
		m.taskList.SelectedTaskIndex() >= 0 {
		taskIdx := m.taskList.SelectedTaskIndex()
		if taskIdx < len(m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks) {
			desc := m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks[taskIdx].Description
			m.form = NewTaskForm("Edit Task", desc)
			m.currentView = ViewEditTask
			m.inForm = true
			m.editingIdx = taskIdx
		}
	}
}

// Form submission and cancellation
func (m *RootModel) submitForm() {
	switch m.currentView {
	case ViewCreateProject:
		m.createProject()
	case ViewEditProject:
		m.updateProject()
	case ViewCreateQuest:
		m.createQuest()
	case ViewEditQuest:
		m.updateQuest()
	case ViewCreateTask:
		m.createTask()
	case ViewEditTask:
		m.updateTask()
	}
	m.exitForm()
}

func (m *RootModel) cancelForm() {
	m.exitForm()
}

func (m *RootModel) exitForm() {
	m.inForm = false
	m.editingIdx = -1
	switch m.currentView {
	case ViewCreateProject, ViewEditProject:
		m.currentView = ViewProjectList
	case ViewCreateQuest, ViewEditQuest:
		m.currentView = ViewDashboard
	case ViewCreateTask, ViewEditTask:
		m.currentView = ViewQuestDetail
	}
}

// CRUD Operations
func (m *RootModel) createProject() {
	values := m.form.GetValues()
	name := strings.TrimSpace(values["Name:"])
	if name != "" {
		domain.CreateProject(&m.projects, name)
		m.updateScreenModels()
	}
}

func (m *RootModel) updateProject() {
	values := m.form.GetValues()
	name := strings.TrimSpace(values["Name:"])
	if name != "" && m.editingIdx >= 0 {
		domain.UpdateProject(&m.projects, m.editingIdx, name)
		m.updateScreenModels()
	}
}

func (m *RootModel) startDeleteProject() {
	idx := m.projectList.SelectedIndex()
	if idx >= 0 && idx < len(m.projects) {
		m.pendingDelete = true
		m.deleteType = "project"
		m.deleteIndices = [3]int{idx, -1, -1}
	}
}

func (m *RootModel) confirmDelete() {
	switch m.deleteType {
	case "project":
		idx := m.deleteIndices[0]
		if idx >= 0 && idx < len(m.projects) {
			domain.DeleteProject(&m.projects, idx)
			if m.selectedProjectIdx == idx {
				m.selectedProjectIdx = -1
			} else if m.selectedProjectIdx > idx {
				m.selectedProjectIdx--
			}
			m.updateScreenModels()
		}
	case "quest":
		if m.selectedProjectIdx >= 0 {
			idx := m.deleteIndices[1]
			if idx >= 0 && idx < len(m.projects[m.selectedProjectIdx].Quests) {
				domain.DeleteQuest(&m.projects, m.selectedProjectIdx, idx)
				if m.selectedQuestIdx == idx {
					m.selectedQuestIdx = -1
				} else if m.selectedQuestIdx > idx {
					m.selectedQuestIdx--
				}
				m.updateScreenModels()
			}
		}
	case "task":
		if m.selectedProjectIdx >= 0 && m.selectedQuestIdx >= 0 {
			idx := m.deleteIndices[2]
			if idx >= 0 && idx < len(m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks) {
				domain.DeleteTask(&m.projects, m.selectedProjectIdx, m.selectedQuestIdx, idx)
				m.updateScreenModels()
			}
		}
	}
	m.pendingDelete = false
}

func (m *RootModel) createQuest() {
	if m.selectedProjectIdx < 0 {
		return
	}
	values := m.form.GetValues()
	title := strings.TrimSpace(values["Title:"])
	if title == "" {
		return
	}

	desc := strings.TrimSpace(values["Description:"])
	priorityStr := strings.TrimSpace(values["Priority (0-10):"])
	priority := 0
	if p, err := strconv.Atoi(priorityStr); err == nil {
		priority = p
	}

	var deadline *time.Time
	deadlineStr := strings.TrimSpace(values["Deadline (YYYY-MM-DD):"])
	if deadlineStr != "" {
		if d, err := time.Parse("2006-01-02", deadlineStr); err == nil {
			deadline = &d
		}
	}

	domain.CreateQuest(&m.projects, m.selectedProjectIdx, title, desc, priority, deadline)
	m.updateScreenModels()
}

func (m *RootModel) updateQuest() {
	if m.selectedProjectIdx < 0 || m.editingIdx < 0 {
		return
	}
	values := m.form.GetValues()
	title := strings.TrimSpace(values["Title:"])
	if title == "" {
		return
	}

	desc := strings.TrimSpace(values["Description:"])
	priorityStr := strings.TrimSpace(values["Priority (0-10):"])
	priority := 0
	if p, err := strconv.Atoi(priorityStr); err == nil {
		priority = p
	}

	var deadline *time.Time
	deadlineStr := strings.TrimSpace(values["Deadline (YYYY-MM-DD):"])
	if deadlineStr != "" {
		if d, err := time.Parse("2006-01-02", deadlineStr); err == nil {
			deadline = &d
		}
	}

	domain.UpdateQuest(&m.projects, m.selectedProjectIdx, m.editingIdx, title, desc, priority, deadline)
	m.updateScreenModels()
}

func (m *RootModel) createTask() {
	if m.selectedProjectIdx < 0 || m.selectedQuestIdx < 0 {
		return
	}
	values := m.form.GetValues()
	desc := strings.TrimSpace(values["Description:"])
	if desc == "" {
		return
	}

	domain.CreateTask(&m.projects, m.selectedProjectIdx, m.selectedQuestIdx, desc)
	m.updateScreenModels()
}

func (m *RootModel) updateTask() {
	if m.selectedProjectIdx < 0 || m.selectedQuestIdx < 0 || m.editingIdx < 0 {
		return
	}
	values := m.form.GetValues()
	desc := strings.TrimSpace(values["Description:"])
	if desc == "" {
		return
	}

	domain.UpdateTask(&m.projects, m.selectedProjectIdx, m.selectedQuestIdx, m.editingIdx, desc)
	m.updateScreenModels()
}

func (m *RootModel) startDeleteTask() {
	idx := m.taskList.SelectedTaskIndex()
	if idx >= 0 && m.selectedProjectIdx >= 0 && m.selectedQuestIdx >= 0 &&
		idx < len(m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks) {
		m.pendingDelete = true
		m.deleteType = "task"
		m.deleteIndices = [3]int{m.selectedProjectIdx, m.selectedQuestIdx, idx}
	}
}

func (m *RootModel) toggleTask() {
	if m.selectedProjectIdx < 0 || m.selectedQuestIdx < 0 {
		return
	}
	taskIdx := m.taskList.SelectedTaskIndex()
	if taskIdx >= 0 && taskIdx < len(m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks) {
		m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks[taskIdx].Done =
			!m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks[taskIdx].Done
		m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].CalculateProgress()
		m.projects[m.selectedProjectIdx].CalculateProgress()
		m.updateScreenModels()
	}
}

func (m *RootModel) saveProjectsCmd() tea.Cmd {
	return func() tea.Msg {
		err := domain.SaveProjects(m.projects)
		return SaveCompleteMsg{Err: err}
	}
}

func (m *RootModel) saveProjects() {
	_ = domain.SaveProjects(m.projects)
}
