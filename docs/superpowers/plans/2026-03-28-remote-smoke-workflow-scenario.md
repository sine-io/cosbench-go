# Remote Smoke Workflow Scenario Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend the existing manual `Remote Smoke Local` workflow so it can trigger the multistage remote smoke scenario through the current helper.

**Architecture:** Keep the existing helper and workflow as the single source of truth. Add a small workflow-contract test, thread a new `scenario` input through the existing run step, and update README examples without changing the default CI path.

**Tech Stack:** GitHub Actions YAML, lightweight Python test, existing `make smoke-remote-local` helper, README documentation

---

### Task 1: Add A Failing Workflow Contract Test

**Files:**
- Create: `scripts/test_remote_smoke_workflow.py`

- [ ] **Step 1: Add a lightweight workflow text-shape test**

Cover:
- the manual workflow declares a `scenario` input
- the run step exports `SMOKE_REMOTE_LOCAL_SCENARIO`
- the workflow still references the existing `backend` input

- [ ] **Step 2: Run the focused test to confirm the red state**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_workflow.py -q
```

Expected:
- failure because the workflow does not yet expose `scenario`

### Task 2: Thread Scenario Through The Manual Workflow

**Files:**
- Modify: `.github/workflows/remote-smoke-local.yml`

- [ ] **Step 1: Add `workflow_dispatch.inputs.scenario`**

Use:
- default `single`
- description that makes `multistage` discoverable

- [ ] **Step 2: Pass scenario into the existing smoke command**

Keep the same workflow entrypoint, but add:
- `SMOKE_REMOTE_LOCAL_SCENARIO='${{ inputs.scenario }}'`

- [ ] **Step 3: Re-run the focused workflow test**

Run:
```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m pytest scripts/test_remote_smoke_workflow.py -q
```

Expected:
- the workflow contract test passes

- [ ] **Step 4: Commit the workflow slice**

Run:
```bash
git add .github/workflows/remote-smoke-local.yml scripts/test_remote_smoke_workflow.py
git commit -m "feat: add remote smoke workflow scenario input"
```

### Task 3: Update Documentation And Verify Local Smoke Paths

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Document manual workflow scenario usage**

Add:
- `gh workflow run "Remote Smoke Local" --repo sine-io/cosbench-go -f scenario=multistage`
- one combined example for `backend=s3, scenario=multistage`

- [ ] **Step 2: Verify default remote smoke still passes locally**

Run:
```bash
timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success with `scenario=single`

- [ ] **Step 3: Verify multistage remote smoke passes locally**

Run:
```bash
SMOKE_REMOTE_LOCAL_SCENARIO=multistage timeout 240s make --no-print-directory smoke-remote-local
```

Expected:
- success with `scenario=multistage`

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
git add README.md docs/superpowers/specs/2026-03-28-remote-smoke-workflow-scenario-design.md docs/superpowers/plans/2026-03-28-remote-smoke-workflow-scenario.md
git commit -m "docs: record remote smoke workflow scenario"
```
