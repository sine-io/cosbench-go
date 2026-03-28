# Smoke Ready Legacy Live Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Teach `smoke-ready` and `smoke-ready-json` to report `Legacy Live Compare` as a separate readiness and latest-run signal.

**Architecture:** Keep the helper structure unchanged. Extend the tracked workflow list with `Legacy Live Compare`, add one new readiness bit and one latest-success bit, then update the existing tests and README to reflect the expanded readiness surface.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing tests for legacy live reporting

**Files:**
- Modify: `scripts/test_smoke_ready.py`
- Test: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Update the mocked workflow set and latest-run payload so they include `Legacy Live Compare`.
Add assertions for:
- `workflows.present["Legacy Live Compare"]`
- `workflows.latest["Legacy Live Compare"]`
- `summary.legacy_live_ready`
- `summary.legacy_live_latest_success`
- text output containing `Legacy Live Compare`

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper does not yet expose legacy live workflow state.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- track `Legacy Live Compare`
- compute the two new summary fields
- print them in text mode

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the operator-facing note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note near the `smoke-ready` description that the readiness view now also includes `Legacy Live Compare`.

- [ ] **Step 2: Read for consistency**

Confirm the README wording stays consistent with the helper’s summary field names and does not blur `Smoke S3` with legacy live compare.

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
- no protocol/runtime changes
- no non-S3 backend expansion
- only readiness/observability for existing workflows changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-legacy-live-design.md docs/superpowers/plans/2026-03-29-smoke-ready-legacy-live.md
git commit -m "feat: add legacy live compare to smoke-ready"
```
