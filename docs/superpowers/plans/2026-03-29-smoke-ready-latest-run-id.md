# Smoke Ready Latest Run ID Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add latest-run-id reporting to `smoke-ready` and `smoke-ready-json`.

**Architecture:** Reuse the already-normalized `database_id` in `workflows.latest` and thread it into the summary block as `*_latest_run_id` for each latest-evidence surface. Keep the change reporting-only.

**Tech Stack:** Python helper scripts, pytest, JSON Schema, Markdown

---

### Task 1: Add failing coverage

**Files:**
- Modify: `scripts/test_smoke_ready.py`
- Modify: `scripts/test_smoke_ready_schema.py`
- Modify: `scripts/test_validate_smoke_ready_schema.py`

- [ ] **Step 1: Write failing tests**

Add assertions for:
- `workflows.latest[*].database_id` already present in mocks
- summary `*_latest_run_id` fields for all latest-evidence surfaces

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: FAIL because summary does not yet expose `*_latest_run_id`.

### Task 2: Implement run-id reporting

**Files:**
- Modify: `scripts/smoke_ready.py`
- Modify: `docs/smoke-ready.schema.json`

- [ ] **Step 1: Write minimal implementation**

Add `*_latest_run_id` fields for:
- `real_endpoint`
- `real_endpoint_matrix`
- `schema_validation`
- `legacy_live`
- `legacy_live_matrix`
- `remote_happy`
- `remote_recovery`

- [ ] **Step 2: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Update docs**

Add short notes that `smoke-ready` now reports latest run ids in addition to existing latest-evidence metadata.

- [ ] **Step 2: Run targeted tests**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py scripts/test_smoke_ready_validate_workflow.py -q`
Expected: PASS

- [ ] **Step 3: Run repo tests**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 4: Run full build**

Run: `go build ./...`
Expected: PASS

- [ ] **Step 5: Compare against overall project goal**

Confirm this stays within the project goal:
- no runtime or protocol changes
- no non-S3 backend expansion
- only reporting metadata changed

- [ ] **Step 6: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py docs/smoke-ready.schema.json README.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-29-smoke-ready-latest-run-id-design.md docs/superpowers/plans/2026-03-29-smoke-ready-latest-run-id.md
git commit -m "feat: report smoke-ready latest run ids"
```
