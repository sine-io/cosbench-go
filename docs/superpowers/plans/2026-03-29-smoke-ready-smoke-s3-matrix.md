# Smoke Ready Smoke S3 Matrix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Teach `smoke-ready` and `smoke-ready-json` to report `Smoke S3 Matrix` as a separate real-endpoint matrix signal.

**Architecture:** Keep the helper structure unchanged. Extend the tracked workflow set with `Smoke S3 Matrix`, add one matrix-readiness bit and one matrix-latest-success bit, then update the existing tests and README wording to cover the expanded real-endpoint surface.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing tests for Smoke S3 Matrix reporting

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Extend the mocked workflow set and latest-run payload with `Smoke S3 Matrix`.
Add assertions for:
- `workflows.present["Smoke S3 Matrix"]`
- `workflows.latest["Smoke S3 Matrix"]`
- `summary.real_endpoint_matrix_ready`
- `summary.real_endpoint_matrix_latest_success`
- text output containing `Smoke S3 Matrix`

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

Add one short note that `smoke-ready` now includes both `Smoke S3` and `Smoke S3 Matrix` as separate real-endpoint signals.

- [ ] **Step 2: Confirm terminology**

Keep the README wording aligned with `real_endpoint_ready` and `real_endpoint_matrix_ready`.

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

Expected: PASS, with `Smoke S3 Matrix` visible.

- [ ] **Step 5: Compare against overall project goal**

Confirm this stays within the project goal:
- no runtime or protocol changes
- no non-S3 backend expansion
- only readiness observability for already-landed workflows changed

- [ ] **Step 6: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-smoke-s3-matrix-design.md docs/superpowers/plans/2026-03-29-smoke-ready-smoke-s3-matrix.md
git commit -m "feat: add smoke s3 matrix to smoke-ready"
```
