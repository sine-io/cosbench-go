# Smoke Ready Legacy Live Matrix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Teach `smoke-ready` and `smoke-ready-json` to report `Legacy Live Compare Matrix` as a separate readiness and latest-run signal.

**Architecture:** Keep the helper structure intact. Extend the workflow list with `Legacy Live Compare Matrix`, add two new summary fields for matrix readiness/latest success, then update the existing tests and README wording to cover the expanded readiness surface.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing tests for legacy live matrix reporting

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Extend the mocked workflow set and latest-run payload with `Legacy Live Compare Matrix`.
Add assertions for:
- `workflows.present["Legacy Live Compare Matrix"]`
- `workflows.latest["Legacy Live Compare Matrix"]`
- `summary.legacy_live_matrix_ready`
- `summary.legacy_live_matrix_latest_success`
- text output containing `Legacy Live Compare Matrix`

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper does not yet track the matrix workflow.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to expose the new workflow and summary fields.

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the operator-facing note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that `smoke-ready` now includes both the single-run and matrix legacy live compare workflows.

- [ ] **Step 2: Confirm terminology**

Keep the README wording aligned with the exact summary field names and avoid collapsing the two legacy signals together.

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
- only readiness observability for already-landed workflows changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-legacy-live-matrix-design.md docs/superpowers/plans/2026-03-29-smoke-ready-legacy-live-matrix.md
git commit -m "feat: add legacy live compare matrix to smoke-ready"
```
