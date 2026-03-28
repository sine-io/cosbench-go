# Smoke Ready Source Fields Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add source attribution fields for the aggregated remote result summaries in `smoke-ready`.

**Architecture:** Keep the current result computation intact. Reuse the existing “pick latest run between single and matrix” logic, capture the selected workflow names as summary fields, then update the test fixture and README wording to expose those fields clearly.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing tests for remote source fields

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Assert that JSON output includes:
- `remote_happy_latest_source`
- `remote_recovery_latest_source`

and that text output includes:
- `Remote Happy Latest Source`
- `Remote Recovery Latest Source`

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper does not yet emit those fields.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- reuse the existing chosen workflow names
- store them in the summary
- print them in text mode

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that the readiness summary now shows which remote workflow provided the latest aggregated remote result.

- [ ] **Step 2: Confirm wording**

Keep the wording scoped to the aggregated remote categories only.

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
- only readiness observability changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-source-fields-design.md docs/superpowers/plans/2026-03-29-smoke-ready-source-fields.md
git commit -m "feat: add smoke-ready source fields"
```
