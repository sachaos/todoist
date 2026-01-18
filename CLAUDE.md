# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a CLI client for Todoist written in Go (requires Go 1.25+). It provides command-line access to Todoist's task management features with support for filtering, projects, labels, and offline caching.

**Note**: This project uses Todoist's Sync API v1 (migrated from v9). The v9 API is scheduled to shut down on Feb 10, 2026.

## Common Commands

### Building and Testing
```bash
# Install the binary to $GOPATH/bin
make install

# Build the binary in the current directory
make build

# Run all tests
make test

# Generate the filter parser (required before building if filter_parser.y is modified)
make prepare
```

### Running Single Tests
```bash
# Run a specific test function
go test -v -run TestFilterParser

# Run tests in a specific file
go test -v ./filter_parser_test.go ./filter_parser.go ./filter_eval.go
```

### Development Workflow
```bash
# Generate filter parser from yacc grammar (automatically done by prepare target)
goyacc -o filter_parser.go filter_parser.y

# Build for multiple platforms (for releases)
make release VERSION=x.y.z
```

## Architecture

### Core Structure

The codebase is organized into two main layers:

1. **CLI Layer (root directory)**: Command handlers, formatting, and user-facing logic
   - `main.go` - CLI setup using urfave/cli/v2 framework, app initialization
   - `list.go` - List tasks with tree traversal and filter support
   - `add.go` - Add new task with priority, labels, project, date
   - `modify.go` - Modify existing task properties
   - `close.go` - Mark task(s) complete
   - `delete.go` - Delete task(s) with prefix completion support
   - `show.go` - Show task details with URL opening (pkg/browser)
   - `quick.go` - Quick add task using REST endpoint
   - `completed.go` - List completed tasks (premium, 90-day default range)
   - `labels.go`, `projects.go`, `add_project.go` - Label/project management
   - `karma.go` - Display user karma
   - `sync.go` - Sync cache with Todoist API
   - `format.go` - Output formatting (colors, dates, priorities)
   - `cache.go` - Local cache management for offline access
   - `filter_parser.y` - Yacc grammar (generates `filter_parser.go`)
   - `filter_eval.go` - Filter expression evaluation
   - `utils.go` - TSVWriter and file utilities

2. **Library Layer (`lib/` directory)**: Core Todoist API client and data models
   - `lib/todoist.go` - HTTP client and API communication (base URL: `api.todoist.com/api/v1/`)
   - `lib/sync.go` - Store structure, tree construction, lookup maps
   - `lib/item.go` - Item model with tree pointers, date handling
   - `lib/project.go` - Project model with tree structure
   - `lib/label.go`, `lib/section.go` - Label and section models
   - `lib/user.go` - User model with profile info
   - `lib/completed.go` - Completed tasks API (90-day range)
   - `lib/command.go` - Command serialization for Sync API
   - `lib/interface.go` - Common interfaces (HaveID, HaveProjectID, etc.)
   - `lib/item_order.go` - Order sorting structures

### Data Model

The application uses a **tree-based data structure** for both items and projects:

- **Store** (`lib/sync.go`): Central data structure containing all synced data
  - Maintains maps (`ItemMap`, `ProjectMap`, `LabelMap`, `SectionMap`) for fast lookups
  - Constructs linked tree structures with `ChildItem`/`BrotherItem` and `ChildProject`/`BrotherProject` pointers
  - `RootItem` and `RootProject` serve as entry points to traverse hierarchies

- **Tree Construction**: `ConstructItemTree()` builds parent-child and sibling relationships from flat API responses

### Filter System

The filter implementation uses a **yacc-based parser** for Todoist's filter syntax:

- `filter_parser.y` - Yacc grammar definition (generates `filter_parser.go`)
- `filter_eval.go` - Expression evaluation logic
- Supports date expressions, project/label filters, boolean operators (`&`, `|`, `!`)
- Filter syntax examples: `(overdue | today) & p1`, `#ProjectName`, `@LabelName`

### Configuration and Caching

- **Config**: Stored in `$HOME/.config/todoist/config.json` (XDG Base Directory compliant)
  - Contains API token, color preferences, date formats
  - Permissions enforced to 0600 for security

- **Cache**: Stored in `$HOME/.cache/todoist/cache.json`
  - Full sync from Todoist API using sync token
  - Enables offline browsing and faster startup

### CLI Commands

Available subcommands (via urfave/cli/v2):
- `list` / `l` - Show tasks (supports filter, priority sorting)
- `show` - Show task detail (with `--browse` flag to open URLs)
- `completed-list` / `c-l` / `cl` - List completed tasks (premium only, 90-day default)
- `add` / `a` - Add task (priority, labels, project, date, reminder)
- `modify` / `m` - Edit task (content, priority, labels, project, date)
- `close` / `c` - Mark task(s) complete
- `delete` / `d` - Delete task(s)
- `labels` - List all labels
- `projects` - List all projects
- `add-project` / `ap` - Add new project
- `karma` - Show user karma
- `sync` / `s` - Sync cache
- `quick` / `q` - Quick add (uses REST endpoint)

Global flags: `--header`, `--color`, `--csv`, `--debug`, `--namespace`, `--indent`, `--project-namespace`

### API Communication

The client uses Todoist's **Sync API v1** (base URL: `https://api.todoist.com/api/v1/`):

- Commands are batched and sent to `POST /sync` endpoint
- Uses bearer token authentication (Authorization header)
- Operations (add, update, close, delete) are queued as commands
- `ExecCommands()` executes command batches atomically
- Quick add uses `POST /tasks/quick` (REST API)
- Completed tasks use `GET /tasks/completed/by_completion_date`

Command types: `item_add`, `item_update`, `item_close`, `item_delete`, `item_move`, `project_add`

### Integration Points

- **peco/fzf integration**: Shell functions in `todoist_functions.sh` provide interactive task selection
- **Shell completion**: Supports bash/zsh autocomplete via urfave/cli
- **Browser integration**: Can open URLs embedded in task content (markdown links)

## Dependencies

Key dependencies (from go.mod):
- `github.com/urfave/cli/v2` (v2.25.1) - CLI framework
- `github.com/fatih/color` (v1.13.0) - Colored output
- `github.com/spf13/viper` (v1.15.0) - Config file management
- `github.com/gofrs/uuid` (v3.2.0) - UUID generation for commands
- `github.com/pkg/browser` - Browser opening for URLs
- `github.com/rkoesters/xdg` - XDG Base Directory support
- `github.com/stretchr/testify` (v1.8.1) - Testing

## Important Notes

- The `filter_parser.go` file is **auto-generated** from `filter_parser.y` - never edit it directly
- When modifying the filter grammar, run `make prepare` to regenerate the parser
- The application automatically syncs on first run if token differs from cached token
- Item and project hierarchies use **in-memory linked structures** rather than recursive lookups
- Priority values in Todoist API are 1-4, where 1 is highest priority (p1); CLI inverts this (p1 → priority 4)
- Date handling supports multiple formats: RFC3339, natural language (via Todoist API), and custom filter syntax
- Filter syntax supports: `#Project`, `##Project` (include children), `@Label`, `p1`-`p4`, date expressions
- Shell integration files: `todoist_functions.sh` (peco), `todoist_functions_fzf.sh` (fzf)
