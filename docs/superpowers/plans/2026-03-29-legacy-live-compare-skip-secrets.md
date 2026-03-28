# Legacy Live Compare Skip-Secrets Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the manual `Legacy Live Compare` workflow skip cleanly when live endpoint secrets are absent, while preserving the existing live run path when they are present.

**Architecture:** Keep the change at the workflow layer. Add a preflight step that materializes a stable skipped artifact and gates the render/run steps behind a workflow output. Update the lightweight workflow tests and the user-facing docs to describe the new behavior.

**Tech Stack:** GitHub Actions YAML, Python pytest workflow-shape tests, Markdown docs

---

### Task 1: Lock the missing-secrets behavior in tests

**Files:**
- Modify: `.github/workflows/legacy-live-compare.yml`
- Test: `scripts/test_legacy_live_compare_workflow.py`

- [ ] **Step 1: Write the failing test**

Extend `scripts/test_legacy_live_compare_workflow.py` to require:
- a preflight step name for credential checking
- gated render/run steps using `if:`
- a stable skipped summary artifact path such as `.artifacts/legacy-live-compare/summary.json`

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_legacy_live_compare_workflow.py -q`
Expected: FAIL because the workflow currently lacks skip preflight logic.

- [ ] **Step 3: Write minimal workflow implementation**

Update `.github/workflows/legacy-live-compare.yml` to:
- preflight required secrets
- write skipped artifact files when absent
- skip render/run when secrets are missing

- [ ] **Step 4: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_legacy_live_compare_workflow.py -q`
Expected: PASS

### Task 2: Document the new skip semantics

**Files:**
- Modify: `README.md`
- Modify: `docs/legacy-live-run-checklist.md`

- [ ] **Step 1: Update README**

Add one short note near the `Legacy Live Compare` trigger example that the workflow will mark itself skipped when required live secrets are not configured.

- [ ] **Step 2: Update legacy live checklist**

Refresh the checklist language so the remote workflow path distinguishes:
- `Smoke S3` style skip-ready behavior when secrets are absent
- `Legacy Live Compare` as an opt-in workflow that now skips cleanly rather than failing on empty config

- [ ] **Step 3: Verify docs are consistent**

Read both updated sections and confirm they use the same secret names and skip terminology.

### Task 3: Verify end to end

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_legacy_live_compare_workflow.py -q`
Expected: PASS

- [ ] **Step 2: Run repo tests**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 3: Run full build**

Run: `go build ./...`
Expected: PASS

- [ ] **Step 4: Compare against overall project goal**

Confirm the change stays within the project goal:
- unified Go service remains unchanged
- no non-S3 backend expansion
- only workflow ergonomics for real-endpoint legacy validation changed

- [ ] **Step 5: Commit**

```bash
git add .github/workflows/legacy-live-compare.yml scripts/test_legacy_live_compare_workflow.py README.md docs/legacy-live-run-checklist.md docs/superpowers/specs/2026-03-29-legacy-live-compare-skip-secrets-design.md docs/superpowers/plans/2026-03-29-legacy-live-compare-skip-secrets.md
git commit -m "ci: skip legacy live compare without secrets"
```
