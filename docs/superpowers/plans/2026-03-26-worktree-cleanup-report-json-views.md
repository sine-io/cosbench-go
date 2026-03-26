# Worktree Cleanup Report JSON Views Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a top-level `views` container to `worktree-cleanup-report-json` while preserving the current top-level section keys.

**Architecture:** Keep the change additive. Build each section object once, publish it under `views`, and continue exposing the existing top-level aliases so current consumers keep working.

**Tech Stack:** Python 3 helper scripts, Go make-target tests, Markdown docs

---

### Task 1: Lock the `views` contract

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Require `views`**

Expect:
- top-level `views`
- `views.merged`
- `views.integrated`
- `views.stale`
- `views.prune_candidates`
- `views.prune_plan`

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because `worktree-cleanup-report-json` does not yet expose the `views` object.

### Task 2: Implement the additive JSON shape

**Files:**
- Modify: `scripts/worktree_cleanup_report.py`

- [ ] **Step 1: Build section payloads once**

Avoid duplicate subprocess calls where possible.

- [ ] **Step 2: Add `views`**

Keep the existing top-level aliases in place.

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 3: Sync docs and run full verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document `views` as the preferred JSON entrypoint**

Make it clear that the top-level aliases remain for compatibility.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
