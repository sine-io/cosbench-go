# Cancel/Abort Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a usable local cancel path for running jobs, including control-plane cancellation, proper cancelled state handling, and a minimal Web cancel action.

**Architecture:** Implement cancellation in the control plane with per-job `context.CancelFunc`, expose a dedicated `cancelled` status, and teach the execution layer to distinguish external cancellation from normal runtime completion. Surface the capability through the existing job detail page and preserve partial metrics/events instead of discarding them.

**Tech Stack:** Go 1.26, package-local tests in `internal/controlplane`, `internal/web`, and `internal/domain/execution`, existing templates under `web/templates`

---

### Task 1: Add Failing Control-Plane and Web Tests

**Files:**
- Modify: `internal/controlplane/manager_test.go`
- Modify: `internal/web/handler_test.go`

- [ ] **Step 1: Write failing cancel tests**

Add tests proving:
- a running job can be cancelled
- cancelled jobs end in `cancelled`
- cancellation writes an event
- partial results remain readable
- the job detail page shows a cancel action for a running job
- `POST /jobs/{id}/cancel` redirects and the job ends in `cancelled`

- [ ] **Step 2: Run the focused tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- failures because there is no cancelled state, no `CancelJob`, and no Web route/button

### Task 2: Implement Cancelled State and Control-Plane Cancellation

**Files:**
- Modify: `internal/domain/job.go`
- Modify: `internal/controlplane/manager.go`

- [ ] **Step 1: Add the cancelled status and runtime cancel bookkeeping**

Implement:
- `JobStatusCancelled`
- per-job cancel function storage in `Manager`
- lifecycle cleanup for completed/cancelled/failed jobs

- [ ] **Step 2: Add `CancelJob(jobID)`**

Behavior:
- only works for currently running jobs
- records a cancellation-requested event
- triggers the stored cancel function

- [ ] **Step 3: Teach `runJob` to classify cancellation separately**

When work exits due to `context.Canceled`:
- mark the active stage `cancelled`
- mark the job `cancelled`
- preserve partial metrics/work summaries/events
- persist snapshots

- [ ] **Step 4: Re-run focused control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- control-plane tests pass with cancellation semantics covered

### Task 3: Fix Execution-Layer Cancellation Semantics

**Files:**
- Modify: `internal/domain/execution/engine.go`
- Modify: `internal/executor/executor.go`
- Modify: `internal/domain/execution/engine_test.go`

- [ ] **Step 1: Add a failing execution test for external cancellation**

Write a test showing:
- external context cancellation returns `context.Canceled`
- timeout-based completion for runtime-controlled work still behaves as normal completion

- [ ] **Step 2: Run execution tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution ./internal/executor
```

Expected:
- failures because cancellation is currently swallowed as a normal stop

- [ ] **Step 3: Implement the minimal cancellation distinction**

Implement:
- engine returns `context.Canceled` for external cancellation
- partial sample summary is preserved in `StageExecutor.RunWork`

- [ ] **Step 4: Re-run execution tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution ./internal/executor
```

Expected:
- execution-layer packages pass with cancellation behavior covered

### Task 4: Add Minimal Web Cancel Flow

**Files:**
- Modify: `internal/web/handler.go`
- Modify: `web/templates/job_detail.html`

- [ ] **Step 1: Add the Web route and handler branch**

Implement:
- `POST /jobs/{id}/cancel`
- redirect back to the job page on success or error

- [ ] **Step 2: Add the running-job cancel button**

Show `Cancel Job` only when `.Job.Status == "running"`.

- [ ] **Step 3: Re-run focused Web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- Web package passes with cancel flow covered

### Task 5: Final Verification and Status Update

**Files:**
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Update board/checklist state**

Reflect that the local cancel/abort gap is closed for the single-process scope.

- [ ] **Step 2: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass

- [ ] **Step 3: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 4: Review scope**

Run:
```bash
git diff -- internal/domain/job.go internal/controlplane internal/domain/execution internal/executor internal/web web/templates/job_detail.html BOARD.md TODO.md
```

Expected:
- the slice stays focused on local cancellation only
