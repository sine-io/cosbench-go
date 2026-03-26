# Worktree Prune Plan Metadata Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add generation-time, base-ref, and current-worktree metadata to the plain-text prune-plan output.

**Architecture:** Keep the change local to `worktree_prune_plan.py`. Reuse the existing timestamp helper and already-available base-ref/current-worktree values, and emit them as comment lines above the command list.

**Tech Stack:** Python 3 helper scripts, Go make-target tests, Markdown docs

---

### Task 1: Lock the text metadata contract

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Require the new header lines**

Expect text prune-plan output to include:
- `# Generated at:`
- `# Base ref:`
- `# Current worktree:`

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because the text prune-plan output does not yet render those header lines.

### Task 2: Implement the metadata

**Files:**
- Modify: `scripts/worktree_prune_plan.py`

- [ ] **Step 1: Reuse the existing timestamp helper**

Generate one UTC timestamp per run.

- [ ] **Step 2: Add the three header lines**

Keep them prefixed with `#` so the output remains copy/paste friendly.

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 3: Sync docs and run full verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the new text metadata**

Note that the text prune-plan output now records generation time, base ref, and current worktree.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
