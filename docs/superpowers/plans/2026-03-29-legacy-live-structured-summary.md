# Legacy Live Structured Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `Legacy Live Compare` and `Legacy Live Compare Matrix` publish structured result summaries in addition to their existing raw summary files.

**Architecture:** Add one small Python summary script that normalizes the current raw legacy summary into a stable `result` shape, wire it into the single-run and matrix workflows, then update the matrix aggregate script and tests so it consumes structured row summaries rather than inferring everything from workflow behavior.

**Tech Stack:** GitHub Actions YAML, Python helper scripts, pytest

---

### Task 1: Add failing tests for structured legacy live summaries

**Files:**
- Create: `scripts/test_summarize_legacy_live_compare.py`
- Modify: `scripts/test_legacy_live_compare_workflow.py`
- Modify: `scripts/test_legacy_live_compare_matrix_workflow.py`
- Modify: `scripts/test_aggregate_legacy_live_compare_matrix.py`

- [ ] **Step 1: Write the failing tests**

Require:
- a summary script that maps legacy raw summary JSON to `executed/skipped/failed`
- `Legacy Live Compare` workflow to call that script and upload the normalized summary
- `Legacy Live Compare Matrix` row workflows to do the same
- aggregate script to emit row statuses from normalized summaries rather than only reading legacy raw summaries

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_summarize_legacy_live_compare.py scripts/test_legacy_live_compare_workflow.py scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py -q`
Expected: FAIL because the summary script and workflow wiring do not exist yet.

- [ ] **Step 3: Write minimal implementation**

Create/update:
- `scripts/summarize_legacy_live_compare.py`
- `.github/workflows/legacy-live-compare.yml`
- `.github/workflows/legacy-live-compare-matrix.yml`
- `scripts/aggregate_legacy_live_compare_matrix.py`

- [ ] **Step 4: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_summarize_legacy_live_compare.py scripts/test_legacy_live_compare_workflow.py scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py -q`
Expected: PASS

### Task 2: Update README note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that legacy live workflows now publish structured normalized summaries in addition to the raw CLI summary.

- [ ] **Step 2: Confirm wording**

Keep the wording focused on emitted artifacts, not downstream `smoke-ready` behavior yet.

### Task 3: Verify and commit

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_summarize_legacy_live_compare.py scripts/test_legacy_live_compare_workflow.py scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py -q`
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
- only evidence/observability around existing legacy live workflows changed

- [ ] **Step 5: Commit**

```bash
git add scripts/summarize_legacy_live_compare.py scripts/test_summarize_legacy_live_compare.py scripts/test_legacy_live_compare_workflow.py scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py scripts/aggregate_legacy_live_compare_matrix.py .github/workflows/legacy-live-compare.yml .github/workflows/legacy-live-compare-matrix.yml README.md docs/superpowers/specs/2026-03-29-legacy-live-structured-summary-design.md docs/superpowers/plans/2026-03-29-legacy-live-structured-summary.md
git commit -m "feat: add structured legacy live summaries"
```
