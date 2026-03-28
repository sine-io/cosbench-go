# Smoke Ready Legacy Live Result Semantics Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `smoke-ready` distinguish legacy live workflow success from actual executed live workload success, and expose explicit legacy live result states.

**Architecture:** Extend the helper’s latest-run metadata with run ids, fetch step-level job details for the two legacy workflows only, derive normalized legacy result states from the execution step conclusions, then surface those states and corrected success booleans in JSON and text output.

**Tech Stack:** Python helper script, pytest, GitHub CLI run metadata

---

### Task 1: Add failing tests for legacy live result semantics

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Add mocked legacy workflow detail payloads where:
- `Legacy Live Compare` has workflow conclusion `success` but execution step conclusion `skipped`
- `Legacy Live Compare Matrix` has workflow conclusion `success` but both row execution steps `skipped`

Assert:
- `legacy_live_latest_success is False`
- `legacy_live_matrix_latest_success is False`
- `legacy_live_latest_result == "skipped"`
- `legacy_live_matrix_latest_result == "skipped"`
- text output includes the two result lines

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper still maps workflow success directly to legacy latest success.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- carry `databaseId` in latest-run metadata
- load run details for legacy workflows
- derive normalized result states
- compute corrected success booleans

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that legacy live signals in `smoke-ready` now distinguish workflow availability from actual executed legacy live results.

- [ ] **Step 2: Confirm terminology**

Ensure README wording matches the new `legacy_live_latest_result` and `legacy_live_matrix_latest_result` field names.

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

- [ ] **Step 4: Run readiness helpers**

Run:
- `make --no-print-directory smoke-ready`
- `make --no-print-directory smoke-ready-json`

Expected: PASS, with legacy result states visible.

- [ ] **Step 5: Compare against overall project goal**

Confirm this stays within the project goal:
- no runtime or protocol changes
- no non-S3 backend expansion
- only readiness signal accuracy changes

- [ ] **Step 6: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-legacy-live-result-semantics-design.md docs/superpowers/plans/2026-03-29-smoke-ready-legacy-live-result-semantics.md
git commit -m "feat: refine legacy live smoke-ready results"
```
