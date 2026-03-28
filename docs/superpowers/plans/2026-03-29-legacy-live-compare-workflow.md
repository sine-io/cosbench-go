# Legacy Live Compare Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a manual workflow that renders one legacy sample into a runnable workload and executes it against a real endpoint with repository secrets.

**Architecture:** Keep the legacy sample files immutable and add a small render helper that substitutes smoke env values into a temporary XML. Use the existing CLI with `-json -quiet -summary-file` to execute the rendered workload. Back the workflow with focused Python tests for both the render helper and the workflow contract.

**Tech Stack:** Python render helper, GitHub Actions YAML, existing CLI entrypoint, Python workflow tests

---

### Task 1: Add Failing Tests For Rendering And Workflow Contract

**Files:**
- Create: `scripts/test_render_legacy_live_compare_workload.py`
- Create: `scripts/test_legacy_live_compare_workflow.py`

- [ ] **Step 1: Add render-helper tests**

Cover:
- rendered S3 sample no longer contains placeholder tokens
- rendered SIO sample no longer contains placeholder tokens
- storage type remains unchanged

- [ ] **Step 2: Add workflow-contract test**

Cover:
- workflow file exists
- `workflow_dispatch` inputs include `fixture`, `backend`, `region`, `path_style`
- workflow runs the render helper
- workflow runs `go run ./cmd/cosbench-go ... -json -quiet -summary-file`
- workflow uploads rendered XML, summary JSON, and log output

- [ ] **Step 3: Run the focused tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_render_legacy_live_compare_workload.py scripts/test_legacy_live_compare_workflow.py -q
```

Expected:
- failure because neither the render helper nor the workflow exists yet

### Task 2: Add The Render Helper And Workflow

**Files:**
- Create: `scripts/render_legacy_live_compare_workload.py`
- Create: `.github/workflows/legacy-live-compare.yml`

- [ ] **Step 1: Add the render helper**

Implement:
- input legacy fixture path
- output rendered workload path
- placeholder substitution from provided values

- [ ] **Step 2: Add the workflow file**

Use:
- `workflow_dispatch`
- required secrets
- inputs for fixture/backend/region/path_style

- [ ] **Step 3: Run the CLI through the rendered workload**

Persist:
- rendered XML
- summary JSON
- run log

- [ ] **Step 4: Re-run the focused tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_render_legacy_live_compare_workload.py scripts/test_legacy_live_compare_workflow.py -q
```

Expected:
- both tests pass

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add scripts/render_legacy_live_compare_workload.py scripts/test_render_legacy_live_compare_workload.py scripts/test_legacy_live_compare_workflow.py .github/workflows/legacy-live-compare.yml
git commit -m "feat: add legacy live compare workflow"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add a README invocation example**

Document:
- `gh workflow run "Legacy Live Compare" --repo sine-io/cosbench-go`

- [ ] **Step 2: Run the full Go test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all Go packages pass

- [ ] **Step 3: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 4: Commit the docs slice**

Run:
```bash
git add README.md docs/superpowers/specs/2026-03-29-legacy-live-compare-workflow-design.md docs/superpowers/plans/2026-03-29-legacy-live-compare-workflow.md
git commit -m "docs: record legacy live compare workflow"
```
