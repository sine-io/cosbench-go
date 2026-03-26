# Worktree Audit Integrated View Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add first-class integrated-only audit views and surface them in the cleanup report.

**Architecture:** Extend `worktree_audit.py` with one additional filter flag, expose it through Make targets, and reuse the new view inside the cleanup report instead of reconstructing integrated rows indirectly from the full audit or prune plan.

**Tech Stack:** Python 3 helper scripts, Make targets, Go make-target tests, Markdown docs

---

### Task 1: Lock the integrated-view contract with failing tests

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Add integrated-only target expectations**

Require:
- `worktree-audit-integrated` to print only integrated rows
- `worktree-audit-integrated-json` to emit only integrated rows
- `worktree-cleanup-report` to include `## Integrated`
- `worktree-cleanup-report-json` to include an `integrated` key

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because the integrated-only targets and cleanup-report section do not exist yet.

### Task 2: Implement the integrated-only view

**Files:**
- Modify: `scripts/worktree_audit.py`
- Modify: `Makefile`

- [ ] **Step 1: Add `--integrated-only` filtering**

Mirror the structure already used by `--merged-only` and `--stale-only`.

- [ ] **Step 2: Add Make targets**

Expose both text and JSON variants.

### Task 3: Wire the view into cleanup report output

**Files:**
- Modify: `scripts/worktree_cleanup_report.py`

- [ ] **Step 1: Add an integrated JSON payload**

Expose the integrated-only audit JSON in `worktree-cleanup-report-json`.

- [ ] **Step 2: Add an `## Integrated` Markdown section**

Keep the existing section ordering readable:
- Summary
- Merged
- Integrated
- Stale
- Prune Plan

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 4: Sync docs and run full verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the new integrated views**

Add brief usage notes for the two new Make targets and the integrated cleanup-report section.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
