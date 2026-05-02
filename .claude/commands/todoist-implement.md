---
description: Read a spec issue, implement the feature, open a PR, and run code review.
model: claude-opus-4-7
---

You are implementing a feature for the todoist CLI. The spec is in a GitHub issue number provided as the argument.

## Step 1 — Read the spec

Run: `gh issue view <issue-number> --repo sachaos/todoist`

Extract:
- The command name, args, and flags
- The API endpoint and HTTP method
- Behavior details (error handling, success output, sync behavior)
- Help text to use in `main.go` Description field

## Step 2 — Research existing patterns

Read these files to understand patterns to follow:
- `close.go` and `lib/item.go` — template for simple mutation commands
- `main.go` — how commands are registered (structure, flags, ArgsUsage, Description)
- `lib/todoist.go` — `doRestApi` and `doApi` helpers; note that 204 No Content is handled as success

## Step 3 — Create a branch

`git checkout -b <feature-branch-name>` using a descriptive kebab-case name.

## Step 4 — Spawn implementation agent (Sonnet)

Brief a general-purpose Sonnet agent with:
- The full spec (paste it verbatim)
- Specific files to create/modify and what each change should do
- Patterns to follow from the existing codebase
- Reminders: do NOT call `Sync(c)` unless the spec says to; do NOT commit or push; do NOT add features beyond the spec
- After implementing: run `make build` and `make test` and report results

Wait for the agent to complete before proceeding.

## Step 5 — Review and fix

Check the agent's output. If `make build` or `make test` failed, fix the issues directly. Verify the implementation matches the spec.

## Step 6 — Commit and push

```
git add <changed files>
git commit -m "<imperative summary>\n\nCo-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
git push -u origin <branch>
```

## Step 7 — Open the PR

```
gh pr create --title "..." --body "..."
```

PR body should include:
- "Closes #<issue-number>"
- Summary of changes (files modified and why)
- Design notes explaining non-obvious decisions (API choice, cache behavior)
- Test plan checklist (build, tests, manual steps)

## Step 8 — Run code review

Invoke the `/code-review:code-review` skill against the new PR number.

Report the PR URL to the user and wait for their manual testing before they call `/todoist-ship`.
