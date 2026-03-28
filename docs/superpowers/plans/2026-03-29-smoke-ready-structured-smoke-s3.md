# Smoke Ready Structured Smoke S3 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `smoke-ready` consume structured `Smoke S3` and `Smoke S3 Matrix` summaries first, while keeping raw-text parsing as a fallback.

**Architecture:** Keep the helper interfaces unchanged. Extend the real-endpoint detail loader to capture row summaries, update result derivation to prefer structured summary fields, then prove backward compatibility with tests that cover both structured and fallback artifact shapes.

**Tech Stack:** Python helper script, pytest, GitHub CLI artifact downloads

---

### Task 1: Add failing tests for structured real-endpoint summaries

**Files:**
- Modify: `scripts/test_smoke_ready.py`

- [ ] **Step 1: Write the failing test**

Change the mocked real-endpoint details so:
- `Smoke S3` includes a structured `summary` object
- `Smoke S3 Matrix` rows use structured row statuses instead of raw `present`

Assert that:
- `real_endpoint_latest_result` still resolves correctly
- `real_endpoint_matrix_latest_result` still resolves correctly
- the helper no longer needs raw text in those structured cases

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: FAIL because the helper still relies on the old fallback assumptions.

- [ ] **Step 3: Write minimal implementation**

Update `scripts/smoke_ready.py` to:
- carry structured single-run summaries
- prefer row `status` for matrix rows when already normalized
- fall back to raw output parsing only when needed

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready.py -q`
Expected: PASS

### Task 2: Update the helper note

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note that `smoke-ready` now consumes the structured `Smoke S3` summary artifacts produced by the workflows.

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
git add scripts/smoke_ready.py scripts/test_smoke_ready.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-structured-smoke-s3-design.md docs/superpowers/plans/2026-03-29-smoke-ready-structured-smoke-s3.md
git commit -m "feat: use structured smoke s3 summaries in smoke-ready"
```
