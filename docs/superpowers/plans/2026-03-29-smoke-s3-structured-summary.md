# Smoke S3 Structured Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `Smoke S3` and `Smoke S3 Matrix` produce structured smoke result summaries in addition to raw output text.

**Architecture:** Add one small Python summary script that classifies raw smoke output, wire it into the single-run and matrix workflows, then update the matrix aggregate script and tests to consume structured row summaries instead of treating every row as only present or missing.

**Tech Stack:** GitHub Actions YAML, Python helper scripts, pytest

---

### Task 1: Add failing tests for structured smoke summaries

**Files:**
- Create: `scripts/test_summarize_smoke_s3_output.py`
- Modify: `scripts/test_smoke_s3_workflow.py`
- Modify: `scripts/test_smoke_s3_matrix_workflow.py`
- Modify: `scripts/test_aggregate_smoke_s3_matrix.py`

- [ ] **Step 1: Write the failing tests**

Require:
- a summary script that maps raw smoke output to `executed`, `skipped`, or `failed`
- `Smoke S3` workflow to call that script and upload the structured summary
- `Smoke S3 Matrix` workflow to do the same per row
- aggregate script to emit row statuses from structured summaries rather than only `present`

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_summarize_smoke_s3_output.py scripts/test_smoke_s3_workflow.py scripts/test_smoke_s3_matrix_workflow.py scripts/test_aggregate_smoke_s3_matrix.py -q`
Expected: FAIL because the summary script and workflow wiring do not exist yet.

- [ ] **Step 3: Write minimal implementation**

Create/update:
- `scripts/summarize_smoke_s3_output.py`
- `.github/workflows/smoke-s3.yml`
- `.github/workflows/smoke-s3-matrix.yml`
- `scripts/aggregate_smoke_s3_matrix.py`

- [ ] **Step 4: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_summarize_smoke_s3_output.py scripts/test_smoke_s3_workflow.py scripts/test_smoke_s3_matrix_workflow.py scripts/test_aggregate_smoke_s3_matrix.py -q`
Expected: PASS

### Task 2: Update README note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that `Smoke S3` and `Smoke S3 Matrix` now publish structured summary artifacts in addition to raw text output.

- [ ] **Step 2: Confirm wording**

Keep the wording focused on emitted artifacts, not on `smoke-ready` behavior.

### Task 3: Verify and commit

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_summarize_smoke_s3_output.py scripts/test_smoke_s3_workflow.py scripts/test_smoke_s3_matrix_workflow.py scripts/test_aggregate_smoke_s3_matrix.py -q`
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
- only evidence/observability around existing real-endpoint smoke changed

- [ ] **Step 5: Commit**

```bash
git add scripts/summarize_smoke_s3_output.py scripts/test_summarize_smoke_s3_output.py scripts/test_smoke_s3_workflow.py scripts/test_smoke_s3_matrix_workflow.py scripts/test_aggregate_smoke_s3_matrix.py scripts/aggregate_smoke_s3_matrix.py .github/workflows/smoke-s3.yml .github/workflows/smoke-s3-matrix.yml README.md docs/superpowers/specs/2026-03-29-smoke-s3-structured-summary-design.md docs/superpowers/plans/2026-03-29-smoke-s3-structured-summary.md
git commit -m "feat: add structured smoke s3 summaries"
```
