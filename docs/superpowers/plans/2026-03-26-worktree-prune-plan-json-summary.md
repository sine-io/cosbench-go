# Worktree Prune Plan JSON Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `worktree-prune-plan-json` return a structured `{summary, rows}` payload and update repo-local consumers, tests, and docs to match.

**Architecture:** Keep the existing prune-plan row generation logic intact and wrap it in a small summary layer for JSON mode only. Update the cleanup-report aggregator and command-level tests to consume the structured payload instead of a bare array.

**Tech Stack:** Python 3 helper scripts, Go command tests, Make targets, Markdown docs

---

### Task 1: Lock the new JSON contract with failing tests

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`
- Modify: `cmd/cosbench-go/worktree_audit_script_test.go`

- [ ] **Step 1: Write the failing test expectations**

Update `TestWorktreePrunePlanJSONTargetRuns` to expect:
- a top-level `summary` object
- a top-level `rows` array
- `summary.total` to be non-negative

Update the direct script test to validate:
- `summary.base_ref`
- `summary.total`
- `summary.integrated`
- `rows[0].state`

- [ ] **Step 2: Run targeted tests to verify they fail**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because `scripts/worktree_prune_plan.py --json` still returns a bare array.

### Task 2: Implement the structured prune-plan payload

**Files:**
- Modify: `scripts/worktree_prune_plan.py`

- [ ] **Step 1: Keep row generation unchanged**

Do not change the filtering rules or text output path.

- [ ] **Step 2: Add JSON summary generation**

Return:
- `summary.base_ref`
- `summary.current_worktree`
- `summary.total`
- `summary.merged`
- `summary.integrated`
- `rows`

- [ ] **Step 3: Run targeted tests to verify they pass**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 3: Update structured consumers and docs

**Files:**
- Modify: `scripts/worktree_cleanup_report.py`
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Preserve the new prune-plan object in cleanup-report JSON**

Ensure `worktree_cleanup_report.py --json` carries through the structured prune-plan payload without flattening it.

- [ ] **Step 2: Update docs**

Document that `make --no-print-directory worktree-prune-plan-json` now returns summary metadata plus rows.

- [ ] **Step 3: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS

### Task 4: Commit

**Files:**
- Commit the modified tests, scripts, and docs together

- [ ] **Step 1: Create a single commit**

```bash
git add cmd/cosbench-go/compare_local_make_test.go \
  cmd/cosbench-go/worktree_audit_script_test.go \
  scripts/worktree_prune_plan.py \
  scripts/worktree_cleanup_report.py \
  README.md AGENTS.md BOARD.md TODO.md \
  docs/superpowers/specs/2026-03-26-worktree-prune-plan-json-summary-design.md \
  docs/superpowers/plans/2026-03-26-worktree-prune-plan-json-summary.md
git commit -m "feat: add prune plan JSON summary"
```
