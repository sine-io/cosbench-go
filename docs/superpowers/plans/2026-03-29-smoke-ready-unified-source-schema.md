# Smoke Ready Unified Source Schema Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add source fields for the remaining non-aggregated latest evidence categories so the `smoke-ready` summary schema is uniform.

**Architecture:** Keep current result, URL, and timestamp logic intact. Add four direct workflow-name mappings to the summary block, surface them in text mode, and update tests and README wording to describe the completed schema.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing tests for the remaining source fields

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Assert that JSON output includes:
- `real_endpoint_latest_source`
- `real_endpoint_matrix_latest_source`
- `legacy_live_latest_source`
- `legacy_live_matrix_latest_source`

Also assert that text output includes the paired source labels.

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because those source fields do not yet exist.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to add the four source fields and print them in text mode.

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that every latest evidence category in `smoke-ready` now exposes result, source, URL, and timestamp.

- [ ] **Step 2: Confirm wording**

Keep the wording focused on schema consistency, not on new behavior.

### Task 3: Verify and commit

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
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
- only summary schema consistency changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-unified-source-schema-design.md docs/superpowers/plans/2026-03-29-smoke-ready-unified-source-schema.md
git commit -m "feat: add smoke-ready source schema fields"
```
