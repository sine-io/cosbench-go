# Worktree Cleanup Report Metadata Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add generation-time and current-worktree metadata to the Markdown cleanup report summary.

**Architecture:** Keep the change local to `worktree_cleanup_report.py`. Reuse the existing structured prune-plan summary for the current-worktree path and compute a single UTC timestamp per report render.

**Tech Stack:** Python 3 helper scripts, Go make-target tests, Markdown docs

---

### Task 1: Lock the summary metadata contract

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Require the new summary lines**

Expect both stdout and the written Markdown report file to contain:
- `- Generated at:`
- `- Current worktree:`

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because the Markdown report does not yet render those lines.

### Task 2: Implement the metadata

**Files:**
- Modify: `scripts/worktree_cleanup_report.py`

- [ ] **Step 1: Add a timestamp helper**

Generate a single UTC timestamp with second precision for each report render.

- [ ] **Step 2: Add the two summary lines**

Use:
- the generated timestamp
- `prune_plan["summary"]["current_worktree"]`

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 3: Sync docs and run full verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the new metadata**

Note that the Markdown cleanup report now records generation time and current-worktree context.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
