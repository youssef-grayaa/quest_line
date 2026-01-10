package domain

import (
	"fmt"
	"sort"
	"time"
)

// CalculateProgress computes the progress of a quest based on completed tasks
func (q *Quest) CalculateProgress() {
	if len(q.Tasks) == 0 {
		q.Progress = 0.0
		return
	}
	completed := 0
	for _, task := range q.Tasks {
		if task.Done {
			completed++
		}
	}
	q.Progress = float64(completed) / float64(len(q.Tasks)) * 100.0
}

// CalculateProgress computes the progress of a project based on quest progresses
func (p *Project) CalculateProgress() {
	if len(p.Quests) == 0 {
		p.Progress = 0.0
		return
	}
	total := 0.0
	for _, quest := range p.Quests {
		total += quest.Progress
	}
	p.Progress = total / float64(len(p.Quests))
}

// IsLeaf returns true if the quest has no sub-quests (assuming no nested quests for now)
func (q *Quest) IsLeaf() bool {
	// Since we don't have sub-quests in the current model, all quests are leaf
	return true
}

// DailyPlanner returns leaf quests ordered by priority (higher first) and deadlines (soonest first)
func DailyPlanner(projects []Project) []Quest {
	var leafQuests []Quest
	for _, project := range projects {
		for _, quest := range project.Quests {
			if quest.IsLeaf() && quest.State == StateActive {
				leafQuests = append(leafQuests, quest)
			}
		}
	}

	sort.Slice(leafQuests, func(i, j int) bool {
		// Higher priority first
		if leafQuests[i].Priority != leafQuests[j].Priority {
			return leafQuests[i].Priority > leafQuests[j].Priority
		}
		// Soonest deadline first
		if leafQuests[i].Deadline != nil && leafQuests[j].Deadline != nil {
			return leafQuests[i].Deadline.Before(*leafQuests[j].Deadline)
		}
		if leafQuests[i].Deadline != nil {
			return true
		}
		if leafQuests[j].Deadline != nil {
			return false
		}
		return false
	})

	return leafQuests
}

// DailyPlannerForProject returns active quests for a specific project
func DailyPlannerForProject(projects []Project, projectIdx int) []Quest {
	if projectIdx < 0 || projectIdx >= len(projects) {
		return []Quest{}
	}
	var leafQuests []Quest
	for _, quest := range projects[projectIdx].Quests {
		if quest.IsLeaf() && quest.State == StateActive {
			leafQuests = append(leafQuests, quest)
		}
	}

	sort.Slice(leafQuests, func(i, j int) bool {
		// Higher priority first
		if leafQuests[i].Priority != leafQuests[j].Priority {
			return leafQuests[i].Priority > leafQuests[j].Priority
		}
		// Soonest deadline first
		if leafQuests[i].Deadline != nil && leafQuests[j].Deadline != nil {
			return leafQuests[i].Deadline.Before(*leafQuests[j].Deadline)
		}
		if leafQuests[i].Deadline != nil {
			return true
		}
		if leafQuests[j].Deadline != nil {
			return false
		}
		return false
	})

	return leafQuests
}

// generateID generates a unique ID using timestamp
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// CreateProject adds a new project to the list
func CreateProject(projects *[]Project, name string) *Project {
	id := generateID()
	p := Project{ID: id, Name: name}
	*projects = append(*projects, p)
	return &(*projects)[len(*projects)-1]
}

// UpdateProject updates an existing project
func UpdateProject(projects *[]Project, index int, name string) {
	if index >= 0 && index < len(*projects) {
		(*projects)[index].Name = name
	}
}

// DeleteProject removes a project and its quests
func DeleteProject(projects *[]Project, index int) {
	if index >= 0 && index < len(*projects) {
		*projects = append((*projects)[:index], (*projects)[index+1:]...)
	}
}

// CreateQuest adds a new quest to a project
func CreateQuest(projects *[]Project, projectIndex int, title, description string, priority int, deadline *time.Time) *Quest {
	if projectIndex < 0 || projectIndex >= len(*projects) {
		return nil
	}
	id := generateID()
	q := Quest{
		ID:          id,
		Title:       title,
		Description: description,
		Priority:    priority,
		Deadline:    deadline,
		State:       StateActive,
	}
	(*projects)[projectIndex].Quests = append((*projects)[projectIndex].Quests, q)
	return &(*projects)[projectIndex].Quests[len((*projects)[projectIndex].Quests)-1]
}

// UpdateQuest updates an existing quest
func UpdateQuest(projects *[]Project, projectIndex, questIndex int, title, description string, priority int, deadline *time.Time) {
	if projectIndex >= 0 && projectIndex < len(*projects) && questIndex >= 0 && questIndex < len((*projects)[projectIndex].Quests) {
		q := &(*projects)[projectIndex].Quests[questIndex]
		q.Title = title
		q.Description = description
		q.Priority = priority
		q.Deadline = deadline
		q.CalculateProgress()
	}
}

// DeleteQuest removes a quest from a project
func DeleteQuest(projects *[]Project, projectIndex, questIndex int) {
	if projectIndex >= 0 && projectIndex < len(*projects) && questIndex >= 0 && questIndex < len((*projects)[projectIndex].Quests) {
		(*projects)[projectIndex].Quests = append((*projects)[projectIndex].Quests[:questIndex], (*projects)[projectIndex].Quests[questIndex+1:]...)
	}
}

// CreateTask adds a new task to a quest
func CreateTask(projects *[]Project, projectIndex, questIndex int, description string) *Task {
	if projectIndex >= 0 && projectIndex < len(*projects) && questIndex >= 0 && questIndex < len((*projects)[projectIndex].Quests) {
		id := generateID()
		t := Task{ID: id, Description: description, Done: false}
		(*projects)[projectIndex].Quests[questIndex].Tasks = append((*projects)[projectIndex].Quests[questIndex].Tasks, t)
		(*projects)[projectIndex].Quests[questIndex].CalculateProgress()
		return &(*projects)[projectIndex].Quests[questIndex].Tasks[len((*projects)[projectIndex].Quests[questIndex].Tasks)-1]
	}
	return nil
}

// UpdateTask updates an existing task
func UpdateTask(projects *[]Project, projectIndex, questIndex, taskIndex int, description string) {
	if projectIndex >= 0 && projectIndex < len(*projects) && questIndex >= 0 && questIndex < len((*projects)[projectIndex].Quests) && taskIndex >= 0 && taskIndex < len((*projects)[projectIndex].Quests[questIndex].Tasks) {
		(*projects)[projectIndex].Quests[questIndex].Tasks[taskIndex].Description = description
		(*projects)[projectIndex].Quests[questIndex].CalculateProgress()
	}
}

// DeleteTask removes a task from a quest
func DeleteTask(projects *[]Project, projectIndex, questIndex, taskIndex int) {
	if projectIndex >= 0 && projectIndex < len(*projects) && questIndex >= 0 && questIndex < len((*projects)[projectIndex].Quests) && taskIndex >= 0 && taskIndex < len((*projects)[projectIndex].Quests[questIndex].Tasks) {
		(*projects)[projectIndex].Quests[questIndex].Tasks = append((*projects)[projectIndex].Quests[questIndex].Tasks[:taskIndex], (*projects)[projectIndex].Quests[questIndex].Tasks[taskIndex+1:]...)
		(*projects)[projectIndex].Quests[questIndex].CalculateProgress()
	}
}

// FindQuestIndices finds the project and quest indices for a given quest ID
func FindQuestIndices(projects []Project, questID string) (int, int) {
	for pIdx, project := range projects {
		for qIdx, quest := range project.Quests {
			if quest.ID == questID {
				return pIdx, qIdx
			}
		}
	}
	return -1, -1
}
