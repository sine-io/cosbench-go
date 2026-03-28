# Smoke Ready Run Status Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Teach `smoke-ready` to report the latest run status of all smoke workflows in addition to their presence.

**Architecture:** Keep `smoke_ready.py` as a read-only helper. Add one lightweight GitHub query for each smoke workflow’s latest run, surface those records under a new JSON section, and summarize remote happy-path and recovery success from the latest workflow results. Use mock env data for focused unit-style tests.

**Tech Stack:** Python helper script, GitHub CLI `gh run list --json ...`, Python tests with mocked environment inputs

---

### Task 1: Add Failing Tests For Latest Run Reporting

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Extend the JSON-mode test**

Require:
- `payload["workflows"]["latest"]` exists
- latest run records exist for all smoke workflows
- summary includes:
  - `remote_happy_latest_success`
  - `remote_recovery_latest_success`

- [ ] **Step 2: Extend the text-mode test**

Require:
- output includes a latest-runs section
- output includes latest status/conclusion lines for smoke workflows

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_ready.py -q
```

Expected:
- failure because latest-run reporting is not implemented yet

### Task 2: Add Latest Run Queries To The Helper

**Files:**
- Modify: `scripts/smoke_ready.py`

- [ ] **Step 1: Add workflow-run loading**

Implement:
- mocked path via env for tests
- real path via `gh run list --workflow <name> --limit 1 --json ...`

- [ ] **Step 2: Extend the payload**

Add:
- `workflows.latest`
- `summary.remote_happy_latest_success`
- `summary.remote_recovery_latest_success`

- [ ] **Step 3: Update text rendering**

Add a `Latest Runs` section showing:
- workflow name
- status / conclusion
- created time when available

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
git commit -m "feat: add smoke ready run status"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README wording**

Describe that `smoke-ready` now summarizes both workflow presence and latest remote smoke evidence status.

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
git add README.md docs/superpowers/specs/2026-03-29-smoke-ready-run-status-design.md docs/superpowers/plans/2026-03-29-smoke-ready-run-status.md
git commit -m "docs: record smoke ready run status"
```
