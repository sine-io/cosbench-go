# Smoke Ready Schema Version Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a stable schema version marker and a dedicated structural contract test for `smoke-ready`.

**Architecture:** Keep the current behavior intact. Add one top-level version field to the helper output, create one narrow contract test that validates structural presence of the current summary schema, and update the README with one short note that the JSON output is versioned.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing contract coverage

**Files:**
- Create: `scripts/test_smoke_ready_schema.py`

- [ ] **Step 1: Write the failing test**

Assert:
- top-level `schema_version == 1`
- top-level keys like `repo`, `required`, `local_env`, `workflows`, `summary`, `blockers` exist
- summary keys for result/source/url/artifact/created_at exist for all current latest evidence categories

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready_schema.py -q`
Expected: FAIL because `schema_version` does not exist yet.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to emit `schema_version: 1`.

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready_schema.py -q`
Expected: PASS

### Task 2: Update README note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that `smoke-ready-json` now includes a top-level schema version for machine consumers.

- [ ] **Step 2: Confirm wording**

Keep the wording scoped to the JSON interface, not to workflow behavior.

### Task 3: Verify and commit

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py -q`
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
- only schema usability and regression protection changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready_schema.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-schema-version-design.md docs/superpowers/plans/2026-03-29-smoke-ready-schema-version.md
git commit -m "feat: version smoke-ready schema"
```
