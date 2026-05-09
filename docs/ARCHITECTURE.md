# Lib/CLI Architecture Guidelines

## Core Principle

`/lib` is a reusable Todoist client library designed to power multiple implementations: the CLI, a TUI interface, a web app, a GUI client, scripts, or anything else. Contributors should design lib code as if other teams will build clients on top of it.

## What Goes Where

| Concern | Lib | CLI |
|---------|-----|-----|
| **Data models** | ✅ | ❌ |
| **HTTP/API communication** | ✅ | ❌ |
| **API operations (Client methods)** | ✅ (AddItem, CloseItem, etc.) | ❌ |
| **Data serialization for API** | ✅ (AddParam, UpdateParam) | ❌ |
| **Data lookup & navigation** | ✅ (FindItem, ConstructItemTree) | ❌ |
| **Color/formatting output** | ❌ | ✅ |
| **Flag parsing & input translation** | ❌ | ✅ |
| **Filter parsing & evaluation** | ❌ | ✅ |
| **Sync() calls** | ❌ | ✅ (CLI handler only) |
| **Error messages to user** | Generic (what happened) | Specific (how to fix) |

## API Color Data

API color metadata (`Project.Color`, `Label.Color`) belongs in lib as part of the data model. Converting that data to terminal colors for output belongs in the CLI layer.
