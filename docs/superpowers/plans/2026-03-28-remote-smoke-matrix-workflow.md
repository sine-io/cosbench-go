# Remote Smoke Matrix Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a scheduled and manually triggerable remote smoke matrix workflow that validates all currently supported `backend × scenario` combinations without gating default CI.

**Architecture:** Keep the existing `Remote Smoke Local` workflow as the targeted single-run entrypoint. Add a separate `Remote Smoke Matrix` workflow with a four-row strategy matrix, per-row artifact names, and the same helper invocation path. Lock the workflow contract with a small Python test and document the new manual entrypoint in README.

**Tech Stack:** GitHub Actions YAML, lightweight Python workflow-contract test, existing `make smoke-remote-local` helper, README documentation

---

### Task 1: Add A Failing Matrix Workflow Contract Test

**Files:**
- Create: `scripts/test_remote_smoke_matrix_workflow.py`

- [ ] **Step 1: Add a workflow-contract test for the matrix workflow**

Cover:
- workflow file exists
- it includes `schedule` and `workflow_dispatch`
- it declares the four expected matrix rows
- it sets `fail-fast: false`
- it passes both `SMOKE_REMOTE_LOCAL_BACKEND` and `SMOKE_REMOTE_LOCAL_SCENARIO`

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py -q
```

Expected:
- failure because the workflow file does not exist yet

### Task 2: Add The Matrix Workflow

**Files:**
- Create: `.github/workflows/remote-smoke-matrix.yml`

- [ ] **Step 1: Add the matrix workflow file**

Use:
- `name: Remote Smoke Matrix`
- triggers:
  - `schedule`
  - `workflow_dispatch`

- [ ] **Step 2: Add the four matrix rows**

Support:
- `s3 + single`
- `s3 + multistage`
- `sio + single`
- `sio + multistage`

- [ ] **Step 3: Add per-row artifact and summary handling**

Ensure:
- `if: always()` on artifact upload
- `if: always()` on summary export
- artifact names are unique per row

- [ ] **Step 4: Re-run the focused workflow-contract test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py -q
```

Expected:
- the matrix workflow contract test passes

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-matrix.yml scripts/test_remote_smoke_matrix_workflow.py
git commit -m "feat: add remote smoke matrix workflow"
```

### Task 3: Verify Local Surface And Document Manual Usage

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the new matrix workflow**

Add:
- one short explanation
- one `gh workflow run "Remote Smoke Matrix"` example

- [ ] **Step 2: Verify the existing remote smoke helper still passes for the full local matrix**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
SMOKE_REMOTE_LOCAL_BACKEND=sio timeout 240s make --no-print-directory smoke-remote-local
SMOKE_REMOTE_LOCAL_BACKEND=sio SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- all four combinations succeed locally

- [ ] **Step 3: Run the full Go test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all Go packages pass

- [ ] **Step 4: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 5: Commit the docs slice**

Run:
```bash
git add README.md docs/superpowers/specs/2026-03-28-remote-smoke-matrix-workflow-design.md docs/superpowers/plans/2026-03-28-remote-smoke-matrix-workflow.md
git commit -m "docs: record remote smoke matrix workflow"
```
