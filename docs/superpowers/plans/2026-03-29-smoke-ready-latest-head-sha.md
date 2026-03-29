# Smoke Ready Latest Head SHA Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add latest evidence head-SHA reporting to `smoke-ready` and `smoke-ready-json`.

**Architecture:** Extend the existing latest-run normalization with `head_sha`, then thread that value into the summary block for each latest-evidence surface. Keep the change reporting-only.

**Tech Stack:** Python helper scripts, pytest, JSON Schema, Markdown

---

### Task 1: Add failing coverage

**Files:**
- Modify: `scripts/test_smoke_ready.py`
- Modify: `scripts/test_smoke_ready_schema.py`
- Modify: `scripts/test_validate_smoke_ready_schema.py`

- [ ] **Step 1: Write failing tests**

Add mocked `headSha` values and assert:
- `workflows.latest[*].head_sha` exists
- summary includes `*_latest_head_sha` for all latest-evidence surfaces

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: FAIL because `smoke_ready.py` does not yet expose `head_sha`.

### Task 2: Implement head-SHA reporting

**Files:**
- Modify: `scripts/smoke_ready.py`
- Modify: `docs/smoke-ready.schema.json`

- [ ] **Step 1: Write minimal implementation**

Add:
- `head_sha` to latest-run normalization
- `headSha` to `gh run list --json ...`
- summary `*_latest_head_sha` fields

- [ ] **Step 2: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Update docs**

Add short notes that `smoke-ready` now reports latest head SHA in addition to the other latest-evidence metadata.

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
git add scripts/smoke_ready.py scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py docs/smoke-ready.schema.json README.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-29-smoke-ready-latest-head-sha-design.md docs/superpowers/plans/2026-03-29-smoke-ready-latest-head-sha.md
git commit -m "feat: report smoke-ready latest head shas"
```
