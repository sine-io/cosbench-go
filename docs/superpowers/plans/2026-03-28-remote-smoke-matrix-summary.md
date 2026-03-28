# Remote Smoke Matrix Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a combined aggregate summary to the non-blocking `Remote Smoke Matrix` workflow without changing helper behavior or default CI.

**Architecture:** Keep the existing four matrix rows intact. Add a small Python aggregation script and a new `aggregate` workflow job that downloads per-row artifacts, aggregates `summary.json` files into one combined report, and writes that report into the GitHub job summary.

**Tech Stack:** GitHub Actions YAML, Python aggregation script, lightweight Python tests, existing remote smoke artifacts

---

### Task 1: Add Failing Tests For Matrix Aggregation Contract

**Files:**
- Modify: `scripts/test_remote_smoke_matrix_workflow.py`
- Create: `scripts/test_aggregate_remote_smoke_matrix.py`

- [ ] **Step 1: Extend the workflow-contract test for the aggregate job**

Cover:
- workflow contains `aggregate` job
- `aggregate` depends on the matrix job
- `aggregate` downloads `remote-smoke-*` artifacts
- `aggregate` calls the aggregation script

- [ ] **Step 2: Add a focused aggregation-script test**

Cover:
- multiple row summaries aggregate into one Markdown and JSON output
- missing row summaries are represented as `missing`

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py scripts/test_aggregate_remote_smoke_matrix.py -q
```

Expected:
- failures because no aggregate job or aggregation script exists yet

### Task 2: Add The Aggregation Script And Workflow Job

**Files:**
- Create: `scripts/aggregate_remote_smoke_matrix.py`
- Modify: `.github/workflows/remote-smoke-matrix.yml`

- [ ] **Step 1: Add the aggregation script**

Implement:
- input directory scan for `summary.json`
- combined JSON output
- combined Markdown output
- `missing` rows when expected artifacts are absent

- [ ] **Step 2: Add the `aggregate` workflow job**

Ensure:
- `needs: remote_smoke_matrix`
- `if: always()`
- artifact download step
- aggregation step
- GitHub summary export step

- [ ] **Step 3: Re-run the focused aggregation tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_matrix_workflow.py scripts/test_aggregate_remote_smoke_matrix.py -q
```

Expected:
- both tests pass

- [ ] **Step 4: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-matrix.yml scripts/aggregate_remote_smoke_matrix.py scripts/test_remote_smoke_matrix_workflow.py scripts/test_aggregate_remote_smoke_matrix.py
git commit -m "feat: add remote smoke matrix summary"
```

### Task 3: Verify Repository Surface And Document The Combined Summary

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the combined matrix summary**

Add one short note that the matrix workflow now emits a combined summary across the four supported combinations.

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

- [ ] **Step 4: Re-run the local smoke matrix**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
SMOKE_REMOTE_LOCAL_BACKEND=sio timeout 240s make --no-print-directory smoke-remote-local
SMOKE_REMOTE_LOCAL_BACKEND=sio SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- all four combinations still succeed locally

- [ ] **Step 5: Commit the docs slice**

Run:
```bash
git add README.md docs/superpowers/specs/2026-03-28-remote-smoke-matrix-summary-design.md docs/superpowers/plans/2026-03-28-remote-smoke-matrix-summary.md
git commit -m "docs: record remote smoke matrix summary"
```
