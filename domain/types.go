package domain

import "time"

type QuestState int

const (
	StateActive QuestState = iota
	StateCompleted
	StateCancelled
)

func (s QuestState) String() string {
	switch s {
	case StateActive:
		return "Active"
	case StateCompleted:
		return "Completed"
	case StateCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

type Task struct {
	ID          string
	Description string
	Done        bool
}

type Quest struct {
	ID          string
	Title       string
	Description string
	Tasks       []Task
	Progress    float64 // 0.0 → 100.0

	Priority int
	Deadline *time.Time
	State    QuestState
}

type Project struct {
	ID       string
	Name     string
	Quests   []Quest
	Progress float64 // 0.0 → 100.0
}
