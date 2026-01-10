package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"quest_line/tui"
)

func main() {
	model := tui.InitialModel()
	program := tea.NewProgram(&model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
