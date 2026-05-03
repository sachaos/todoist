---
description: Read a spec issue, implement the feature, open a PR, and run code review.
model: claude-sonnet-4-6
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
- **Tests**: always add a test for the zero-input/no-args guard (e.g. command with no IDs, or a flag that requires another flag) using `newTestContext()` from `show_test.go`. Add any other pure-logic tests that don't require HTTP. Do NOT introduce HTTP mocking infrastructure.
- After implementing: run `make build` and `make test` and report results

Wait for the agent to complete before proceeding.

## Step 5 — Review and fix

Check the agent's output. If `make build` or `make test` failed, fix the issues directly. Verify the implementation matches the spec.

## Step 6 — Commit and push

```
git add <changed files>
git commit -m "<imperative summary>\n\nCo-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
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

## Step 9 — Prompt for manual testing

Tell the user the PR URL and present a manual test checklist derived from the spec. At minimum include:

- The happy path (valid ID, expected success behavior)
- The no-args case (should show usage, not a cryptic error)
- An invalid ID (should surface a clear API error message)
- Any edge cases called out in the spec (recurring tasks, already-active tasks, etc.)

Then explicitly ask: **"Please run through the test checklist above. When you're done, call `/todoist-ship <PR-number>` to merge."**

Do not proceed further — the skill ends here and waits for the user.
