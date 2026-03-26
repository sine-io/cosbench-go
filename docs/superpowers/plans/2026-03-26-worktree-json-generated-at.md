# Worktree JSON Generated-At Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a top-level `generated_at` timestamp to the machine-readable worktree JSON helpers.

**Architecture:** Reuse a small UTC timestamp helper in each Python script and keep the new metadata strictly top-level so existing `summary` and `rows` payloads do not need to change shape.

**Tech Stack:** Python 3 helper scripts, Go make-target tests, Markdown docs

---

### Task 1: Lock the timestamp contract with failing tests

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Add `generated_at` expectations**

Require a non-empty top-level `generated_at` field in:
- `worktree-audit-json`
- `worktree-prune-plan-json`
- `worktree-cleanup-report-json`

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because those JSON helpers do not yet emit `generated_at`.

### Task 2: Implement `generated_at`

**Files:**
- Modify: `scripts/worktree_audit.py`
- Modify: `scripts/worktree_prune_plan.py`
- Modify: `scripts/worktree_cleanup_report.py`

- [ ] **Step 1: Add a UTC timestamp helper**

Use RFC 3339 / ISO 8601 with `Z` suffix and second precision.

- [ ] **Step 2: Add the field to each JSON payload**

Do not alter existing nested summary structures.

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

Briefly note that the machine-readable worktree helpers now carry generation timestamps.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
