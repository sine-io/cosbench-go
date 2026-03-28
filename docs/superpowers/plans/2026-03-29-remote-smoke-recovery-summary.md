# Remote Smoke Recovery Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a dedicated `remote-smoke-recovery-summary` artifact to the `Remote Smoke Recovery` workflow.

**Architecture:** Keep the workflow’s raw artifact upload and job summary behavior unchanged. Add one small builder script that copies `summary.json` and `summary.md` into a stable directory, then upload that directory as a second artifact.

**Tech Stack:** GitHub Actions YAML, Python builder script, Python workflow contract test

---

### Task 1: Add Failing Tests For The New Summary Artifact Contract

**Files:**
- Modify: `scripts/test_remote_smoke_recovery_workflow.py`
- Create: `scripts/test_build_remote_smoke_recovery_summary.py`

- [ ] **Step 1: Extend the workflow-contract test**

Require:
- the workflow invokes a recovery summary builder script
- the workflow uploads `remote-smoke-recovery-summary`

- [ ] **Step 2: Add a focused builder-script test**

Require:
- given `.artifacts/remote-smoke/summary.json`
- and `.artifacts/remote-smoke/summary.md`
- the script creates `.artifacts/remote-smoke-recovery-summary/` with both files

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_workflow.py scripts/test_build_remote_smoke_recovery_summary.py -q
```

Expected:
- failure because the builder script and summary artifact step do not exist yet

### Task 2: Add The Builder Script And Workflow Upload

**Files:**
- Create: `scripts/build_remote_smoke_recovery_summary.py`
- Modify: `.github/workflows/remote-smoke-recovery.yml`

- [ ] **Step 1: Add the builder script**

Implement:
- input: `.artifacts/remote-smoke/`
- output: `.artifacts/remote-smoke-recovery-summary/`
- copy:
  - `summary.json`
  - `summary.md`

- [ ] **Step 2: Invoke the builder script from the workflow**

Run it after the main smoke step and before artifact upload.

- [ ] **Step 3: Add the summary artifact upload**

Artifact name:
- `remote-smoke-recovery-summary`

- [ ] **Step 4: Re-run the focused tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_recovery_workflow.py scripts/test_build_remote_smoke_recovery_summary.py -q
```

Expected:
- both tests pass

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-recovery.yml scripts/build_remote_smoke_recovery_summary.py scripts/test_remote_smoke_recovery_workflow.py scripts/test_build_remote_smoke_recovery_summary.py
git commit -m "feat: add remote smoke recovery summary artifact"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the new summary artifact briefly**

Add one short note that the recovery workflow now emits a dedicated summary artifact.

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
git add README.md docs/superpowers/specs/2026-03-29-remote-smoke-recovery-summary-design.md docs/superpowers/plans/2026-03-29-remote-smoke-recovery-summary.md
git commit -m "docs: record remote smoke recovery summary"
```
