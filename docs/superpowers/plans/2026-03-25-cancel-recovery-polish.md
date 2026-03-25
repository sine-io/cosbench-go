# Cancel Recovery Polish Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Distinguish restart recovery of plain running jobs from jobs that had already entered user-requested cancellation.

**Architecture:** Introduce a transient `cancelling` state that is persisted immediately when a cancellation request is accepted. Keep `running -> interrupted` for unexpected restart recovery, but recover persisted `cancelling` jobs as `cancelled`. Expose `cancelling` in the Web UI and stop offering a second cancel action once that state is set.

**Tech Stack:** Go 1.26, package-local tests in `internal/controlplane` and `internal/web`, existing snapshot persistence

---

### Task 1: Add Failing Control-Plane and Web Tests

**Files:**
- Modify: `internal/controlplane/manager_test.go`
- Modify: `internal/web/handler_test.go`

- [ ] **Step 1: Write failing tests**

Add tests proving:
- `CancelJob()` moves a running job into `cancelling` before final cancellation
- a persisted `cancelling` job recovers as `cancelled`
- a persisted plain `running` job still recovers as `interrupted`
- a `cancelling` job detail page does not show `Cancel Job`

- [ ] **Step 2: Run focused tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- failures because there is no `cancelling` state or recovery path yet

### Task 2: Implement the `cancelling` State

**Files:**
- Modify: `internal/domain/job.go`
- Modify: `internal/controlplane/manager.go`

- [ ] **Step 1: Add `JobStatusCancelling`**

Introduce a distinct persisted status for jobs/stages whose cancel request has been accepted but whose goroutine has not yet finished.

- [ ] **Step 2: Persist `cancelling` immediately in `CancelJob()`**

Behavior:
- running job becomes `cancelling`
- active stage becomes `cancelling` when applicable
- cancellation-requested event is persisted
- stored cancel function is invoked

- [ ] **Step 3: Teach restart recovery to distinguish `running` vs `cancelling`**

Behavior:
- persisted `running` → `interrupted`
- persisted `cancelling` → `cancelled` with an explanatory recovery note

- [ ] **Step 4: Re-run control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- control-plane tests pass with polished recovery semantics

### Task 3: Update Web Status Presentation

**Files:**
- Modify: `internal/web/handler.go`
- Modify: `web/templates/job_detail.html`
- Modify: `web/static/app.css`

- [ ] **Step 1: Add `cancelling` presentation**

Behavior:
- `cancelling` status renders visibly
- `Cancel Job` button is hidden once job status is `cancelling`

- [ ] **Step 2: Re-run Web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- Web package passes with the new status handling

### Task 4: Final Verification and Status Update

**Files:**
- Modify: `BOARD.md`
- Modify: `TODO.md`
- Review only: `docs/storage-layout.md`

- [ ] **Step 1: Update board/checklist state**

Reflect that restart/recovery polish for cancelled jobs is no longer open.

- [ ] **Step 2: Review storage-layout wording**

If needed, clarify that:
- `running` jobs recover as `interrupted`
- `cancelling` jobs recover as `cancelled`

- [ ] **Step 3: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass

- [ ] **Step 4: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 5: Review scope**

Run:
```bash
git diff -- internal/domain/job.go internal/controlplane internal/web web/templates/job_detail.html web/static/app.css BOARD.md TODO.md docs/storage-layout.md
```

Expected:
- the slice stays limited to cancellation recovery polish
