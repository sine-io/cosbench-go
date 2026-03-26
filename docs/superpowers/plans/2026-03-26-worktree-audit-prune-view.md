# Worktree Audit Prune View Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a prune-candidates audit view and wire it into cleanup report outputs.

**Architecture:** Reuse the same row classification used by the full audit and express prune eligibility as one additional filter. Expose that filter through Make targets and pass the resulting audit-form rows into cleanup report outputs, keeping `prune_plan` reserved for executable cleanup commands.

**Tech Stack:** Python 3 helper scripts, Make targets, Go make-target tests, Markdown docs

---

### Task 1: Lock the prune-view contract with failing tests

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Add prune-only target expectations**

Require:
- `worktree-audit-prune`
- `worktree-audit-prune-json`
- `## Prune Candidates` in cleanup-report markdown
- `prune_candidates` in cleanup-report JSON

- [ ] **Step 2: Run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: FAIL because those targets and report sections do not exist yet.

### Task 2: Implement the prune-only audit view

**Files:**
- Modify: `scripts/worktree_audit.py`
- Modify: `Makefile`

- [ ] **Step 1: Add `--prune-only` filtering**

Match the same eligibility rules already used by `worktree_prune_plan.py`.

- [ ] **Step 2: Add Make targets**

Expose both text and JSON prune-only views.

### Task 3: Surface prune candidates in cleanup report

**Files:**
- Modify: `scripts/worktree_cleanup_report.py`

- [ ] **Step 1: Add a `prune_candidates` JSON payload**

Keep `prune_plan` as the command-oriented object.

- [ ] **Step 2: Add an `## Prune Candidates` Markdown section**

Place it before `## Prune Plan`.

- [ ] **Step 3: Re-run the targeted test**

Run: `GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go`

Expected: PASS

### Task 4: Sync docs and run full verification

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the new prune-candidates view**

Add brief usage notes for the text and JSON prune-only targets and the cleanup-report section.

- [ ] **Step 2: Run full verification**

Run: `GO=$(which go || echo /snap/bin/go) go test ./...`

Expected: PASS

Run: `GO=$(which go || echo /snap/bin/go) go build ./...`

Expected: PASS
