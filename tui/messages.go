package tui

import "quest_line/domain"

// DataChangedMsg is sent when data changes (for persistence)
type DataChangedMsg struct {
	Projects []domain.Project
}

// ItemToggleMsg is sent when a task/item is toggled
type ItemToggleMsg struct {
	ProjectIdx int
	QuestIdx   int
	TaskIdx    int
}

// ItemDeleteMsg is sent when an item is deleted
type ItemDeleteMsg struct {
	ItemType   string // "project", "quest", "task"
	ProjectIdx int
	QuestIdx   int
	TaskIdx    int
}

// SaveCompleteMsg is sent when save is complete
type SaveCompleteMsg struct {
	Err error
}
