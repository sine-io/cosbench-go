# Worktree Cleanup Report Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `Integrated`, `Stale`, and `Prune candidates` counts to the Markdown summary emitted by `make worktree-cleanup-report`.

**Architecture:** Reuse the existing structured audit summary and structured prune-plan JSON rather than computing additional counts from the rendered text sections. Keep the text sections unchanged and limit the behavior change to the summary block.

**Tech Stack:** Python 3 helper scripts, Go make-target tests, Markdown docs

---

### Task 1: Lock the new Markdown summary contract

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Write the failing expectation**

Require the report output and written report file to include:
- `- Integrated:`
- `- Stale:`
- `- Prune candidates:`

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because the Markdown summary does not yet render those lines.

### Task 2: Implement the summary expansion

**Files:**
- Modify: `scripts/worktree_cleanup_report.py`

- [ ] **Step 1: Reuse structured inputs**

Load the structured prune-plan JSON once and keep the current text sections unchanged.

- [ ] **Step 2: Add the three summary lines**

Render:
- `Integrated`
- `Stale`
- `Prune candidates`

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 3: Sync docs and verify

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Update docs**

Document that the Markdown cleanup report summary now includes integrated, stale, and prune-candidate counts.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
