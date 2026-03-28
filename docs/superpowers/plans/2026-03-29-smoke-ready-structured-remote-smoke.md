# Smoke Ready Structured Remote Smoke Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `smoke-ready` consume structured remote smoke and remote recovery summaries first, and expose explicit remote result states.

**Architecture:** Extend the helper to download the already-produced remote summary artifacts, compare latest run timestamps between single-run and matrix workflows, derive normalized remote result states from those structured payloads, and keep workflow-conclusion logic only as a compatibility fallback.

**Tech Stack:** Python helper script, pytest, GitHub CLI artifact downloads

---

### Task 1: Add failing tests for structured remote smoke summaries

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Change the mocked remote details so:
- `Remote Smoke Local` and `Remote Smoke Matrix` provide structured summaries with explicit `overall`
- `Remote Smoke Recovery` and `Remote Smoke Recovery Matrix` do the same

Assert:
- `remote_happy_latest_success` and `remote_recovery_latest_success` come from the structured result
- `remote_happy_latest_result` and `remote_recovery_latest_result` exist
- text output includes the two new result lines

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper still uses workflow conclusions only.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- download structured remote summary artifacts
- derive normalized remote result states
- compute corrected success booleans

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that remote smoke readiness now also consumes structured remote summary artifacts.

- [ ] **Step 2: Confirm wording**

Keep the wording focused on evidence consumption, not workflow behavior changes.

### Task 3: Verify and commit

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

- [ ] **Step 2: Run repo tests**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 3: Run full build**

Run: `go build ./...`
Expected: PASS

- [ ] **Step 4: Compare against overall project goal**

Confirm this stays within the project goal:
- no runtime or protocol changes
- no non-S3 backend expansion
- only readiness evidence consumption changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-structured-remote-smoke-design.md docs/superpowers/plans/2026-03-29-smoke-ready-structured-remote-smoke.md
git commit -m "feat: use structured remote smoke summaries in smoke-ready"
```
