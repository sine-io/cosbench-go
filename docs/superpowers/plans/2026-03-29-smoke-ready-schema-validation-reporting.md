# Smoke Ready Schema Validation Reporting Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `Smoke Ready Validate` as a first-class signal in `smoke-ready` and `smoke-ready-json`.

**Architecture:** Extend the existing workflow discovery and summary pipeline with one additional contract-surface workflow. Reuse the current pattern of downloading workflow artifacts, derive one normalized result from `validation.json`, and expose that result through the same summary shape already used elsewhere.

**Tech Stack:** Python helper scripts, pytest, JSON Schema, Markdown

---

### Task 1: Add failing smoke-ready coverage

**Files:**
- Modify: `scripts/test_smoke_ready.py`
- Modify: `scripts/test_smoke_ready_schema.py`

- [ ] **Step 1: Write the failing tests**

Update mocked workflow inputs so they include `Smoke Ready Validate`, then assert:
- `workflows.present["Smoke Ready Validate"] == true`
- `workflows.latest["Smoke Ready Validate"]` is present
- summary includes schema-validation fields
- `schema_validation_latest_result == "validated"`

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py -q`
Expected: FAIL because `smoke_ready.py` does not expose the new workflow yet.

### Task 2: Implement schema-validation reporting

**Files:**
- Modify: `scripts/smoke_ready.py`
- Modify: `docs/smoke-ready.schema.json`

- [ ] **Step 1: Write minimal implementation**

Add:
- workflow constant and workflow list entry
- artifact mapping
- loader for `smoke-ready-validate-output/validation.json`
- result classification logic
- summary fields
- text output lines

- [ ] **Step 2: Run tests to verify they pass**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Update docs**

Add short notes that:
- `Smoke Ready Validate` is now part of the readiness surface
- `smoke-ready` can report the latest schema-validation evidence

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
- only readiness/reporting and contract-evidence visibility changed

- [ ] **Step 6: Commit**

```bash
git add scripts/smoke_ready.py scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py docs/smoke-ready.schema.json README.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-29-smoke-ready-schema-validation-reporting-design.md docs/superpowers/plans/2026-03-29-smoke-ready-schema-validation-reporting.md
git commit -m "feat: report smoke-ready schema validation state"
```
