# Smoke Ready Result URLs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add direct run URL fields to the `smoke-ready` summary block so operators can jump straight from summary data to the underlying workflow runs.

**Architecture:** Keep the current result/source computation intact. Add summary URL fields by reusing the existing `workflows.latest` URLs, with a small lookup helper for aggregated remote categories, then update tests and README wording accordingly.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing tests for summary URL fields

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Assert that JSON output includes:
- `real_endpoint_latest_url`
- `real_endpoint_matrix_latest_url`
- `legacy_live_latest_url`
- `legacy_live_matrix_latest_url`
- `remote_happy_latest_url`
- `remote_recovery_latest_url`

Also assert the text output includes the paired URL labels.

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the summary does not yet expose those URL fields.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- derive the six URL fields
- print them in the text summary

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that the readiness summary now includes direct run URLs for the latest evidence.

- [ ] **Step 2: Confirm wording**

Keep the wording focused on the summary surface, not on workflow changes.

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
- only readiness usability changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-result-urls-design.md docs/superpowers/plans/2026-03-29-smoke-ready-result-urls.md
git commit -m "feat: add smoke-ready result urls"
```
