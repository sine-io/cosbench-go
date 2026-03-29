# Smoke Ready Current Head Branch Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add current checkout branch reporting to `smoke-ready` and `smoke-ready-json`.

**Architecture:** Resolve the current branch once, expose it as top-level `current_head_branch`, and keep the rest of the evidence-selection logic unchanged. This is a reporting-only addition.

**Tech Stack:** Python helper scripts, pytest, JSON Schema, Markdown

---

### Task 1: Add failing coverage

**Files:**
- Modify: `scripts/test_smoke_ready.py`
- Modify: `scripts/test_smoke_ready_schema.py`
- Modify: `scripts/test_validate_smoke_ready_schema.py`

- [ ] **Step 1: Write failing tests**

Add:
- `SMOKE_READY_MOCK_CURRENT_HEAD_BRANCH`
- assertions for top-level `current_head_branch`

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: FAIL because helper does not yet expose the new field.

### Task 2: Implement current branch reporting

**Files:**
- Modify: `scripts/smoke_ready.py`
- Modify: `docs/smoke-ready.schema.json`

- [ ] **Step 1: Write minimal implementation**

Add:
- current-branch resolution helper
- top-level `current_head_branch`

- [ ] **Step 2: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Update docs**

Add short notes that `smoke-ready` now reports current checkout branch alongside current head SHA.

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
git add scripts/smoke_ready.py scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py docs/smoke-ready.schema.json README.md AGENTS.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-29-smoke-ready-current-head-branch-design.md docs/superpowers/plans/2026-03-29-smoke-ready-current-head-branch.md
git commit -m "feat: report smoke-ready current head branch"
```
