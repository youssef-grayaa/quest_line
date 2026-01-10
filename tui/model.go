package tui

import (
	"quest_line/domain"
)

// RootModel is the main application model
type RootModel struct {
	// Navigation
	currentView View
	keymap      KeyMap
	help        HelpModel

	// Data
	projects []domain.Project

	// Screen models
	projectSelection ProjectSelectionModel
	dashboard        DashboardModel
	projectList      ProjectListModel
	taskList         QuestDetailModel
	form             FormModel
	inForm           bool

	// List keys
	listKeys     *listKeyMap
	delegateKeys *delegateKeyMap

	// Navigation context (using indices for safety)
	selectedProjectIdx int
	selectedQuestIdx   int
	editingIdx         int // for edit operations

	// Delete confirmation
	pendingDelete bool
	deleteType    string // "project", "quest", "task"
	deleteIndices [3]int // projectIdx, questIdx, taskIdx

	// Error handling
	errorMsg string
}

// InitialModel creates the initial root model
func InitialModel() RootModel {
	projects, _ := domain.LoadProjects()

	// If no projects loaded, add a sample project
	if len(projects) == 0 {
		sampleProject := domain.Project{
			ID:       "sample-project",
			Name:     "Sample Project",
			Quests:   []domain.Quest{},
			Progress: 0,
		}
		projects = []domain.Project{sampleProject}
		// Save the initial sample project
		_ = domain.SaveProjects(projects)
	}

	// Calculate progress for all quests and projects (in case loaded from JSON without progress)
	for i := range projects {
		for j := range projects[i].Quests {
			projects[i].Quests[j].CalculateProgress()
		}
		projects[i].CalculateProgress()
	}

	keymap := DefaultKeyMap()
	help := NewHelpModel()
	listKeys := newListKeyMap()
	delegateKeys := newDelegateKeyMap()

	// Determine initial view and selected project
	var currentView View
	var selectedProjectIdx int
	projectSelection := NewProjectSelectionModel(projects, keymap)

	if len(projects) == 0 {
		// Shouldn't happen, but handle
		currentView = ViewDashboard
		selectedProjectIdx = -1
	} else if len(projects) == 1 {
		currentView = ViewDashboard
		selectedProjectIdx = 0
	} else {
		currentView = ViewProjectSelection
		selectedProjectIdx = -1
	}

	// Create screen models
	dashboard := NewDashboardModel(projects, selectedProjectIdx, keymap)
	projectList := NewProjectListModel(projects, keymap)
	taskList := NewQuestDetailModel(projects, -1, -1, keymap)

	return RootModel{
		currentView:        currentView,
		keymap:             keymap,
		help:               help,
		projects:           projects,
		projectSelection:   projectSelection,
		dashboard:          dashboard,
		projectList:        projectList,
		taskList:           taskList,
		listKeys:           listKeys,
		delegateKeys:       delegateKeys,
		inForm:             false,
		selectedProjectIdx: selectedProjectIdx,
		selectedQuestIdx:   -1,
		editingIdx:         -1,
		pendingDelete:      false,
	}
}
