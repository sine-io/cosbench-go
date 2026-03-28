# Smoke Local MinIO Migration Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the moto-backed `smoke-local` helper with a MinIO-backed helper so local smoke evidence exercises a real S3-compatible server without extra Python package installation.

**Architecture:** Keep the existing `make smoke-local` interface and mock-mode tests, but change the helper internals to manage a local MinIO binary under `.artifacts/minio/`. Remove Python package pinning and any workflow setup that only existed for moto.

**Tech Stack:** Make, Python 3 stdlib, Go smoke tests, MinIO binary, GitHub Actions

---

### Task 1: Lock the visible contract change

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Add a failing assertion that `smoke-local` identifies MinIO in its summary**
- [ ] **Step 2: Run the targeted smoke-local tests to confirm they fail**

### Task 2: Migrate the helper to MinIO

**Files:**
- Modify: `scripts/smoke_local.py`
- Modify: `cmd/cosbench-go/script_dependency_test.go`
- Delete: `requirements-smoke-local.txt`

- [ ] **Step 1: Replace moto bootstrap/install logic with MinIO download/bootstrap logic**
- [ ] **Step 2: Keep `SMOKE_LOCAL_MOCK` behavior intact for fast tests**
- [ ] **Step 3: Run the targeted smoke-local tests and make them pass**
- [ ] **Step 4: Run `make --no-print-directory smoke-local` and confirm the real helper still passes**

### Task 3: Simplify workflow and docs

**Files:**
- Modify: `.github/workflows/smoke-local.yml`
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`
- Modify: `docs/legacy-live-run-checklist.md`

- [ ] **Step 1: Remove workflow steps that only existed for moto/pip dependencies**
- [ ] **Step 2: Update docs to describe `smoke-local` as MinIO-backed**
- [ ] **Step 3: Run `make validate`**
