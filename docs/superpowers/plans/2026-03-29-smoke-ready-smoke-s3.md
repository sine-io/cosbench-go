# Smoke Ready Smoke S3 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Teach `smoke-ready` to include the new `Smoke S3` workflow in workflow presence, latest-run reporting, and summary output.

**Architecture:** Keep the current `smoke_ready.py` helper structure, but extend the workflow set with `Smoke S3`, add one extra summary flag for latest real-endpoint smoke success, and update the existing tests to cover the new field and workflow row.

**Tech Stack:** Python helper script, Python tests, existing `gh run list` integration

---

### Task 1: Add Failing Tests For Smoke S3 Reporting

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Extend the JSON-mode test**

Require:
- `workflows.present["Smoke S3"]`
- `workflows.latest["Smoke S3"]`
- `summary.real_endpoint_latest_success`

- [ ] **Step 2: Extend the text-mode test**

Require:
- `Smoke S3` appears in workflows and latest-runs sections
- `Real Endpoint Latest Success` appears in the summary section

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_ready.py -q
```

Expected:
- failure because `Smoke S3` is not yet in the helper’s workflow list or summary

### Task 2: Extend The Helper For Smoke S3

**Files:**
- Modify: `scripts/smoke_ready.py`

- [ ] **Step 1: Add `Smoke S3` to the workflow set**

- [ ] **Step 2: Add `real_endpoint_latest_success` to summary output**

Set it from the latest `Smoke S3` run conclusion.

- [ ] **Step 3: Update text rendering**

Show:
- `Smoke S3` in the workflow list
- `Smoke S3` in latest runs
- `Real Endpoint Latest Success`

- [ ] **Step 4: Re-run the focused tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_ready.py -q
```

Expected:
- tests pass

- [ ] **Step 5: Commit the helper slice**

Run:
```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py
git commit -m "feat: add smoke ready smoke s3 status"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README wording**

Describe that `smoke-ready` now includes latest status for real-endpoint smoke as well.

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

- [ ] **Step 4: Smoke the helper in text and JSON modes**

Run:
```bash
make --no-print-directory smoke-ready
make --no-print-directory smoke-ready-json
```

Expected:
- both commands succeed

- [ ] **Step 5: Commit the docs slice**

Run:
```bash
git add README.md docs/superpowers/specs/2026-03-29-smoke-ready-smoke-s3-design.md docs/superpowers/plans/2026-03-29-smoke-ready-smoke-s3.md
git commit -m "docs: record smoke ready smoke s3 status"
```
