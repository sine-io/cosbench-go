# Smoke Ready JSON Schema Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Publish a repository-owned JSON Schema for `smoke-ready-json` and validate helper output against it.

**Architecture:** Keep the helper output unchanged. Add one committed schema document under `docs/`, extend the existing schema contract test to load that file and validate the mocked helper payload with `jsonschema`, and add one short README pointer for machine consumers.

**Tech Stack:** Python helper script, pytest, `jsonschema`, Markdown, JSON Schema Draft 2020-12

---

### Task 1: Add failing schema-file contract coverage

**Files:**
- Modify: `scripts/test_smoke_ready_schema.py`

- [ ] **Step 1: Write the failing test change**

Update the contract test so it:
- loads `docs/smoke-ready.schema.json`
- validates the helper payload with `jsonschema`
- still asserts `schema_version == 1`

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready_schema.py -q`
Expected: FAIL because `docs/smoke-ready.schema.json` does not exist yet.

### Task 2: Add the schema document

**Files:**
- Create: `docs/smoke-ready.schema.json`

- [ ] **Step 1: Write minimal schema**

Define:
- top-level required keys
- object shapes for `local_env`, `repo_secrets`, `workflows`, and `summary`
- required summary keys for the current latest evidence blocks

- [ ] **Step 2: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready_schema.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note pointing machine consumers to `docs/smoke-ready.schema.json`.

- [ ] **Step 2: Run targeted tests**

Run: `python3 -m pytest scripts/test_smoke_ready.py scripts/test_smoke_ready_schema.py -q`
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
- only machine-readable smoke-ready contract documentation changed

- [ ] **Step 6: Commit**

```bash
git add docs/smoke-ready.schema.json scripts/test_smoke_ready_schema.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-json-schema-design.md docs/superpowers/plans/2026-03-29-smoke-ready-json-schema.md
git commit -m "feat: add smoke-ready json schema"
```
