package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all keybindings for the application
type KeyMap struct {
	Dashboard key.Binding
	Projects  key.Binding
	QuestList key.Binding
	Up        key.Binding
	Down      key.Binding
	Enter     key.Binding
	Create    key.Binding
	Edit      key.Binding
	Delete    key.Binding
	Toggle    key.Binding
	Tab       key.Binding
	ShiftTab  key.Binding
	Submit    key.Binding
	Cancel    key.Binding
	Help      key.Binding
	Quit      key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Dashboard: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "dashboard"),
		),
		Projects: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "projects"),
		),
		QuestList: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "quest list"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/submit"),
		),
		Create: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "delete"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev field"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Help: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns the short help info for a view
func (k KeyMap) ShortHelpForView(view View) []key.Binding {
	switch view {
	case ViewProjectSelection:
		return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
	case ViewDashboard:
		createQuestKey := key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "create quest"))
		return []key.Binding{k.Up, k.Down, k.Enter, createQuestKey, k.Edit, k.Delete, k.Help, k.Quit}
	case ViewProjectList:
		return []key.Binding{k.Up, k.Down, k.Enter, k.Create, k.Edit, k.Delete, k.Dashboard, k.Help, k.Quit}
	case ViewQuestDetail:
		return []key.Binding{k.Up, k.Down, k.Enter, k.Create, k.Edit, k.Delete, k.Dashboard, k.Help, k.Quit}
	default:
		return []key.Binding{k.Help, k.Quit}
	}
}

// ShortHelp returns the default short help (for backward compatibility)
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns the full help info
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Create, k.Edit, k.Delete},
		{k.Toggle, k.Dashboard, k.Projects, k.QuestList},
		{k.Tab, k.ShiftTab, k.Submit, k.Cancel},
		{k.Help, k.Quit},
	}
}

// HelpModel wraps the help component
type HelpModel struct {
	help help.Model
	keys KeyMap
}

// NewHelpModel creates a new help model
func NewHelpModel() HelpModel {
	hm := HelpModel{
		help: help.New(),
		keys: DefaultKeyMap(),
	}
	hm.help.ShowAll = false
	return hm
}

// Update updates the help model
func (h HelpModel) Update(msg interface{}) HelpModel {
	// Help model doesn't handle messages
	return h
}

// viewKeyMap wraps KeyMap with view-specific ShortHelp
type viewKeyMap struct {
	KeyMap
	view View
}

// ShortHelp returns view-specific short help
func (v viewKeyMap) ShortHelp() []key.Binding {
	return v.KeyMap.ShortHelpForView(v.view)
}

// View renders the help text for a view
func (h HelpModel) ViewFor(view View) string {
	vkm := viewKeyMap{KeyMap: h.keys, view: view}
	return h.help.View(vkm)
}

// View renders the default help text
func (h HelpModel) View() string {
	return h.help.View(h.keys)
}

// ToggleHelp toggles between short and full help
func (h *HelpModel) ToggleHelp() {
	h.help.ShowAll = !h.help.ShowAll
}
