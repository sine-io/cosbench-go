# Workflow Upload Artifact V7 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade all workflow `upload-artifact` uses to `v7.0.0` and lock the version with a generic workflow contract test.

**Architecture:** Keep workflow behavior unchanged and only update action pins. Add one scan-based Python test that inspects `.github/workflows/*.yml` for `actions/upload-artifact` references and requires them all to be `v7.0.0`.

**Tech Stack:** GitHub Actions YAML, Python workflow contract test

---

### Task 1: Add A Failing Generic Version Test

**Files:**
- Create: `scripts/test_workflow_action_versions.py`

- [ ] **Step 1: Add a test that scans workflow files for `upload-artifact` pins**

Require:
- every `actions/upload-artifact@...` reference equals `v7.0.0`

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_workflow_action_versions.py -q
```

Expected:
- failure because some workflows still use `v6.0.0`

### Task 2: Upgrade Workflow Action Pins

**Files:**
- Modify: `.github/workflows/compare-local.yml`
- Modify: `.github/workflows/remote-smoke-local.yml`
- Modify: `.github/workflows/remote-smoke-matrix.yml`
- Modify: `.github/workflows/smoke-local.yml`

- [ ] **Step 1: Upgrade all `actions/upload-artifact` references to `v7.0.0`**

- [ ] **Step 2: Re-run the focused version test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_workflow_action_versions.py -q
```

Expected:
- the version test passes

- [ ] **Step 3: Commit the workflow slice**

Run:
```bash
git add .github/workflows scripts/test_workflow_action_versions.py
git commit -m "chore: upgrade workflow artifact uploads"
```

### Task 3: Final Verification

**Files:**
- No additional files

- [ ] **Step 1: Run the full Go test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all Go packages pass

- [ ] **Step 2: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly
