---
description: Merge a PR, check off the item on the meta feature-gaps issue, and update memory.
model: claude-haiku-4-5-20251001
---

You are closing out a completed feature PR for the todoist CLI. The PR number is provided as the argument.

## Step 1 — Identify the meta-issue

Check your memory for the active meta-issue number. If not found, ask the user before continuing.

## Step 2 — Verify the PR is ready

Run: `gh pr view <PR-number> --repo sachaos/todoist`

- If it is a draft or closed without merging, report that and stop.
- If it is already merged, skip Step 3 and continue with Step 4.
- Otherwise confirm it is open and proceed.

## Step 3 — Merge

```
gh pr merge <PR-number> --repo sachaos/todoist --squash --delete-branch
```

## Step 4 — Check off item on the meta-issue

Run: `gh issue view <meta-issue> --repo sachaos/todoist --json body --jq '.body'`

Find the line corresponding to the merged feature (match by keyword from the PR title). Change `- [ ]` to `- [x]` for that line only.

```
gh issue edit <meta-issue> --repo sachaos/todoist --body "<updated body>"
```

## Step 5 — Update memory

Update your memory for this project: move the completed item from open to completed, recording the spec issue number, PR number, and today's date.

## Step 6 — Report

Tell the user:
- The PR was merged
- Which meta-issue item was checked off
- What open items remain
