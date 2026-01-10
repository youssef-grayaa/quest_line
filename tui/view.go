package tui

import (
	"fmt"
)

// View renders the appropriate screen
func (m *RootModel) View() string {
	if m.pendingDelete {
		return appStyle.Render(m.viewDeleteConfirmation())
	}

	switch m.currentView {
	case ViewProjectSelection:
		return appStyle.Render(m.projectSelection.View() + "\n\n" + m.help.ViewFor(m.currentView))
	case ViewDashboard:
		return appStyle.Render(m.dashboard.View() + "\n\n" + m.help.ViewFor(m.currentView))
	case ViewProjectList:
		return appStyle.Render(m.projectList.View() + "\n\n" + m.help.ViewFor(m.currentView))
	case ViewQuestDetail:
		return appStyle.Render(m.taskList.View() + "\n\n" + m.help.ViewFor(m.currentView))
	case ViewCreateProject, ViewEditProject, ViewCreateQuest, ViewEditQuest, ViewCreateTask, ViewEditTask:
		return appStyle.Render(m.form.View() + "\n\n" + m.help.ViewFor(m.currentView))
	default:
		return appStyle.Render("Unknown view")
	}
}

// viewDeleteConfirmation shows the delete confirmation prompt
func (m *RootModel) viewDeleteConfirmation() string {
	var itemName string
	switch m.deleteType {
	case "project":
		if m.deleteIndices[0] >= 0 && m.deleteIndices[0] < len(m.projects) {
			itemName = m.projects[m.deleteIndices[0]].Name
		}
	case "quest":
		if m.selectedProjectIdx >= 0 && m.deleteIndices[1] >= 0 &&
			m.selectedProjectIdx < len(m.projects) &&
			m.deleteIndices[1] < len(m.projects[m.selectedProjectIdx].Quests) {
			itemName = m.projects[m.selectedProjectIdx].Quests[m.deleteIndices[1]].Title
		}
	case "task":
		if m.selectedProjectIdx >= 0 && m.selectedQuestIdx >= 0 && m.deleteIndices[2] >= 0 &&
			m.selectedProjectIdx < len(m.projects) &&
			m.selectedQuestIdx < len(m.projects[m.selectedProjectIdx].Quests) &&
			m.deleteIndices[2] < len(m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks) {
			itemName = m.projects[m.selectedProjectIdx].Quests[m.selectedQuestIdx].Tasks[m.deleteIndices[2]].Description
		}
	}

	return fmt.Sprintf("Delete %s '%s'? (y/n)", m.deleteType, itemName)
}
