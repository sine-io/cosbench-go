# Smoke Ready Latest Age Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add latest-evidence age reporting to `smoke-ready` and `smoke-ready-json`.

**Architecture:** Reuse the existing `generated_at` plus each surface’s `*_latest_created_at` to compute integer `*_latest_age_seconds`. Keep the change reporting-only and avoid changing workflow metadata collection.

**Tech Stack:** Python helper scripts, pytest, JSON Schema, Markdown

---

### Task 1: Add failing coverage

**Files:**
- Modify: `scripts/test_smoke_ready.py`
- Modify: `scripts/test_smoke_ready_schema.py`

- [ ] **Step 1: Write failing tests**

Add assertions for:
- `*_latest_age_seconds` keys in the summary
- representative values computed from `payload["generated_at"]` and mocked `*_latest_created_at`

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py -q`
Expected: FAIL because summary does not yet expose age fields.

### Task 2: Implement age reporting

**Files:**
- Modify: `scripts/smoke_ready.py`
- Modify: `docs/smoke-ready.schema.json`

- [ ] **Step 1: Write minimal implementation**

Add:
- helper to parse timestamps
- helper to compute age seconds from `generated_at`
- summary `*_latest_age_seconds` fields

- [ ] **Step 2: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Update docs**

Add short notes that `smoke-ready` now reports latest age seconds in addition to the other latest-evidence metadata.

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
git add scripts/smoke_ready.py scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py docs/smoke-ready.schema.json README.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-29-smoke-ready-latest-age-design.md docs/superpowers/plans/2026-03-29-smoke-ready-latest-age.md
git commit -m "feat: report smoke-ready latest ages"
```
