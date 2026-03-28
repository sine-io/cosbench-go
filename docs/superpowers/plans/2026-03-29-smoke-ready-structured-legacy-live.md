# Smoke Ready Structured Legacy Live Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `smoke-ready` consume structured `Legacy Live Compare` and `Legacy Live Compare Matrix` summaries first, while keeping workflow-step inspection as a fallback.

**Architecture:** Keep the helper interface unchanged. Extend the legacy detail loader to capture normalized summary artifacts, update result derivation to prefer those structured fields, then prove backward compatibility with tests that cover both structured and fallback artifact shapes.

**Tech Stack:** Python helper script, pytest, GitHub CLI artifact downloads

---

### Task 1: Add failing tests for structured legacy live summaries

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Change the mocked legacy details so:
- `Legacy Live Compare` includes a structured `result` object
- `Legacy Live Compare Matrix` uses aggregate row statuses instead of only job-step details

Assert that:
- `legacy_live_latest_result` still resolves correctly
- `legacy_live_matrix_latest_result` still resolves correctly
- the helper no longer needs step details in those structured cases

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper still relies on workflow-step inspection first.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- download and load structured legacy summary artifacts
- prefer artifact result fields for legacy live result derivation
- fall back to the existing step-based path only when needed

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that `smoke-ready` now consumes the structured legacy live summaries produced by the workflows.

- [ ] **Step 2: Confirm wording**

Keep the wording focused on evidence consumption, not workflow changes.

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
- only readiness evidence consumption changed

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-structured-legacy-live-design.md docs/superpowers/plans/2026-03-29-smoke-ready-structured-legacy-live.md
git commit -m "feat: use structured legacy live summaries in smoke-ready"
```
