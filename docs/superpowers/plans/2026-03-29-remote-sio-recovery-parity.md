# Remote SIO Recovery Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the missing `sio + recovery` remote smoke combination without changing the existing recovery orchestration or workflow topology.

**Architecture:** Reuse the current `scenario=recovery` helper path and add only the missing `sio` fixture plus fixture-selection support. Validate the new path locally against MinIO and document how to trigger the same combination through the existing parameterized `Remote Smoke Local` workflow.

**Tech Stack:** XML workload fixtures, Python smoke helper/tests, local MinIO, existing README workflow examples

---

### Task 1: Add The SIO Recovery Fixture And Failing Helper Tests

**Files:**
- Create: `testdata/workloads/remote-smoke-sio-recovery-two-driver.xml`
- Modify: `scripts/test_smoke_remote_local.py`

- [ ] **Step 1: Add the minimal SIO recovery fixture**

Mirror the S3 recovery fixture shape with:
- `storage type="sio"`
- one stage
- one work
- `workers="2"`
- `delay` operation with `duration=45s`

- [ ] **Step 2: Add failing helper tests**

Cover:
- `fixture_for_selection("sio", "recovery")` selects the new fixture
- the new fixture contains the expected delay-based recovery shape

- [ ] **Step 3: Run focused helper tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- failure because `sio + recovery` is not supported yet

### Task 2: Extend Helper Selection For SIO Recovery

**Files:**
- Modify: `scripts/smoke_remote_local.py`

- [ ] **Step 1: Add a dedicated SIO recovery fixture constant**

- [ ] **Step 2: Extend `fixture_for_selection(backend, scenario)`**

Support:
- `s3 + recovery`
- `sio + recovery`

Do not change recovery orchestration or summary shape in this slice.

- [ ] **Step 3: Re-run the focused helper tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- helper tests pass

- [ ] **Step 4: Commit the helper slice**

Run:
```bash
git add testdata/workloads/remote-smoke-sio-recovery-two-driver.xml scripts/test_smoke_remote_local.py scripts/smoke_remote_local.py
git commit -m "feat: add remote sio recovery smoke parity"
```

### Task 3: Verify Local Recovery Paths And Document Usage

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add one README example for SIO recovery**

Document:
- local command for `sio + recovery`
- one `gh workflow run "Remote Smoke Local" ... -f backend=sio -f scenario=recovery` example

- [ ] **Step 2: Re-run the existing S3 recovery path**

Run:
```bash
SMOKE_REMOTE_LOCAL_SCENARIO=recovery timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success

- [ ] **Step 3: Run the new SIO recovery path**

Run:
```bash
SMOKE_REMOTE_LOCAL_BACKEND=sio SMOKE_REMOTE_LOCAL_SCENARIO=recovery timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success with recovery observed

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
git add README.md docs/superpowers/specs/2026-03-29-remote-sio-recovery-parity-design.md docs/superpowers/plans/2026-03-29-remote-sio-recovery-parity.md
git commit -m "docs: record remote sio recovery parity"
```
