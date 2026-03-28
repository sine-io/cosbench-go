# Remote Multi-Process MinIO Smoke Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a repeatable local smoke helper that validates the remote controller/driver protocol across one controller-only process, two driver-only processes, and one local MinIO instance.

**Architecture:** Keep the current runtime modes and remote APIs unchanged. Add one minimal remote-smoke workload fixture, a purpose-built orchestration helper, and stable artifact generation under `.artifacts/remote-smoke/`. Use the helper to prove cross-process registration, work-unit claims, execution, and aggregation on a real local S3-compatible endpoint.

**Tech Stack:** Go 1.26 runtime binaries, Python orchestration helper, existing `cmd/server`, existing controller and driver HTTP APIs, local MinIO

---

### Task 1: Add Failing Fixture And Helper-Level Smoke Tests

**Files:**
- Create: `testdata/workloads/remote-smoke-s3-two-driver.xml`
- Create: `scripts/test_smoke_remote_local.py`

- [ ] **Step 1: Add the minimal remote smoke fixture**

Create a fixture with:
- `storage type="s3"`
- one stage
- one work
- `workers="2"`
- small `totalOps`
- `write`-only workload shape

- [ ] **Step 2: Add focused helper tests**

Cover:
- artifact directory shape
- summary rendering shape
- fixed failure when a required process or check is missing

- [ ] **Step 3: Run the focused Python tests to verify the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- failures because the helper does not exist yet

### Task 2: Implement The Remote Smoke Orchestration Helper

**Files:**
- Create: `scripts/smoke_remote_local.py`
- Modify: `Makefile`

- [ ] **Step 1: Implement process orchestration**

The helper should:
- start MinIO
- start one `controller-only`
- start two `driver-only`
- generate ports, temp data dirs, and one shared token

- [ ] **Step 2: Implement readiness checks**

Wait for:
- MinIO socket/health
- controller HTTP availability
- both driver HTTP surfaces

- [ ] **Step 3: Add a make target for local invocation**

Suggested target:
```make
smoke-remote-local:
```

- [ ] **Step 4: Run focused helper tests again**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- helper-level tests now pass

- [ ] **Step 5: Commit the orchestration slice**

Run:
```bash
git add testdata/workloads/remote-smoke-s3-two-driver.xml scripts/test_smoke_remote_local.py scripts/smoke_remote_local.py Makefile
git commit -m "feat: add remote multi-process smoke helper"
```

### Task 3: Add Controller Submission And Verification Flow

**Files:**
- Modify: `scripts/smoke_remote_local.py`

- [ ] **Step 1: Implement workload submission/start flow**

The helper should:
- submit the fixture to the controller
- start the created job
- poll job status until terminal

- [ ] **Step 2: Implement remote smoke checks**

Required checks:
- two healthy registered drivers
- at least two claimed/completed work units or attempts
- both drivers participated
- job status is `succeeded`
- metrics are non-zero

- [ ] **Step 3: Implement controller/driver API fetches**

Fetch at least:
- controller matrix/detail/timeline
- driver self and mission views

- [ ] **Step 4: Re-run the helper in local smoke mode**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- the helper completes successfully against local MinIO and local controller/driver processes

### Task 4: Add Stable Artifact Outputs

**Files:**
- Modify: `scripts/smoke_remote_local.py`

- [ ] **Step 1: Write `summary.json` and `summary.md`**

Include at least:
- controller URL
- driver URLs
- job id
- job status
- drivers seen
- units claimed
- drivers participated
- operation count
- byte count
- checks
- overall result

- [ ] **Step 2: Persist logs under `.artifacts/remote-smoke/`**

Capture:
- controller log
- driver1 log
- driver2 log
- minio log

- [ ] **Step 3: Make failures still emit summaries**

Even when the helper exits non-zero, it should still write the summary files with failure context.

- [ ] **Step 4: Re-run local remote smoke**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success path writes complete artifacts

### Task 5: Final Verification And Docs Refresh

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Document the new smoke path**

Describe:
- `smoke-remote-local`
- what it validates
- where artifacts are written

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

- [ ] **Step 4: Run the final local remote smoke one more time**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- the end-to-end local remote smoke passes with final artifact output

- [ ] **Step 5: Review final scope**

Run:
```bash
git diff -- testdata/workloads scripts Makefile README.md docs/migration-gap-analysis.md
```

Expected:
- the slice stays focused on the multi-process MinIO remote smoke path

- [ ] **Step 6: Commit the docs slice**

Run:
```bash
git add README.md docs/migration-gap-analysis.md
git commit -m "docs: record remote multi-process smoke path"
```
