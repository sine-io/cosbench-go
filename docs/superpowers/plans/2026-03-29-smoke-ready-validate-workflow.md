# Smoke Ready Validate Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a manual GitHub workflow that runs `smoke-ready` schema validation and publishes its result as an artifact and job summary.

**Architecture:** Keep the existing helper and schema validator unchanged. Add one dedicated workflow that captures both the raw `smoke-ready-json` payload and the validator output into `.artifacts/smoke-ready-validate/`, uploads that directory, and writes the validation JSON into the GitHub job summary.

**Tech Stack:** GitHub Actions YAML, Python helper scripts, Make, Markdown

---

### Task 1: Add failing workflow contract coverage

**Files:**
- Create: `scripts/test_smoke_ready_validate_workflow.py`

- [ ] **Step 1: Write the failing test**

Assert that `.github/workflows/smoke-ready-validate.yml`:
- exists
- uses `workflow_dispatch`
- sets `GH_TOKEN`
- runs `make --no-print-directory smoke-ready-json`
- runs `make --no-print-directory smoke-ready-validate-json`
- uploads `.artifacts/smoke-ready-validate`
- writes `validation.json` to `$GITHUB_STEP_SUMMARY`

- [ ] **Step 2: Run test to verify it fails**

Run: `python3 -m pytest scripts/test_smoke_ready_validate_workflow.py -q`
Expected: FAIL because the workflow file does not exist yet.

### Task 2: Implement the workflow

**Files:**
- Create: `.github/workflows/smoke-ready-validate.yml`

- [ ] **Step 1: Write minimal workflow**

Add:
- `workflow_dispatch`
- checkout
- `GH_TOKEN: ${{ github.token }}`
- shell steps to produce `.artifacts/smoke-ready-validate/smoke-ready.json`
- shell steps to produce `.artifacts/smoke-ready-validate/validation.json`
- artifact upload
- job summary write step

- [ ] **Step 2: Run test to verify it passes**

Run: `python3 -m pytest scripts/test_smoke_ready_validate_workflow.py -q`
Expected: PASS

### Task 3: Document and verify

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README**

Add one short note describing how to trigger `Smoke Ready Validate` with `gh workflow run`.

- [ ] **Step 2: Run targeted tests**

Run: `python3 -m pytest scripts/test_validate_smoke_ready_schema.py scripts/test_smoke_ready_validate_workflow.py -q`
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
- only contract tooling and remote workflow ergonomics changed

- [ ] **Step 6: Commit**

```bash
git add .github/workflows/smoke-ready-validate.yml scripts/test_smoke_ready_validate_workflow.py README.md docs/superpowers/specs/2026-03-29-smoke-ready-validate-workflow-design.md docs/superpowers/plans/2026-03-29-smoke-ready-validate-workflow.md
git commit -m "feat: add smoke-ready validate workflow"
```
