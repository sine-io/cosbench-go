# Remote Recovery Smoke Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a local/manual remote smoke scenario that proves mission lease expiry and reassignment when one driver disappears after claiming work.

**Architecture:** Extend the existing `smoke_remote_local.py` helper with a new `recovery` scenario and one dedicated S3 recovery fixture. Keep the current controller/driver topology and artifact layout, but add orchestration that intentionally stops driver-1 after claim, waits for reassignment, and records recovery-specific evidence in the summary.

**Tech Stack:** Python smoke helper, XML workload fixture, existing controller snapshot files, local MinIO via `make smoke-remote-local`

---

### Task 1: Add The Recovery Fixture And Failing Helper Tests

**Files:**
- Create: `testdata/workloads/remote-smoke-s3-recovery-two-driver.xml`
- Modify: `scripts/test_smoke_remote_local.py`

- [ ] **Step 1: Add the minimal recovery fixture**

Use:
- `storage type="s3"`
- one stage
- one work
- `workers="2"`
- small `totalOps`
- write-only operation shape

- [ ] **Step 2: Add failing helper tests**

Cover:
- `scenario=recovery` selects the new fixture
- recovery summaries include:
  - `recovery_observed`
  - `reclaimed_units`

- [ ] **Step 3: Run focused helper tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- failure because `scenario=recovery` is not implemented yet

### Task 2: Implement Recovery Scenario In The Helper

**Files:**
- Modify: `scripts/smoke_remote_local.py`

- [ ] **Step 1: Add fixture selection for `scenario=recovery`**

Support:
- `backend=s3`
- `scenario=recovery`

- [ ] **Step 2: Add recovery orchestration**

Implement:
- wait for a claim by driver-1
- stop driver-1
- wait for reassignment to driver-2
- wait for final job completion

- [ ] **Step 3: Extend summary output with recovery fields**

Add:
- `recovery_observed`
- `reclaimed_units`

- [ ] **Step 4: Re-run the focused helper tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- helper selection and summary-shape tests pass

- [ ] **Step 5: Commit the helper slice**

Run:
```bash
git add testdata/workloads/remote-smoke-s3-recovery-two-driver.xml scripts/test_smoke_remote_local.py scripts/smoke_remote_local.py
git commit -m "feat: add remote recovery smoke"
```

### Task 3: Verify Local Recovery Smoke And Document Usage

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document the local recovery smoke command**

Add:
- `SMOKE_REMOTE_LOCAL_SCENARIO=recovery make --no-print-directory smoke-remote-local`

- [ ] **Step 2: Re-run existing happy-path smoke**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- default `single` scenario still passes

- [ ] **Step 3: Run recovery smoke**

Run:
```bash
SMOKE_REMOTE_LOCAL_SCENARIO=recovery timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success
- `recovery_observed=true`
- `reclaimed_units >= 1`

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
git add README.md docs/superpowers/specs/2026-03-29-remote-recovery-smoke-design.md docs/superpowers/plans/2026-03-29-remote-recovery-smoke.md
git commit -m "docs: record remote recovery smoke"
```
