# Remote Smoke Recovery Matrix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a non-blocking remote recovery matrix workflow that validates `s3 + recovery` and `sio + recovery` on GitHub-hosted runners.

**Architecture:** Keep the current helper and the dedicated `Remote Smoke Recovery` workflow intact. Add a separate `Remote Smoke Recovery Matrix` workflow with two matrix rows, per-row artifacts, and one aggregate job backed by a dedicated Python aggregation script for recovery summaries.

**Tech Stack:** GitHub Actions YAML, Python workflow contract tests, Python aggregation script

---

### Task 1: Add Failing Tests For The Recovery Matrix Contract

**Files:**
- Create: `scripts/test_remote_smoke_recovery_matrix_workflow.py`
- Create: `scripts/test_aggregate_remote_smoke_recovery_matrix.py`

- [ ] **Step 1: Add a workflow-contract test**

Cover:
- workflow file exists
- `workflow_dispatch` and `schedule` are present
- matrix rows include:
  - `backend=s3`, `scenario=recovery`
  - `backend=sio`, `scenario=recovery`
- aggregate job exists

- [ ] **Step 2: Add an aggregation-script test**

Cover:
- two recovery row summaries aggregate into Markdown and JSON output
- missing rows are represented as `missing`

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_matrix_workflow.py scripts/test_aggregate_remote_smoke_recovery_matrix.py -q
```

Expected:
- failure because the workflow and aggregation script do not exist yet

### Task 2: Add The Recovery Matrix Workflow

**Files:**
- Create: `.github/workflows/remote-smoke-recovery-matrix.yml`
- Create: `scripts/aggregate_remote_smoke_recovery_matrix.py`

- [ ] **Step 1: Add the dedicated aggregation script**

Implement:
- expected rows:
  - `s3 + recovery`
  - `sio + recovery`
- aggregate JSON + Markdown output

- [ ] **Step 2: Add the matrix workflow**

Use:
- `workflow_dispatch`
- `schedule`
- matrix rows for `s3` and `sio`
- `fail-fast: false`

- [ ] **Step 3: Add aggregate job**

Ensure:
- downloads row artifacts
- runs the aggregation script
- writes summary
- uploads aggregate artifact

- [ ] **Step 4: Re-run the focused tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_matrix_workflow.py scripts/test_aggregate_remote_smoke_recovery_matrix.py -q
```

Expected:
- both tests pass

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-recovery-matrix.yml scripts/test_remote_smoke_recovery_matrix_workflow.py scripts/aggregate_remote_smoke_recovery_matrix.py scripts/test_aggregate_remote_smoke_recovery_matrix.py
git commit -m "feat: add remote smoke recovery matrix"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the new recovery matrix workflow**

Add:
- one short note
- one `gh workflow run "Remote Smoke Recovery Matrix"` example

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
git add README.md docs/superpowers/specs/2026-03-29-remote-smoke-recovery-matrix-design.md docs/superpowers/plans/2026-03-29-remote-smoke-recovery-matrix.md
git commit -m "docs: record remote smoke recovery matrix"
```
