# Smoke Ready Schema Validation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a direct helper and Make targets for validating current `smoke-ready-json` output against the published schema.

**Architecture:** Keep `scripts/smoke_ready.py` focused on producing payloads. Add one separate validator script that shells out to `scripts/smoke_ready.py --json`, validates that output against `docs/smoke-ready.schema.json`, and emits either text or JSON status. Expose that helper through new Make targets.

**Tech Stack:** Python helper scripts, pytest, `jsonschema`, Make, Markdown

---

### Task 1: Add failing validator coverage

**Files:**
- Create: `scripts/test_validate_smoke_ready_schema.py`

- [ ] **Step 1: Write the failing test**

Add a focused test that:
- runs `python3 scripts/validate_smoke_ready_schema.py --json`
- provides the same mocked smoke-ready environment used by existing tests
- expects `valid == true`
- expects `schema_path == "docs/smoke-ready.schema.json"`

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_validate_smoke_ready_schema.py -q`
Expected: FAIL because the validator helper does not exist yet.

### Task 2: Implement the validator helper

**Files:**
- Create: `scripts/validate_smoke_ready_schema.py`

- [ ] **Step 1: Write minimal implementation**

Implement a helper that:
- shells out to `scripts/smoke_ready.py --json`
- loads `docs/smoke-ready.schema.json`
- validates the payload with `jsonschema`
- prints text by default
- prints machine-readable JSON with `--json`

- [ ] **Step 2: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_validate_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Expose Make targets and docs

**Files:**
- Modify: `Makefile`
- Modify: `README.md`

- [ ] **Step 1: Add Make targets**

Add:
- `smoke-ready-validate`
- `smoke-ready-validate-json`

- [ ] **Step 2: Update README**

Add one short note describing the new validation entrypoints.

### Task 4: Verify and commit

**Files:**
- Verify only

- [ ] **Step 1: Run targeted tests**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py -q`
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
- only machine-readable contract tooling changed

- [ ] **Step 5: Commit**

```bash
git add scripts/validate_smoke_ready_schema.py scripts/test_validate_smoke_ready_schema.py Makefile README.md docs/superpowers/specs/2026-03-29-smoke-ready-schema-validation-design.md docs/superpowers/plans/2026-03-29-smoke-ready-schema-validation.md
git commit -m "feat: add smoke-ready schema validation helper"
```
