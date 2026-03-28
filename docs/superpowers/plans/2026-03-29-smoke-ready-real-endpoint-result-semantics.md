# Smoke Ready Real Endpoint Result Semantics Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `smoke-ready` distinguish real-endpoint workflow success from actual executed live smoke success, and expose explicit real-endpoint result states.

**Architecture:** Extend latest-run metadata with run ids, fetch or download just enough artifact data for `Smoke S3` and `Smoke S3 Matrix`, derive normalized result states from the smoke outputs, then surface those states and corrected success booleans in JSON and text output.

**Tech Stack:** Python helper script, pytest, GitHub CLI metadata and artifact downloads

---

### Task 1: Add failing tests for real-endpoint result semantics

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Add mocked real-endpoint detail payloads where:
- `Smoke S3` latest run concludes `success` but its output shows only skipped smoke tests
- `Smoke S3 Matrix` latest run concludes `success` but both rows in aggregate summary show skipped smoke output

Assert:
- `real_endpoint_latest_success is False`
- `real_endpoint_matrix_latest_success is False`
- `real_endpoint_latest_result == "skipped"`
- `real_endpoint_matrix_latest_result == "skipped"`
- text output includes the two new result lines

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper still maps workflow success directly to real-endpoint latest success.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- fetch or mock artifact-backed real-endpoint details
- derive normalized real-endpoint result states
- compute corrected success booleans

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that real-endpoint smoke readiness now distinguishes workflow availability from the latest smoke execution result.

- [ ] **Step 2: Confirm terminology**

Ensure README wording matches `real_endpoint_latest_result` and `real_endpoint_matrix_latest_result`.

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

Expected: PASS, with corrected real-endpoint result semantics.

- [ ] **Step 5: Compare against overall project goal**

Confirm this stays within the project goal:
- no runtime or protocol changes
- no non-S3 backend expansion
- only readiness signal accuracy changes

- [ ] **Step 6: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-real-endpoint-result-semantics-design.md docs/superpowers/plans/2026-03-29-smoke-ready-real-endpoint-result-semantics.md
git commit -m "feat: refine real-endpoint smoke-ready results"
```
