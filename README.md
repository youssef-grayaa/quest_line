# Quick Start - Quest Line TUI

## Building and Running

```bash
# Build the application
go build -o quest_line .

# Run the application
./quest_line
```

## Key Controls

### Global Shortcuts

| Key | Action |
|-----|--------|
| `d` | Go to Dashboard |
| `p` | Go to Projects |
| `l` | Go to Quest List |
| `?` | Toggle Help |
| `q` | Quit |

### Navigation (Project/Quest Lists)

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Open/Select |
| `c` | Create new |
| `e` | Edit selected |
| `x` | Delete selected |

### Quest Detail View

| Key | Action |
|-----|--------|
| `↑` / `k` | Move to previous task |
| `↓` / `j` | Move to next task |
| `Space` | Toggle task completion |
| `c` | Create new task |
| `e` | Edit selected task |
| `x` | Delete selected task |
| `l` | Back to quest list |

### Form Input

| Key | Action |
|-----|--------|
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Enter` | Submit form (when on last field) |
| `Esc` | Cancel form |

## File Structure

```
quest_line/
├── main.go                 # Application entry point
├── domain/
│   ├── types.go           # Project, Quest, Task types
│   ├── quest.go           # Domain logic (CRUD, calculations)
│   └── persistence.go     # JSON save/load
└── tui/
    ├── root.go            # RootModel (main coordinator)
    ├── screens.go         # Screen models (Dashboard, ProjectList, etc.)
    ├── form.go            # FormModel (input handling)
    ├── keybindings.go     # KeyMap and help integration
    ├── messages.go        # Custom message types
    └── [deprecated]
        ├── model.go       # Documentation only
        ├── update.go      # Documentation only
        └── view.go        # Documentation only

quests.json               # Data file (auto-created)
ARCHITECTURE.md           # Full architecture documentation
REFACTORING.md           # Before/after comparison
IMPLEMENTATION.md        # Detailed implementation guide
README.md                # This file
```

## Architecture at a Glance

### State Management

```go
RootModel
├── projects []domain.Project        // All data
├── currentView View                  // Current screen
├── dashboard DashboardModel         // Dashboard state
├── projectList ProjectListModel     // Project list state
├── questList QuestListModel         // Quest list state
├── questDetail QuestDetailModel     // Quest detail state
└── form FormModel                    // Form state
```

### Message Flow

```
User Input (tea.KeyMsg)
    ↓
RootModel.Update()
    ├─ Global handlers (quit, help, navigate)
    └─ View-specific handlers (up, down, create, etc.)
    ↓
Domain Operations (if needed)
    ├─ CreateProject, UpdateProject, DeleteProject
    ├─ CreateQuest, UpdateQuest, DeleteQuest
    └─ CreateTask, UpdateTask, DeleteTask
    ↓
Refresh Models
    ├─ dashboard.Update(DataChangedMsg)
    ├─ projectList.Update(DataChangedMsg)
    ├─ questList.Update(DataChangedMsg)
    └─ questDetail.Update(DataChangedMsg)
    ↓
Save Data (async via tea.Cmd)
    └─ SaveCompleteMsg (completion notification)
```

### Screen Navigation

```
Dashboard
├─ p → ProjectList
│     ├─ Enter → QuestList (filtered by project)
│     │         ├─ Enter → QuestDetail
│     │         │         ├─ c → CreateTask
│     │         │         ├─ e → EditTask
│     │         │         └─ x → Delete Task
│     │         ├─ c → CreateQuest
│     │         ├─ e → EditQuest
│     │         └─ x → Delete Quest
│     ├─ c → CreateProject
│     ├─ e → EditProject
│     └─ x → Delete Project
└─ l → QuestList (all quests)
        └─ [same as above]
```

## Common Tasks

### Create a Project

1. Press `p` to go to Projects
2. Press `c` to create
3. Enter project name
4. Press `Enter` to submit

### Add a Quest to a Project

1. Press `p` to go to Projects
2. Press `Enter` to open a project
3. Press `c` to create a quest
4. Fill in:
   - Title (required)
   - Description
   - Priority (0-10)
   - Deadline (YYYY-MM-DD, optional)
5. Press `Enter` on the last field to submit

### Add Tasks to a Quest

1. Open a project (press `Enter` on project)
2. Open a quest (press `Enter` on quest)
3. Press `c` to create a task
4. Enter task description
5. Press `Enter` to submit

### Mark a Task as Done

1. Open a quest (see above)
2. Navigate to a task (↑/↓ keys)
3. Press `Space` to toggle completion

### View Today's Active Quests

1. Press `d` to go to Dashboard
2. See all active quests sorted by priority and deadline

## Features

- ✅ Multi-level CRUD operations (Projects → Quests → Tasks)
- ✅ Task completion tracking with progress calculation
- ✅ Project-based quest organization
- ✅ Quest state management (Active, Completed, Cancelled)
- ✅ Priority and deadline support
- ✅ Daily planner view (Dashboard)
- ✅ Persistent JSON storage
- ✅ Full keyboard navigation
- ✅ Integrated help system
- ✅ Color-coded UI with selection indicators

## Limitations

- No synchronization with external services
- No nested quests (flat hierarchy only)
- No categories/tags for quests
- No recurring tasks
- Single user only (no authentication/multi-user)
- No undo/redo functionality
- No search across all items

## Data Format

Projects are stored in `quests.json`:

```json
[
  {
    "ID": "1234567890",
    "Name": "My Project",
    "Quests": [
      {
        "ID": "1234567891",
        "Title": "My Quest",
        "Description": "Do something important",
        "Tasks": [
          {
            "ID": "1234567892",
            "Description": "Subtask 1",
            "Done": false
          }
        ],
        "Progress": 0.0,
        "Priority": 5,
        "Deadline": "2026-01-31T00:00:00Z",
        "State": 0
      }
    ]
  }
]
```

## Development Tips

### Understanding the Code

1. Start with `main.go` - simple entry point
2. Read `tui/root.go` - main coordinator
3. Read `tui/screens.go` - individual screen models
4. Read `tui/form.go` - form handling
5. Check `domain/` for business logic

### Extending the App

**Add a new screen:**
1. Create a new model struct in `screens.go`
2. Add it to `RootModel` in `root.go`
3. Add navigation in `handleViewSpecificInput()`
4. Implement `Update()` and `View()` methods

**Add a new message type:**
1. Define in `messages.go`
2. Handle in `RootModel.Update()`
3. Update screen models to handle it

**Add new keybindings:**
1. Add to `KeyMap` in `keybindings.go`
2. Include help text in binding creation
3. Handle in appropriate input method

### Debugging

Print debug info to stderr (doesn't interfere with TUI):

```go
import "log"

// In any method
log.Printf("Debug: selectedIdx=%d", m.selectedIdx)
```

Run with:
```bash
./quest_line 2> debug.log
```

View debug output:
```bash
cat debug.log
```

## Performance Notes

- All data loaded into memory at startup
- No pagination (suitable for ~1000 items)
- Synchronous JSON save (blocking)
- No caching/indexing
- For larger datasets, consider:
  - Database instead of JSON
  - Async persistence
  - Pagination UI
  - Search/filter indexing

## Troubleshooting

**Application won't start:**
- Check Go version (requires Go 1.19+)
- Verify dependencies: `go mod tidy`
- Check for terminal issues: `echo $TERM`

**Data not saving:**
- Check file permissions on `quests.json`
- Ensure write access to directory
- Check disk space
- Look for error messages (run with stderr)

**Keyboard shortcuts not working:**
- Some terminals remap keys (check terminal settings)
- Try alternative key combos:
  - `up/down` instead of `k/j`
  - `ctrl+c` instead of `q` for quit

**UI glitches:**
- Resize terminal
- Export `TERM=xterm-256color`
- Try different terminal emulator

## Next Steps

- Read `ARCHITECTURE.md` for detailed design
- Read `IMPLEMENTATION.md` for code patterns
- Explore `examples/` in bubbletea repo for inspiration
- Check `domain/quest.go` for business logic patterns

## Resources

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubbles Documentation](https://github.com/charmbracelet/bubbles)
- [Lipgloss for Styling](https://github.com/charmbracelet/lipgloss)
- [Elm Architecture Guide](https://guide.elm-lang.org/architecture/)

Happy questing!
