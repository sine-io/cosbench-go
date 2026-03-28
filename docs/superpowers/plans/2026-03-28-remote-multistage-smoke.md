# Remote Multistage Smoke Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend the existing local remote smoke helper so it can prove multi-stage remote progression with one controller-only process, two driver-only processes, and local MinIO.

**Architecture:** Keep the current helper and artifact contract as the source of truth. Add a dedicated S3 multistage fixture, parameterize helper execution with `scenario=single|multistage`, add stage-aware checks and summary fields, and verify both scenarios locally against MinIO without changing the workflow surface in this slice.

**Tech Stack:** Python smoke helper, XML workload fixtures, existing Go controller and driver APIs, local MinIO via `make smoke-remote-local`

---

### Task 1: Add The Multistage Fixture And Failing Helper Tests

**Files:**
- Create: `testdata/workloads/remote-smoke-s3-multistage-two-driver.xml`
- Modify: `scripts/test_smoke_remote_local.py`

- [ ] **Step 1: Add the minimal multistage S3 fixture**

Create a workload with:
- `storage type="s3"`
- two `workstage` blocks
- one write-only `work` per stage
- `workers="2"` in each stage
- distinct object ranges between stages

- [ ] **Step 2: Add failing helper tests for scenario selection and summary shape**

Cover:
- `scenario=single` keeps the current single-stage fixture
- `scenario=multistage` selects the new multistage fixture
- summaries include `scenario`
- summaries can include stage metadata without breaking JSON serialization

- [ ] **Step 3: Run focused Python tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- failures because the helper only knows the current single-stage fixture shape today

### Task 2: Parameterize The Helper And Add Multistage Checks

**Files:**
- Modify: `scripts/smoke_remote_local.py`
- Modify: `Makefile`

- [ ] **Step 1: Add scenario selection to the helper**

Support:
- `SMOKE_REMOTE_LOCAL_SCENARIO=single`
- `SMOKE_REMOTE_LOCAL_SCENARIO=multistage`

Default remains:
- `single`

- [ ] **Step 2: Select the fixture from backend plus scenario**

Keep the current backend parameter intact, but route to the multistage fixture when:
- `backend=s3`
- `scenario=multistage`

- [ ] **Step 3: Extend the summary contract carefully**

Add:
- `scenario`
- `stage_names`
- `stages_seen`

Do not remove or rename existing summary keys.

- [ ] **Step 4: Add multistage-specific checks**

For `scenario=multistage`, validate:
- at least two succeeded stages
- mission snapshots span both stage names
- stage ordering barrier is preserved through timestamp comparison

- [ ] **Step 5: Re-run focused helper tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- helper tests pass with both single-stage and multistage fixture selection

- [ ] **Step 6: Commit the helper slice**

Run:
```bash
git add testdata/workloads/remote-smoke-s3-multistage-two-driver.xml scripts/test_smoke_remote_local.py scripts/smoke_remote_local.py Makefile
git commit -m "feat: add remote multistage smoke scenario"
```

### Task 3: Verify Local Single-Stage And Multistage Remote Smoke

**Files:**
- Review only: `.artifacts/remote-smoke/summary.json`
- Review only: `.artifacts/remote-smoke/summary.md`

- [ ] **Step 1: Re-run the current default remote smoke**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- `overall=pass`
- `scenario=single`

- [ ] **Step 2: Run the multistage remote smoke**

Run:
```bash
SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- `overall=pass`
- `scenario=multistage`
- at least two stage names recorded

- [ ] **Step 3: Verify summary artifacts include scenario and stage metadata**

Check:
- `.artifacts/remote-smoke/summary.json`
- `.artifacts/remote-smoke/summary.md`

Expected:
- explicit `scenario`
- multistage run includes stage names and stage count

### Task 4: Final Verification And Documentation

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Document multistage local remote smoke usage**

Describe:
- default single-stage behavior
- `SMOKE_REMOTE_LOCAL_SCENARIO=multistage`
- summary artifact additions

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

- [ ] **Step 4: Re-run final multistage smoke for fresh evidence**

Run:
```bash
SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- multistage remote smoke passes with fresh summary output

- [ ] **Step 5: Review final scope**

Run:
```bash
git diff -- testdata/workloads scripts Makefile README.md docs/migration-gap-analysis.md
```

Expected:
- the slice stays focused on multistage local remote smoke coverage

- [ ] **Step 6: Commit the docs slice**

Run:
```bash
git add README.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-28-remote-multistage-smoke-design.md docs/superpowers/plans/2026-03-28-remote-multistage-smoke.md
git commit -m "docs: record remote multistage smoke"
```
