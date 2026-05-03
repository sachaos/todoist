---
description: Draft a user-facing spec for a feature from the meta feature-gaps issue, drive design decisions interactively, then post as a sub-issue.
model: claude-opus-4-7
---

You are helping implement features for the todoist CLI. Your job is to research the codebase, draft a user-facing spec, drive interactive design decisions, and post the approved spec as a new GitHub sub-issue of the active meta-issue.

## Step 1 — Identify the meta-issue

Check your memory for the active meta-issue number (look for a memory about "meta-issue" or "feature gaps tracker"). If not found, ask the user: "Which GitHub issue is the current feature-gaps tracker for this repo?" Save their answer to memory before continuing.

## Step 2 — Pick the feature

If the user provided an **issue number** as the argument (e.g. `#303` or `303`), read that issue with `gh issue view <number> --repo sachaos/todoist` and use it as the spec destination — skip creating a new issue in Step 6 and instead update the existing one.

If the user provided a feature name (text), use that as the topic.

Otherwise:
- Run `gh issue view <meta-issue> --repo sachaos/todoist` to get the current list
- Show the user the unchecked items and ask which one to spec

## Step 3 — Research

Spawn two agents in parallel (both Sonnet):

**Agent A (Explore):** Research the codebase for patterns relevant to this feature:
- How similar existing commands are structured (`close.go`, `delete.go`, `modify.go` as references)
- What `lib/` functions already exist that might be relevant
- What API communication patterns are used (`doRestApi` vs `ExecCommands`/sync)
- What flags/options similar commands use

**Agent B (API research):** Search the OpenAPI spec at `todoist-openapi.json` in the repo root (fetch a fresh copy if missing: `curl -s https://developer.todoist.com/openapi.json -o todoist-openapi.json`). Find:
- The relevant REST endpoint(s) and HTTP methods
- Required parameters and response codes (especially: does it return 200 or 204?)
- Any side effects documented in the API description
- Whether a sync API command equivalent exists

## Step 4 — Draft the spec

Write a user-facing spec covering:
- **Command shape**: name, aliases (or lack thereof), args, flags
- **Behavior**: what it does, success/failure cases, side effects
- **Help text**: what `--help` should show, including caveats
- **API choice**: REST vs sync command, and why
- **Cache behavior**: does it need to call `Sync()` after? Why or why not?
- **UX considerations**: watch for flag name collisions, upper/lowercase issues, conflicts with existing flags (per CLAUDE.md)
- **Out of scope**: what this PR won't include

## Step 5 — Resolve design decisions

Present the spec, then work through this standard punch list interactively. Skip questions with obvious answers from the research:

1. **Command name** — matches Todoist's terminology? Any collision with existing commands?
2. **Aliases** — short alias or none? Existing aliases in use: `l, a, m, c, d, s, q, ap, cl, c-l`
3. **Flag names** — any UX problems (similar names, upper/lowercase collisions, conflicts)?
4. **API approach** — REST or sync command? Batch or sequential?
5. **Error behavior** — stop on first error, or continue and report all?
6. **Success output** — silent (like `close`/`delete`) or print confirmation?
7. **Cache/sync** — call `Sync()` after, or not? Why?
8. **Side effects** — does the API do anything surprising (e.g. affects parent tasks)?
9. **Partial failure recovery** — what does the user do if only some IDs succeed?

Resolve each one before moving on. State your recommendation for each and ask the user to confirm or redirect.

## Step 6 — Post the issue

Once the user approves the spec:

**If the user pointed to an existing issue in Step 2**, update that issue's body with the final spec:
```
gh issue edit <number> --repo sachaos/todoist --body "<spec>"
```

**Otherwise**, create a new issue:
- Title: concise, imperative ("Add `X` command to...")
- Body: the full spec in markdown, opening with "Proposal for [item] tracked in #<meta-issue>"
- Attempt to link as a sub-issue: `gh api -X POST /repos/sachaos/todoist/issues/<meta-issue>/sub_issues -F sub_issue_id=<new-issue-id>`
- If the sub_issues API returns 404, fall back to posting a comment on the meta-issue: `gh issue comment <meta-issue> --repo sachaos/todoist --body "Spec posted: #<new-issue-id>"`

Report the issue URL to the user.
