# Smoke Ready Artifact Fields Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add artifact-name fields to the `smoke-ready` summary block so operators can jump from summary data to the correct downloadable artifact name.

**Architecture:** Keep current result/source/url/timestamp logic intact. Add six artifact fields by mapping known workflow artifact names, using the already-selected aggregated source workflow names for the remote categories, then update tests and README wording.

**Tech Stack:** Python helper script, pytest, Markdown

---

### Task 1: Add failing tests for summary artifact fields

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Assert that JSON output includes:
- `real_endpoint_latest_artifact`
- `real_endpoint_matrix_latest_artifact`
- `legacy_live_latest_artifact`
- `legacy_live_matrix_latest_artifact`
- `remote_happy_latest_artifact`
- `remote_recovery_latest_artifact`

Also assert that text output includes the paired artifact labels.

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the summary does not yet expose those artifact fields.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to derive and print the six artifact fields.

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that the readiness summary now also names the latest artifact to download.

- [ ] **Step 2: Confirm wording**

Keep the wording focused on summary usability and artifact discoverability.

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
- only summary usability changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-artifact-fields-design.md docs/superpowers/plans/2026-03-29-smoke-ready-artifact-fields.md
git commit -m "feat: add smoke-ready artifact fields"
```
