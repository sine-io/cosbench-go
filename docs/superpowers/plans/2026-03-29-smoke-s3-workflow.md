# Smoke S3 Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a dedicated manual workflow that runs the existing real-endpoint `make smoke-s3` path using repository secrets and a few optional workflow inputs.

**Architecture:** Keep the existing smoke test implementation unchanged. Add one `workflow_dispatch` workflow that exports the required `COSBENCH_SMOKE_*` environment variables, runs `GO=go make smoke-s3 | tee smoke-s3-output.txt`, uploads the output file, and writes it to the job summary. Lock the workflow contract with a focused Python test.

**Tech Stack:** GitHub Actions YAML, Python workflow contract test, existing `make smoke-s3`

---

### Task 1: Add A Failing Workflow Contract Test

**Files:**
- Create: `scripts/test_smoke_s3_workflow.py`

- [ ] **Step 1: Add a workflow-contract test**

Cover:
- workflow file exists
- `workflow_dispatch` exists
- inputs include:
  - `backend`
  - `region`
  - `path_style`
  - `bucket_prefix`
- secrets are threaded through:
  - `COSBENCH_SMOKE_ENDPOINT`
  - `COSBENCH_SMOKE_ACCESS_KEY`
  - `COSBENCH_SMOKE_SECRET_KEY`
- run step includes `GO=go make smoke-s3`
- artifact upload and job summary steps exist

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_s3_workflow.py -q
```

Expected:
- failure because the workflow file does not exist yet

### Task 2: Add The Manual Smoke S3 Workflow

**Files:**
- Create: `.github/workflows/smoke-s3.yml`

- [ ] **Step 1: Add the workflow file**

Use:
- `name: Smoke S3`
- `workflow_dispatch`

- [ ] **Step 2: Add the workflow inputs and environment mapping**

Map:
- secrets to required `COSBENCH_SMOKE_*`
- inputs to optional `COSBENCH_SMOKE_*`

- [ ] **Step 3: Add the run, summary, and artifact steps**

Ensure:
- output is tee’d to `smoke-s3-output.txt`
- summary consumes that file
- artifact upload uses `smoke-s3-output`

- [ ] **Step 4: Re-run the focused contract test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_s3_workflow.py -q
```

Expected:
- the workflow test passes

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add .github/workflows/smoke-s3.yml scripts/test_smoke_s3_workflow.py
git commit -m "feat: add smoke s3 workflow"
```

### Task 3: Final Verification And Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add a README invocation example**

Document:
- `gh workflow run "Smoke S3" --repo sine-io/cosbench-go`

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
git add README.md docs/superpowers/specs/2026-03-29-smoke-s3-workflow-design.md docs/superpowers/plans/2026-03-29-smoke-s3-workflow.md
git commit -m "docs: record smoke s3 workflow"
```
