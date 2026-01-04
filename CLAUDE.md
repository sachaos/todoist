# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a CLI client for Todoist written in Go. It provides command-line access to Todoist's task management features with support for filtering, projects, labels, and offline caching.

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
   - `main.go` - CLI setup using urfave/cli/v2 framework
   - `list.go`, `add.go`, `close.go`, etc. - Command implementations
   - `format.go` - Output formatting (TSV/CSV writers)
   - `cache.go` - Local cache management for offline access

2. **Library Layer (`lib/` directory)**: Core Todoist API client and data models
   - `lib/todoist.go` - HTTP client and API communication
   - `lib/sync.go` - Store structure and tree construction
   - `lib/item.go`, `lib/project.go`, `lib/label.go` - Data models
   - `lib/command.go` - Command serialization for Todoist Sync API
   - `lib/interface.go` - Common interfaces and helper functions

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

### API Communication

The client uses Todoist's **Sync API** (not the REST API):

- Commands are batched and sent to `/sync` endpoint
- Uses bearer token authentication
- Operations (add, update, close, delete) are queued as commands
- `ExecCommands()` executes command batches atomically

### Integration Points

- **peco/fzf integration**: Shell functions in `todoist_functions.sh` provide interactive task selection
- **Shell completion**: Supports bash/zsh autocomplete via urfave/cli
- **Browser integration**: Can open URLs embedded in task content (markdown links)

## Important Notes

- The `filter_parser.go` file is **auto-generated** from `filter_parser.y` - never edit it directly
- When modifying the filter grammar, run `make prepare` to regenerate the parser
- The application automatically syncs on first run if token differs from cached token
- Item and project hierarchies use **in-memory linked structures** rather than recursive lookups
- Priority values in Todoist API are 1-4, where 1 is highest priority (p1)
- Date handling supports multiple formats: RFC3339, natural language (via Todoist API), and custom filter syntax
