# Worktree JSON View Envelope Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add additive `views` envelopes to `worktree-audit-json` and `worktree-prune-plan-json`.

**Architecture:** Keep the change strictly additive. Build the existing summary/rows payload once, then expose it under a named `views` key while preserving the current top-level aliases.

**Tech Stack:** Python 3 helper scripts, Go make-target tests, Markdown docs

---

### Task 1: Lock the `views` envelope contract

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Require `views.audit` and `views.prune_plan`**

Add expectations for:
- `views.audit.summary`
- `views.audit.rows`
- `views.prune_plan.summary`
- `views.prune_plan.rows`

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because the single-view helpers do not yet emit `views`.

### Task 2: Implement the additive envelopes

**Files:**
- Modify: `scripts/worktree_audit.py`
- Modify: `scripts/worktree_prune_plan.py`

- [ ] **Step 1: Build the existing payload once**

Do not duplicate summary/rows computation.

- [ ] **Step 2: Add the `views` wrapper**

Keep the top-level aliases for compatibility.

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 3: Sync docs and run full verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document `views` as the preferred entrypoint**

Clarify that `summary` and `rows` remain available for compatibility.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
