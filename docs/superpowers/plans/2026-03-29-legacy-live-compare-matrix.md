# Legacy Live Compare Matrix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a manual `Legacy Live Compare Matrix` workflow with per-row artifacts and one aggregate summary/artifact across the representative `s3` and `sio` legacy fixtures.

**Architecture:** Keep the existing single-run workflow intact. Add a separate matrix workflow that inlines the same preflight/render/run pattern for two fixed rows, then aggregate row artifacts with a small Python script modeled after the existing smoke matrix helpers.

**Tech Stack:** GitHub Actions YAML, Python aggregation script, pytest workflow/aggregation tests, Markdown docs

---

### Task 1: Add failing tests for the matrix workflow contract

**Files:**
- Create: `scripts/test_legacy_live_compare_matrix_workflow.py`
- Create: `scripts/test_aggregate_legacy_live_compare_matrix.py`

- [ ] **Step 1: Write the failing workflow-shape test**

Require:
- `name: Legacy Live Compare Matrix`
- two fixed rows for `s3` and `sio`
- per-row artifact names
- an `aggregate` job
- use of `actions/download-artifact@v8.0.1`
- call to `scripts/aggregate_legacy_live_compare_matrix.py`

- [ ] **Step 2: Write the failing aggregation test**

Require the aggregation script to:
- find row `summary.json`
- render `summary.md`
- write `summary.json`
- classify rows as `executed`, `skipped`, or `missing`

- [ ] **Step 3: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py -q`
Expected: FAIL because none of the new files exist yet.

- [ ] **Step 4: Implement the minimal workflow and aggregation script**

Create:
- `.github/workflows/legacy-live-compare-matrix.yml`
- `scripts/aggregate_legacy_live_compare_matrix.py`

- [ ] **Step 5: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py -q`
Expected: PASS

### Task 2: Document the new entrypoint

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add workflow usage**

Add one `gh workflow run "Legacy Live Compare Matrix"` example near the existing single-run legacy workflow example.

- [ ] **Step 2: Note skip semantics**

Add a short note that the matrix rows inherit the same clean `skipped` behavior when `COSBENCH_SMOKE_*` secrets are absent.

### Task 3: Verify and commit

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py -q`
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
- only workflow automation around existing legacy live validation

- [ ] **Step 5: Commit**

```bash
git add .github/workflows/legacy-live-compare-matrix.yml scripts/aggregate_legacy_live_compare_matrix.py scripts/test_legacy_live_compare_matrix_workflow.py scripts/test_aggregate_legacy_live_compare_matrix.py README.md docs/superpowers/specs/2026-03-29-legacy-live-compare-matrix-design.md docs/superpowers/plans/2026-03-29-legacy-live-compare-matrix.md
git commit -m "feat: add legacy live compare matrix workflow"
```
