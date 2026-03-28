# Remote Smoke Download Artifact Upgrade Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade `actions/download-artifact` in `Remote Smoke Matrix` to the current official release and lock that version in tests.

**Architecture:** Keep the workflow structure unchanged. Update only the `download-artifact` step version and extend the existing workflow contract test so future edits do not regress to the Node20-targeting release.

**Tech Stack:** GitHub Actions YAML, Python workflow contract test

---

### Task 1: Add A Failing Version Assertion

**Files:**
- Modify: `scripts/test_remote_smoke_matrix_workflow.py`

- [ ] **Step 1: Add a version assertion for `download-artifact`**

Require:
- `uses: actions/download-artifact@v8.0.1`

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py -q
```

Expected:
- failure because the workflow still references `v6.0.0`

### Task 2: Upgrade The Workflow Step

**Files:**
- Modify: `.github/workflows/remote-smoke-matrix.yml`

- [ ] **Step 1: Upgrade `actions/download-artifact` to `v8.0.1`**

- [ ] **Step 2: Re-run the focused workflow-contract test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py -q
```

Expected:
- the test passes

- [ ] **Step 3: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-matrix.yml scripts/test_remote_smoke_matrix_workflow.py
git commit -m "chore: upgrade remote smoke artifact download"
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
