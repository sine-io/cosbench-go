# Worktree JSON Meta Envelope Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a consistent top-level `meta` envelope to the machine-readable worktree helpers.

**Architecture:** Keep the change additive. Reuse existing timestamp/base-ref/current-worktree values where available, expose them under `meta`, and preserve current top-level fields for compatibility.

**Tech Stack:** Python 3 helper scripts, Go make-target tests, Markdown docs

---

### Task 1: Lock the `meta` envelope contract

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Require `meta`**

Expect:
- `meta.generated_at`
- `meta.base_ref`
- `meta.current_worktree`

for audit, prune-plan, and cleanup-report JSON outputs.

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because those helpers do not yet emit `meta`.

### Task 2: Implement the additive envelopes

**Files:**
- Modify: `scripts/worktree_audit.py`
- Modify: `scripts/worktree_prune_plan.py`
- Modify: `scripts/worktree_cleanup_report.py`

- [ ] **Step 1: Build the current metadata once**

Reuse the existing generated timestamps and current-worktree values.

- [ ] **Step 2: Add `meta`**

Keep the current top-level fields unchanged.

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 3: Sync docs and run full verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document `meta` as the preferred metadata entrypoint**

Make it clear that the older top-level fields remain for compatibility.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
