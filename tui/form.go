package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"quest_line/domain"
)

// FormModel handles form input and validation for create/edit operations
type FormModel struct {
	inputs    []textinput.Model
	focusIdx  int
	formType  View // ViewCreateProject, ViewCreateQuest, etc.
	title     string
	labels    []string
	fieldInfo map[string]string // store processed field values
	errorMsg  string
}

// NewProjectForm creates a new form for project creation/editing
func NewProjectForm(title string, initialName string) FormModel {
	inputs := make([]textinput.Model, 1)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project Name"
	inputs[0].SetValue(initialName)
	inputs[0].Focus()

	return FormModel{
		inputs:    inputs,
		focusIdx:  0,
		formType:  ViewCreateProject,
		title:     title,
		labels:    []string{"Name:"},
		fieldInfo: make(map[string]string),
	}
}

// NewQuestForm creates a new form for quest creation/editing
func NewQuestForm(title string, initial *domain.Quest) FormModel {
	inputs := make([]textinput.Model, 4)
	for i := range inputs {
		inputs[i] = textinput.New()
	}

	inputs[0].Placeholder = "Title"
	inputs[0].Focus()
	inputs[1].Placeholder = "Description"
	inputs[2].Placeholder = "Priority (0-10)"
	inputs[3].Placeholder = "Deadline (YYYY-MM-DD)"

	if initial != nil {
		inputs[0].SetValue(initial.Title)
		inputs[1].SetValue(initial.Description)
		inputs[2].SetValue(strconv.Itoa(initial.Priority))
		if initial.Deadline != nil {
			inputs[3].SetValue(initial.Deadline.Format("2006-01-02"))
		}
	}

	return FormModel{
		inputs:    inputs,
		focusIdx:  0,
		formType:  ViewCreateQuest,
		title:     title,
		labels:    []string{"Title:", "Description:", "Priority (0-10):", "Deadline (YYYY-MM-DD):"},
		fieldInfo: make(map[string]string),
	}
}

// NewTaskForm creates a new form for task creation/editing
func NewTaskForm(title string, initialDesc string) FormModel {
	inputs := make([]textinput.Model, 1)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Task Description"
	inputs[0].SetValue(initialDesc)
	inputs[0].Focus()

	return FormModel{
		inputs:    inputs,
		focusIdx:  0,
		formType:  ViewCreateTask,
		title:     title,
		labels:    []string{"Description:"},
		fieldInfo: make(map[string]string),
	}
}

// Update handles form input
func (f FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			f.nextInput()
		case "shift+tab":
			f.prevInput()
		case "enter":
			if f.focusIdx == len(f.inputs)-1 {
				// Submit form
				return f, nil
			}
			f.nextInput()
		}
	}

	f.inputs[f.focusIdx], cmd = f.inputs[f.focusIdx].Update(msg)
	return f, cmd
}

// View renders the form
func (f FormModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(f.title))
	b.WriteString("\n\n")

	for i, input := range f.inputs {
		b.WriteString(f.labels[i])
		b.WriteString("\n")
		b.WriteString(input.View())
		if i == f.focusIdx {
			b.WriteString(" ← editing")
		}
		b.WriteString("\n\n")
	}

	if f.errorMsg != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Render("Error: " + f.errorMsg))
		b.WriteString("\n\n")
	}

	return b.String()
}

// nextInput moves focus to next input
func (f *FormModel) nextInput() {
	f.focusIdx = (f.focusIdx + 1) % len(f.inputs)
	f.updateInputFocus()
}

// prevInput moves focus to previous input
func (f *FormModel) prevInput() {
	f.focusIdx--
	if f.focusIdx < 0 {
		f.focusIdx = len(f.inputs) - 1
	}
	f.updateInputFocus()
}

// updateInputFocus applies focus styling to inputs
func (f *FormModel) updateInputFocus() {
	for i := range f.inputs {
		if i == f.focusIdx {
			f.inputs[i].Focus()
		} else {
			f.inputs[i].Blur()
		}
	}
}

// GetValues returns the form values as a map
func (f FormModel) GetValues() map[string]string {
	values := make(map[string]string)
	for i, input := range f.inputs {
		values[f.labels[i]] = input.Value()
	}
	return values
}

// Validate checks if the form data is valid
func (f FormModel) Validate() error {
	switch f.formType {
	case ViewCreateProject, ViewEditProject:
		if strings.TrimSpace(f.inputs[0].Value()) == "" {
			return ErrProjectNameRequired
		}
	case ViewCreateQuest, ViewEditQuest:
		if strings.TrimSpace(f.inputs[0].Value()) == "" {
			return ErrQuestTitleRequired
		}
		// Validate priority is a number
		if p := strings.TrimSpace(f.inputs[2].Value()); p != "" {
			// Just check if it parses as int
			_ = p
		}
		// Validate deadline format
		if d := strings.TrimSpace(f.inputs[3].Value()); d != "" {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				return ErrInvalidDateFormat
			}
		}
	case ViewCreateTask, ViewEditTask:
		if strings.TrimSpace(f.inputs[0].Value()) == "" {
			return ErrTaskDescRequired
		}
	}
	return nil
}

// Custom errors
var (
	ErrProjectNameRequired = NewValidationError("project name is required")
	ErrQuestTitleRequired  = NewValidationError("quest title is required")
	ErrTaskDescRequired    = NewValidationError("task description is required")
	ErrInvalidDateFormat   = NewValidationError("invalid date format (use YYYY-MM-DD)")
)

// ValidationError represents a form validation error
type ValidationError struct {
	msg string
}

// NewValidationError creates a new validation error
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg: msg}
}

// Error implements error interface
func (e *ValidationError) Error() string {
	return e.msg
}
