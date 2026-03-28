# Remote SIO Multistage Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the missing `sio + multistage` remote smoke combination to the existing helper and workflow surface.

**Architecture:** Keep the current `backend + scenario` helper interface and the existing manual workflow unchanged. Add one `sio` multistage fixture, extend fixture selection logic, and verify the new path locally against MinIO without touching controller or driver code.

**Tech Stack:** XML workload fixtures, Python smoke helper and tests, existing README workflow examples, local MinIO via `make smoke-remote-local`

---

### Task 1: Add The SIO Multistage Fixture And Failing Helper Tests

**Files:**
- Create: `testdata/workloads/remote-smoke-sio-multistage-two-driver.xml`
- Modify: `scripts/test_smoke_remote_local.py`

- [ ] **Step 1: Add the minimal SIO multistage fixture**

Mirror the existing S3 multistage fixture shape with:
- `storage type="sio"`
- two stages
- one write-only work per stage
- `workers="2"` per stage
- distinct object ranges across stages

- [ ] **Step 2: Add failing helper tests for the new combination**

Cover:
- `fixture_for_selection("sio", "multistage")` selects the new fixture
- the new fixture contains two stages and two workers per stage
- existing `s3 + multistage` and `sio + single` selection tests still hold

- [ ] **Step 3: Run focused helper tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- failure because `sio + multistage` is still rejected today

### Task 2: Extend Helper Selection To Support SIO Multistage

**Files:**
- Modify: `scripts/smoke_remote_local.py`

- [ ] **Step 1: Add a dedicated SIO multistage fixture constant**

Keep naming aligned with the existing fixture constants.

- [ ] **Step 2: Extend `fixture_for_selection(backend, scenario)`**

Support:
- `s3 + single`
- `s3 + multistage`
- `sio + single`
- `sio + multistage`

Do not add new helper inputs or new summary fields.

- [ ] **Step 3: Re-run focused helper tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- helper selection tests pass

- [ ] **Step 4: Commit the helper slice**

Run:
```bash
git add testdata/workloads/remote-smoke-sio-multistage-two-driver.xml scripts/test_smoke_remote_local.py scripts/smoke_remote_local.py
git commit -m "feat: add remote sio multistage smoke parity"
```

### Task 3: Verify Local MinIO Smoke And Document Workflow Usage

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the new workflow invocation example**

Add:
- `gh workflow run "Remote Smoke Local" --repo sine-io/cosbench-go -f backend=sio -f scenario=multistage`

- [ ] **Step 2: Verify the existing single-stage SIO path still passes**

Run:
```bash
SMOKE_REMOTE_LOCAL_BACKEND=sio timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success with `scenario=single`

- [ ] **Step 3: Verify the new SIO multistage path**

Run:
```bash
SMOKE_REMOTE_LOCAL_BACKEND=sio SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success with:
  - `job_status=succeeded`
  - `stages_seen=2`
  - `stage_barrier=pass`

- [ ] **Step 4: Run the full Go test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all Go packages pass

- [ ] **Step 5: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 6: Commit the docs slice**

Run:
```bash
git add README.md docs/superpowers/specs/2026-03-28-remote-sio-multistage-parity-design.md docs/superpowers/plans/2026-03-28-remote-sio-multistage-parity.md
git commit -m "docs: record remote sio multistage parity"
```
