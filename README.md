# Quest Line

A terminal-based quest and task management application built with Go and Bubble Tea.

## Installation

```bash
# Clone or download the repository
# Build the application
go build -o quest_line .

# Run
./quest_line
```

## Usage
<img width="1384" height="738" alt="1" src="https://github.com/user-attachments/assets/f88f2bb9-81a1-47e7-9c0a-aad0da442135" />

### Navigation
- `d` - Dashboard (active quests overview)
- `p` - Projects list
- `q` - Quit

### Lists (Projects/Quests)
- `↑/k` - Move up
- `↓/j` - Move down
- `Enter` - Select/open
- `c` - Create new
- `e` - Edit selected
- `x` - Delete selected

### Quest Details
- `↑/k` - Navigate tasks
- `↓/j` - Navigate tasks
- `Enter` - Toggle task completion
- `c` - Create task
- `e` - Edit task
- `x` - Delete task
- `d` - Back to dashboard
<img width="1381" height="736" alt="2" src="https://github.com/user-attachments/assets/522c7218-0695-4b09-8745-946336a4c23d" />

### Forms
- `Tab` - Next field
- `Enter` - Submit (advances through fields)
- `Esc` - Cancel
<img width="1381" height="739" alt="3" src="https://github.com/user-attachments/assets/fe991363-8999-4d3e-b4c1-7a2c3133f859" />

## Features

- **Project Management**: Organize quests into projects
- **Quest Tracking**: Create and manage quests with priorities and deadlines
- **Task Management**: Break quests into actionable tasks
- **Progress Tracking**: Automatic progress calculation
- **Dashboard**: Daily overview of active quests
- **Persistent Storage**: JSON-based data storage
- **Keyboard-Driven**: Full keyboard navigation

## Data Storage

Data is stored in `quests.json`:

```json
[
  {
    "ID": "sample-project",
    "Name": "Sample Project",
    "Quests": [],
    "Progress": 50
  }
]
```

## Development

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

## Quick Start

1. Launch the app: `./quest_line`
2. Select a project (or create one with `c`)
3. Create quests and tasks
4. Use the dashboard to track progress

Happy questing! 🚀
