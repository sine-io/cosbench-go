# Remote Smoke Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a manual GitHub Actions workflow that runs the existing remote multi-process MinIO smoke helper, uploads its artifacts, and publishes the smoke summary in the GitHub job summary.

**Architecture:** Keep the current smoke helper as the source of truth. Add one `workflow_dispatch` workflow that executes `make --no-print-directory smoke-remote-local`, preserves `.artifacts/remote-smoke/` on both success and failure, and reuses `summary.md` when present instead of duplicating smoke logic inside workflow YAML.

**Tech Stack:** GitHub Actions workflow YAML, existing `Makefile`, existing `scripts/smoke_remote_local.py`, repository docs in `README.md`

---

### Task 1: Add Failing Workflow-Level Coverage

**Files:**
- Create: `.github/workflows/remote-smoke-local.yml`
- Modify: `README.md`

- [ ] **Step 1: Write the workflow skeleton with no-op placeholders**

Create the workflow file with:
- `workflow_dispatch`
- checkout
- Go setup
- placeholder remote smoke step
- placeholder artifact upload

Do not wire the real smoke command yet.

- [ ] **Step 2: Add a documentation note for the new workflow entrypoint**

Document that a dedicated remote smoke workflow exists and is intended for manual execution.

- [ ] **Step 3: Run a local YAML sanity check**

Run:
```bash
python3 - <<'PY'
from pathlib import Path
import yaml
path = Path(".github/workflows/remote-smoke-local.yml")
print(yaml.safe_load(path.read_text(encoding="utf-8"))["name"])
PY
```

Expected:
- the workflow parses as valid YAML

### Task 2: Wire The Workflow To The Existing Remote Smoke Helper

**Files:**
- Modify: `.github/workflows/remote-smoke-local.yml`

- [ ] **Step 1: Replace placeholder execution with the real smoke command**

Run:
- `make --no-print-directory smoke-remote-local`

Make the workflow fail if the command exits non-zero.

- [ ] **Step 2: Add artifact upload for `.artifacts/remote-smoke/`**

Upload on both success and failure.

- [ ] **Step 3: Add summary export from `.artifacts/remote-smoke/summary.md`**

If the file exists:
- append it to `$GITHUB_STEP_SUMMARY`

If the file does not exist:
- write a short fallback note explaining that the helper failed before summary emission

- [ ] **Step 4: Run another local YAML sanity check**

Run:
```bash
python3 - <<'PY'
from pathlib import Path
import yaml
path = Path(".github/workflows/remote-smoke-local.yml")
data = yaml.safe_load(path.read_text(encoding="utf-8"))
print(data["on"])
print([step.get("name") for step in data["jobs"]["remote_smoke"]["steps"]])
PY
```

Expected:
- valid YAML with visible smoke, artifact, and summary steps

- [ ] **Step 5: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-local.yml
git commit -m "feat: add remote smoke workflow"
```

### Task 3: Harden The Helper Only If The Workflow Needs It

**Files:**
- Modify: `scripts/smoke_remote_local.py`

- [ ] **Step 1: Review the helper for GitHub-runner portability issues**

Check only for concrete runner-facing problems such as:
- path assumptions
- permission handling for MinIO
- missing artifact emission on failure

- [ ] **Step 2: If needed, make the smallest portability fix**

Do not expand the smoke scope. Keep this slice focused on workflow compatibility.

- [ ] **Step 3: Run the local remote smoke helper**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- helper still passes locally and writes `.artifacts/remote-smoke/summary.json`
- helper still passes locally and writes `.artifacts/remote-smoke/summary.md`

- [ ] **Step 4: Commit the helper tweak only if a real fix was required**

Run:
```bash
git add scripts/smoke_remote_local.py
git commit -m "fix: harden remote smoke helper for workflow execution"
```

Skip this commit if no helper change was needed.

### Task 4: Final Verification And Documentation

**Files:**
- Modify: `README.md`
- Review only: `.github/workflows/remote-smoke-local.yml`
- Review only: `.artifacts/remote-smoke/summary.md`

- [ ] **Step 1: Document manual workflow execution**

Add a concrete `gh workflow run ...` example and explain what artifacts the workflow uploads.

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

- [ ] **Step 4: Run the local remote smoke helper one more time**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success
- `.artifacts/remote-smoke/summary.md` exists for workflow summary reuse

- [ ] **Step 5: Review final scope**

Run:
```bash
git diff -- .github/workflows README.md scripts/smoke_remote_local.py Makefile
```

Expected:
- the slice stays focused on manual remote smoke workflow automation

- [ ] **Step 6: Commit the docs slice**

Run:
```bash
git add README.md docs/superpowers/specs/2026-03-28-remote-smoke-workflow-design.md docs/superpowers/plans/2026-03-28-remote-smoke-workflow.md
git commit -m "docs: record remote smoke workflow"
```
