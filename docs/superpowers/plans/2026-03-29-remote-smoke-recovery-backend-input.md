# Remote Smoke Recovery Backend Input Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Parameterize the dedicated `Remote Smoke Recovery` workflow so it can run either `s3` or `sio` while keeping `scenario=recovery` fixed.

**Architecture:** Keep the existing workflow file, job structure, artifact behavior, and recovery summary builder unchanged. Add only a `backend` workflow input, thread it into the existing run step, and extend the workflow contract test to lock the new interface.

**Tech Stack:** GitHub Actions YAML, Python workflow contract test, README documentation

---

### Task 1: Add A Failing Workflow Contract Test

**Files:**
- Modify: `scripts/test_remote_smoke_recovery_workflow.py`

- [ ] **Step 1: Extend the workflow-contract test**

Require:
- `inputs.backend` exists
- default backend is `s3`
- the run step uses `${{ inputs.backend }}`
- `SMOKE_REMOTE_LOCAL_SCENARIO=recovery` remains fixed

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_workflow.py -q
```

Expected:
- failure because the workflow still hardcodes `s3`

### Task 2: Parameterize The Workflow

**Files:**
- Modify: `.github/workflows/remote-smoke-recovery.yml`

- [ ] **Step 1: Add `workflow_dispatch.inputs.backend`**

Use:
- default `s3`
- no scenario input

- [ ] **Step 2: Thread `backend` into the run step**

Keep:
- `SMOKE_REMOTE_LOCAL_SCENARIO=recovery`

- [ ] **Step 3: Re-run the focused workflow test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_workflow.py -q
```

Expected:
- the workflow test passes

- [ ] **Step 4: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-recovery.yml scripts/test_remote_smoke_recovery_workflow.py
git commit -m "feat: parameterize remote smoke recovery backend"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add README examples**

Document:
- `gh workflow run "Remote Smoke Recovery" --repo sine-io/cosbench-go`
- `gh workflow run "Remote Smoke Recovery" --repo sine-io/cosbench-go -f backend=sio`

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
git add README.md docs/superpowers/specs/2026-03-29-remote-smoke-recovery-backend-input-design.md docs/superpowers/plans/2026-03-29-remote-smoke-recovery-backend-input.md
git commit -m "docs: record remote smoke recovery backend input"
```
