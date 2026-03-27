# Smoke Ready Helper Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `make smoke-ready` and `make smoke-ready-json` so contributors can check local and GitHub smoke readiness without manually stitching together commands.

**Architecture:** Follow the existing repo pattern for `compare-local` and `worktree-*` helpers: add one focused Python script under `scripts/`, expose it through Make targets, and cover it with make-target tests in `cmd/cosbench-go/compare_local_make_test.go`. Keep the helper read-only and diagnostic-only.

**Tech Stack:** Make, Python 3, Go test suite, GitHub CLI

---

### Task 1: Add failing make-target tests

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Write a failing test for `smoke-ready`**
- [ ] **Step 2: Run the targeted test to confirm it fails because the target does not exist**
- [ ] **Step 3: Write a failing test for `smoke-ready-json`**
- [ ] **Step 4: Run the targeted test to confirm it fails because the target does not exist**

### Task 2: Implement the helper

**Files:**
- Create: `scripts/smoke_ready.py`
- Modify: `Makefile`

- [ ] **Step 1: Add a Python helper that reports local env presence, repo secret presence, workflow presence, readiness, and blockers**
- [ ] **Step 2: Expose `make smoke-ready` and `make smoke-ready-json` through the existing Python helper pattern**
- [ ] **Step 3: Run the targeted tests and make them pass**

### Task 3: Document the helper

**Files:**
- Modify: `README.md`
- Modify: `docs/legacy-live-run-checklist.md`

- [ ] **Step 1: Add concise usage notes for `smoke-ready` and `smoke-ready-json`**
- [ ] **Step 2: Re-run targeted tests plus a smoke invocation of the helper**
- [ ] **Step 3: Run `make validate`**
