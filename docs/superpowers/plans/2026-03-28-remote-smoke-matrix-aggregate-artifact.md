# Remote Smoke Matrix Aggregate Artifact Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a downloadable aggregate artifact to the `Remote Smoke Matrix` workflow without changing helper behavior or default CI.

**Architecture:** Keep the existing matrix rows and aggregate summary generation intact. Add one upload step in the aggregate job, lock it with the workflow-contract test, and document the new downloadable artifact in README.

**Tech Stack:** GitHub Actions YAML, existing Python workflow-contract test, README documentation

---

### Task 1: Add A Failing Workflow Contract Assertion

**Files:**
- Modify: `scripts/test_remote_smoke_matrix_workflow.py`

- [ ] **Step 1: Extend the workflow-contract test**

Cover:
- aggregate job contains an upload step
- artifact name is `remote-smoke-matrix-aggregate`
- upload path is `.artifacts/remote-smoke-matrix-aggregate`

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py -q
```

Expected:
- failure because the aggregate upload step does not exist yet

### Task 2: Add Aggregate Artifact Upload

**Files:**
- Modify: `.github/workflows/remote-smoke-matrix.yml`

- [ ] **Step 1: Add the aggregate artifact upload step**

Ensure:
- stable artifact name
- upload path points at `.artifacts/remote-smoke-matrix-aggregate`
- current summary export behavior stays intact

- [ ] **Step 2: Re-run the focused workflow-contract test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py -q
```

Expected:
- workflow-contract test passes

- [ ] **Step 3: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-matrix.yml scripts/test_remote_smoke_matrix_workflow.py
git commit -m "feat: add remote smoke matrix aggregate artifact"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the aggregate artifact**

Add one short note that the matrix workflow uploads a downloadable aggregate artifact in addition to the combined summary.

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
git add README.md docs/superpowers/specs/2026-03-28-remote-smoke-matrix-aggregate-artifact-design.md docs/superpowers/plans/2026-03-28-remote-smoke-matrix-aggregate-artifact.md
git commit -m "docs: record remote smoke matrix aggregate artifact"
```
