# Smoke S3 Matrix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a manual matrix workflow that runs `make smoke-s3` for both `s3` and `sio` and publishes an aggregate summary.

**Architecture:** Keep the existing `Smoke S3` workflow intact as the single-backend entrypoint. Add a separate `Smoke S3 Matrix` workflow with two matrix rows, per-row artifacts, and one aggregate job backed by a small Python aggregation script for the output text files.

**Tech Stack:** GitHub Actions YAML, Python workflow contract tests, Python aggregation script

---

### Task 1: Add Failing Tests For The Smoke S3 Matrix Contract

**Files:**
- Create: `scripts/test_smoke_s3_matrix_workflow.py`
- Create: `scripts/test_aggregate_smoke_s3_matrix.py`

- [ ] **Step 1: Add a workflow-contract test**

Cover:
- workflow file exists
- `workflow_dispatch` exists
- matrix rows include `backend=s3` and `backend=sio`
- aggregate job exists
- run step uses `GO=go make smoke-s3`

- [ ] **Step 2: Add an aggregation-script test**

Cover:
- per-row text outputs aggregate into one Markdown and JSON output
- missing rows are represented as `missing`

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_s3_matrix_workflow.py scripts/test_aggregate_smoke_s3_matrix.py -q
```

Expected:
- failure because the workflow and aggregation script do not exist yet

### Task 2: Add The Smoke S3 Matrix Workflow

**Files:**
- Create: `.github/workflows/smoke-s3-matrix.yml`
- Create: `scripts/aggregate_smoke_s3_matrix.py`

- [ ] **Step 1: Add the aggregation script**

Implement:
- read row output files
- emit aggregate JSON + Markdown

- [ ] **Step 2: Add the manual matrix workflow**

Use:
- `workflow_dispatch`
- matrix rows for `backend=s3` and `backend=sio`
- existing `COSBENCH_SMOKE_*` secrets and optional inputs

- [ ] **Step 3: Add aggregate job**

Ensure:
- downloads row artifacts
- runs the aggregation script
- writes summary
- uploads aggregate artifact

- [ ] **Step 4: Re-run the focused tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_s3_matrix_workflow.py scripts/test_aggregate_smoke_s3_matrix.py -q
```

Expected:
- both tests pass

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add .github/workflows/smoke-s3-matrix.yml scripts/aggregate_smoke_s3_matrix.py scripts/test_smoke_s3_matrix_workflow.py scripts/test_aggregate_smoke_s3_matrix.py
git commit -m "feat: add smoke s3 matrix workflow"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the new workflow**

Add:
- one short note
- one `gh workflow run "Smoke S3 Matrix"` example

- [ ] **Step 2: Run the full Go test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all Go packages pass

- [ ] **Step 3: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 4: Commit the docs slice**

Run:
```bash
git add README.md docs/superpowers/specs/2026-03-29-smoke-s3-matrix-design.md docs/superpowers/plans/2026-03-29-smoke-s3-matrix.md
git commit -m "docs: record smoke s3 matrix"
```
