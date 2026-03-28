# Remote SIO Smoke Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend the existing remote multi-process MinIO smoke path so the same helper and workflow can validate both `s3` and `sio`.

**Architecture:** Keep one shared smoke helper and one shared manual workflow. Add a dedicated SIO fixture, parameterize the helper with `backend=s3|sio`, surface the backend in summary artifacts, and thread the backend choice through the GitHub workflow input without changing the existing artifact contract.

**Tech Stack:** Python smoke helper, GitHub Actions workflow YAML, existing `Makefile`, existing remote smoke fixture structure under `testdata/workloads/`

---

### Task 1: Add Failing SIO Fixture And Helper Tests

**Files:**
- Create: `testdata/workloads/remote-smoke-sio-two-driver.xml`
- Modify: `scripts/test_smoke_remote_local.py`

- [ ] **Step 1: Add the minimal SIO remote smoke fixture**

Mirror the existing S3 remote smoke fixture shape, but use:
- `storage type="sio"`
- MinIO-compatible config
- `workers="2"`
- small `totalOps`
- write-only workload shape

- [ ] **Step 2: Add failing helper-level tests for backend selection**

Cover:
- `backend=sio` selects the SIO fixture
- summaries include `backend`
- invalid or unknown backend is rejected cleanly

- [ ] **Step 3: Run focused Python tests to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- failures because the helper only knows the S3 remote smoke path today

### Task 2: Parameterize The Local Remote Smoke Helper

**Files:**
- Modify: `scripts/smoke_remote_local.py`
- Modify: `Makefile`

- [ ] **Step 1: Add backend selection to the helper**

Support:
- `SMOKE_REMOTE_LOCAL_BACKEND=s3`
- `SMOKE_REMOTE_LOCAL_BACKEND=sio`

Default remains:
- `s3`

- [ ] **Step 2: Select the fixture and backend-specific settings from that parameter**

At minimum:
- fixture path
- backend name in summaries

- [ ] **Step 3: Thread the backend into the make target**

`smoke-remote-local` should preserve the current default while allowing environment override.

- [ ] **Step 4: Re-run focused helper tests**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_smoke_remote_local.py -q
```

Expected:
- helper-level backend-selection tests pass

- [ ] **Step 5: Commit the helper slice**

Run:
```bash
git add testdata/workloads/remote-smoke-sio-two-driver.xml scripts/test_smoke_remote_local.py scripts/smoke_remote_local.py Makefile
git commit -m "feat: parameterize remote smoke backend"
```

### Task 3: Verify Local S3 And SIO Remote Smoke

**Files:**
- Review only: `.artifacts/remote-smoke/summary.json`
- Review only: `.artifacts/remote-smoke/summary.md`

- [ ] **Step 1: Re-run S3 remote smoke**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- `overall=pass`
- `backend=s3`

- [ ] **Step 2: Run SIO remote smoke**

Run:
```bash
SMOKE_REMOTE_LOCAL_BACKEND=sio timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- `overall=pass`
- `backend=sio`

- [ ] **Step 3: Verify summary artifacts include backend**

Check:
- `.artifacts/remote-smoke/summary.json`
- `.artifacts/remote-smoke/summary.md`

Expected:
- explicit backend recorded in both files

### Task 4: Add Manual Workflow Backend Input

**Files:**
- Modify: `.github/workflows/remote-smoke-local.yml`

- [ ] **Step 1: Add `workflow_dispatch.inputs.backend`**

Use explicit options or a clear default:
- `s3`
- `sio`

- [ ] **Step 2: Pass the backend through to the helper**

The workflow should invoke the same make/helper path rather than duplicating backend logic in YAML.

- [ ] **Step 3: Run a local YAML sanity check**

Run:
```bash
python3 - <<'PY'
from pathlib import Path
import yaml
path = Path(".github/workflows/remote-smoke-local.yml")
data = yaml.safe_load(path.read_text(encoding="utf-8"))
print(data["on"])
print(data["jobs"]["remote_smoke"]["steps"][2]["run"])
PY
```

Expected:
- workflow parses
- backend input and smoke command are visible

- [ ] **Step 4: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-local.yml
git commit -m "feat: add remote smoke backend input"
```

### Task 5: Final Verification And Documentation

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Document S3/SIO remote smoke parity**

Describe:
- local helper backend selection
- manual workflow backend input
- artifact output

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

- [ ] **Step 4: Re-run final local SIO remote smoke**

Run:
```bash
SMOKE_REMOTE_LOCAL_BACKEND=sio timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success with explicit `backend=sio` in summary output

- [ ] **Step 5: Review final scope**

Run:
```bash
git diff -- testdata/workloads scripts Makefile .github/workflows README.md docs/migration-gap-analysis.md
```

Expected:
- the slice stays focused on S3/SIO parity for the existing remote smoke path

- [ ] **Step 6: Commit the docs slice**

Run:
```bash
git add README.md docs/migration-gap-analysis.md docs/superpowers/specs/2026-03-28-remote-sio-smoke-parity-design.md docs/superpowers/plans/2026-03-28-remote-sio-smoke-parity.md
git commit -m "docs: record remote sio smoke parity"
```
