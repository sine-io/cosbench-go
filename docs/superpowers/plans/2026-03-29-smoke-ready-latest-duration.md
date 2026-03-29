# Smoke Ready Latest Duration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add latest workflow duration reporting to `smoke-ready` and `smoke-ready-json`.

**Architecture:** Capture `startedAt` and `updatedAt` in the internal latest-run normalization, compute integer duration seconds, and expose only the derived `*_latest_duration_seconds` fields through the summary block. Keep the change reporting-only.

**Tech Stack:** Python helper scripts, pytest, JSON Schema, Markdown

---

### Task 1: Add failing coverage

**Files:**
- Modify: `scripts/test_smoke_ready.py`
- Modify: `scripts/test_smoke_ready_schema.py`
- Modify: `scripts/test_validate_smoke_ready_schema.py`

- [ ] **Step 1: Write failing tests**

Add mocked `startedAt` / `updatedAt` values and assert:
- summary `*_latest_duration_seconds` exists for all latest-evidence surfaces
- representative values match the mocked elapsed seconds

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: FAIL because summary does not yet expose duration fields.

### Task 2: Implement duration reporting

**Files:**
- Modify: `scripts/smoke_ready.py`
- Modify: `docs/smoke-ready.schema.json`

- [ ] **Step 1: Write minimal implementation**

Add:
- internal normalization for `started_at` and `updated_at`
- helper to compute duration seconds
- summary `*_latest_duration_seconds` fields

- [ ] **Step 2: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Update docs**

Add short notes that `smoke-ready` now reports latest duration seconds in addition to the existing latest-evidence metadata.

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
git add scripts/smoke_ready.py scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py docs/smoke-ready.schema.json README.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-29-smoke-ready-latest-duration-design.md docs/superpowers/plans/2026-03-29-smoke-ready-latest-duration.md
git commit -m "feat: report smoke-ready latest durations"
```
