# Smoke Ready Upgrade Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expand `smoke-ready` so it reports the repository’s full smoke workflow surface and separates local, remote happy-path, and remote recovery readiness.

**Architecture:** Keep `smoke_ready.py` as a lightweight availability checker. Replace the single workflow constant with a small workflow map, add richer summary fields, and add focused tests for both text and JSON output using the existing mock environment hooks.

**Tech Stack:** Python helper script, Python tests, existing Makefile entrypoints

---

### Task 1: Add Failing Tests For Expanded Smoke Ready Output

**Files:**
- Create: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Add a JSON-mode test**

Cover:
- all smoke workflows appear in `payload["workflows"]["present"]`
- `summary` includes:
  - `local_env_ready`
  - `local_workflow_ready`
  - `remote_happy_ready`
  - `remote_recovery_ready`
  - `ready`

- [ ] **Step 2: Add a text-mode test**

Cover:
- output includes all smoke workflow names
- summary section includes the new readiness labels

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_ready.py -q
```

Expected:
- failure because the helper still assumes only `Smoke Local`

### Task 2: Expand The Smoke Ready Helper

**Files:**
- Modify: `scripts/smoke_ready.py`

- [ ] **Step 1: Replace the single workflow constant with a named workflow set**

Track:
- `Smoke Local`
- `Remote Smoke Local`
- `Remote Smoke Matrix`
- `Remote Smoke Recovery`
- `Remote Smoke Recovery Matrix`

- [ ] **Step 2: Add richer summary fields**

Implement:
- `local_env_ready`
- `local_workflow_ready`
- `remote_happy_ready`
- `remote_recovery_ready`
- `ready`

- [ ] **Step 3: Update text rendering**

Ensure the human-readable output shows:
- all workflows
- the expanded summary fields

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
git commit -m "feat: expand smoke ready coverage"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README wording for `smoke-ready`**

Describe that it now summarizes local, remote happy-path, and remote recovery readiness.

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
git add README.md docs/superpowers/specs/2026-03-29-smoke-ready-upgrade-design.md docs/superpowers/plans/2026-03-29-smoke-ready-upgrade.md
git commit -m "docs: record smoke ready upgrade"
```
