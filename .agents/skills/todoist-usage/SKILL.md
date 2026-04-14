---
name: todoist-usage
description: How to interact with the Todoist CLI to manage a user's tasks, projects, and labels. Make sure to use this skill whenever the user mentions Todoist, asks about their tasks, wants to add a reminder, or needs to check their to-do list, even if they don't explicitly ask you to 'use Todoist'.
---

# Todoist CLI Usage Workflow

You are equipped to manage the user's Todoist tasks using the `todoist` CLI. This skill gives you the procedural knowledge required to use this tool effectively.

## 1. Syncing State (Crucial First Step)

The CLI relies on a local cache for speed. This means the cache can become stale if tasks were updated on another device.
Before relying on task IDs, reading the task list, or if the output appears to be out of sync with what the user expects, **always consider running a sync**.

```bash
todoist sync
```

## 2. Listing and Filtering Tasks

The primary way to view tasks is via `todoist list` (or `todoist l`). 

### Using Filters

Filters are extremely powerful. Rather than pulling all tasks and parsing them yourself, push the filtering down to the CLI. The syntax is based on Todoist's official filter queries.

**Why filter?** Dumping all tasks is slow and context-heavy. Always try to filter down to what the user cares about.

**Examples of filters:**
- Due scenarios: `todoist list --filter '(overdue | today)'`
- Project specific: `todoist list --filter '#ProjectName'`
- Label specific: `todoist list --filter '@LabelName'`
- Compound filters: `todoist list --filter '(overdue | today) & p1'`

*Note for Completed Tasks:* To see tasks that have already been checked off, use the dedicated command `todoist completed-list` (or `cl`). This differs from standard lists and defaults to a 90-day range (premium users only).

## 3. Adding and Modifying Tasks

When adding a task, you can pass arguments directly. Quick add allows for natural language text processing.

**Example 1: Standard Add**
```bash
todoist add "Buy milk" -p p1 -d "tomorrow at 10:00"
```

**Example 2: Quick Add (Leveraging Natural Language)**
```bash
todoist quick "Buy milk #Groceries @Urgent tomorrow"
```

To modify an existing task, use `todoist modify`.
To show deep details about a specific task, use `todoist show <task-id>`.

## 4. Closing and Deleting Tasks

When the user asks to complete a task, use the `close` subcommand. Make sure you use the correct alphanumeric task ID.
```bash
todoist close <alphanumeric-id>
```

If it needs to be permanently deleted instead of just completed:
```bash
todoist delete <alphanumeric-id>
```

## 5. Navigating Projects and Labels

If you need to discover what projects or labels exist (for example, to construct a good filter or to categorize a new task), use:
- `todoist projects`
- `todoist labels`

To create a new project:
- `todoist add-project`

## Important Edge Cases & Considerations

- **Priority Inversion**: In the Todoist API, priorities are 1-4 with 1 being the highest (e.g. `p1`). The CLI interprets `p1` as highest, but internally it might display as priority 4. Default to using the `p1`, `p2`, `p3`, `p4` flags as the user expects.
- **Alphanumeric IDs**: Because the tool was migrated to the V1 Sync API, task and project IDs are alphanumeric hashes, not integers. 
- **Formatting Output**: You can append `--color` to commands if the output is going directly to a terminal, but if you are scraping the text, leave it off to avoid ANSI escape sequences cluttering the output.
