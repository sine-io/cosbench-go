# Remote Smoke Recovery Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a dedicated manual GitHub workflow for the `s3 + recovery` remote smoke path.

**Architecture:** Keep the existing helper and `Remote Smoke Local` workflow unchanged. Add one new `workflow_dispatch` workflow that hardcodes `backend=s3` and `scenario=recovery`, plus a lightweight contract test and one README invocation example.

**Tech Stack:** GitHub Actions YAML, Python workflow contract test, README documentation

---

### Task 1: Add A Failing Workflow Contract Test

**Files:**
- Create: `scripts/test_remote_smoke_recovery_workflow.py`

- [ ] **Step 1: Add a workflow-contract test**

Cover:
- workflow file exists
- it is `workflow_dispatch`
- it runs `SMOKE_REMOTE_LOCAL_BACKEND=s3`
- it runs `SMOKE_REMOTE_LOCAL_SCENARIO=recovery`
- it uploads `.artifacts/remote-smoke/`
- it writes the smoke summary

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_workflow.py -q
```

Expected:
- failure because the workflow file does not exist yet

### Task 2: Add The Recovery Workflow

**Files:**
- Create: `.github/workflows/remote-smoke-recovery.yml`

- [ ] **Step 1: Add the workflow file**

Use:
- `name: Remote Smoke Recovery`
- `workflow_dispatch`

- [ ] **Step 2: Fix the command to `backend=s3`, `scenario=recovery`**

Run:
- `SMOKE_REMOTE_LOCAL_BACKEND=s3`
- `SMOKE_REMOTE_LOCAL_SCENARIO=recovery`
- `GO=go make --no-print-directory smoke-remote-local`

- [ ] **Step 3: Add artifact upload and summary export**

Match the current `Remote Smoke Local` conventions.

- [ ] **Step 4: Re-run the focused contract test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_workflow.py -q
```

Expected:
- the workflow contract test passes

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-recovery.yml scripts/test_remote_smoke_recovery_workflow.py
git commit -m "feat: add remote smoke recovery workflow"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add one invocation example**

Document:
- `gh workflow run "Remote Smoke Recovery" --repo sine-io/cosbench-go`

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
git add README.md docs/superpowers/specs/2026-03-29-remote-smoke-recovery-workflow-design.md docs/superpowers/plans/2026-03-29-remote-smoke-recovery-workflow.md
git commit -m "docs: record remote smoke recovery workflow"
```
