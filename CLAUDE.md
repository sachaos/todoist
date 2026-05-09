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

**For detailed lib/CLI boundary guidelines and contributor guidance, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).**

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

## Architectural Guidelines

These rules reflect hard-won patterns. Follow them when adding or modifying commands.

### REST vs Sync API

Use the **Sync API** (`ExecCommands`) for operations on **active tasks** that live in the local cache — add, update, close, delete, move. These can be batched atomically.

Use the **REST API** (`doRestApi`) when:
- The operation targets resources **not in the local cache** (e.g. completed tasks, which are fetched live and never cached)
- No sync command type exists for the operation
- The endpoint is one-shot by nature (e.g. `POST /tasks/{id}/reopen`)

### When to call Sync(c) after a mutation

**Call `Sync(c)`** when the command modifies an active task that is currently in the local cache. This keeps `todoist list` accurate immediately after the operation. Examples: `close`, `delete`, `add`, `modify`.

**Do NOT call `Sync(c)`** when the command operates on completed tasks or resources that were never in the cache. The cache stays internally consistent because it never held those items. Examples: `reopen`. Add a comment in the handler explaining the omission so future maintainers don't add it back by mistake.

### doRestApi behavior

`doRestApi` in `lib/todoist.go` treats both **200 OK** and **204 No Content** as success. When adding a new REST-based lib function, always check the OpenAPI spec (at `todoist-openapi.json` in the repo root) to confirm which status code the endpoint returns — don't assume 200.

### Lib layer conventions

Functions in `lib/item.go` that operate on multiple IDs take `[]string`, not a single `string`. This keeps the lib API consistent and lets callers decide whether to loop or not. Example signatures:
```go
func (c *Client) CloseItem(ctx context.Context, ids []string) error
func (c *Client) DeleteItem(ctx context.Context, ids []string) error
```

### Error handling conventions

When a CLI handler loops over multiple IDs, **stop on first error** and wrap the error with the failing ID for clarity:
```go
if err := client.SomeOperation(ctx, id); err != nil {
    return fmt.Errorf("failed to <action> task %s: %w", id, err)
}
```

This lets users re-run the command with the remaining IDs after fixing the cause.

### Command registration conventions

In `main.go`, every command must have:
- `ArgsUsage` — describes expected arguments (e.g. `"<Item ID> [<Item ID>...]"`)
- `Description` — multi-line string shown in `--help`; explain non-obvious behavior (cache effects, side effects, recovery from partial failure)

When a command is invoked with no arguments and arguments are required, return a descriptive error with usage guidance — not the generic `CommandFailed`:
```go
if c.Args().Len() == 0 {
    return fmt.Errorf("no task IDs provided\nUsage: todoist <cmd> <Item ID> [<Item ID>...]\n...")
}
```

### Existing command aliases

Check `main.go` for the current list of registered aliases before adding a new one. New commands may omit an alias rather than introduce a confusing one.

### Cache architecture

The local cache (`~/.cache/todoist/cache.json`) stores only **active** tasks. Completed tasks are never cached — they are always fetched live from `GET /tasks/completed/by_completion_date`. Reopening a completed task does not require a cache update; the task will appear in `todoist list` after the next `todoist sync`.

### Testing approach

The test infrastructure supports two patterns:

1. **Pure logic and read-only handlers**: Use `newTestContext()` from `show_test.go` with an in-memory `Store`. No HTTP calls needed. Always test the zero-args guard this way.
2. **Mutation handlers**: These require HTTP mocking (no `httptest.Server` infrastructure currently exists). Follow the existing convention — mutation handlers (`close`, `delete`, `add`) have no tests. Do not introduce HTTP mocking infrastructure without explicit discussion.

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
- **OpenAPI Specification**: The full Todoist API v1 OpenAPI spec is available at `https://developer.todoist.com/openapi.json` and saved locally as `todoist-openapi.json` in the repository root. This can be used for code generation, documentation, and API tool integrations.
- When adding new command-line flags, watch out for UX problems: very similar flag names, upper/lowercase collisions, and conflicts with existing flags that could confuse users

## AI Disclosure Guidelines

When working with GitHub:

- **Creating or significantly modifying issues/PRs**: Add the `AI-generated` label to the issue or PR
- **Posting comments on issues/PRs**: End the comment with `_(generated by Claude)_` to disclose AI involvement, even if posting on behalf of the user
