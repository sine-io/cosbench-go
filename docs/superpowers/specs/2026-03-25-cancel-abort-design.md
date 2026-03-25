# Cancel/Abort Path Design

## Goal

Add a usable local cancel path for running jobs in `cosbench-go` without expanding into remote controller/driver abort orchestration.

This slice should let an operator cancel a currently running local job, preserve partial state, and surface the result cleanly through the existing job detail, history, event, and snapshot paths.

## Scope Boundary

### In Scope

- a distinct `cancelled` job/stage state for user-initiated cancellation
- control-plane support for cancelling a running local job
- runtime propagation of cancellation through the execution path
- partial result preservation when cancellation happens after work has begun
- a minimal Web entrypoint to cancel a running job
- focused tests for control-plane and Web behavior

### Out of Scope

- remote driver abort commands
- storage-adapter level `abort()` support
- distributed mission cancellation semantics
- retry/rollback logic after cancel
- new CLI cancel commands

## Why This Slice

The repository already has:

- `created`, `running`, `succeeded`, `failed`, and restart-recovery `interrupted`
- preflight validation before start
- persistent events and partial stage state snapshots
- a job detail page with a start action

What it still lacks is the operator-facing ability to stop a running job deliberately. That gap is explicitly called out in the migration board as the missing `cancel / abort path`.

## Recommended Approach

Implement a **local-only cancellation loop** in the control plane and expose it through the current job detail UI.

This is the smallest coherent closure because:

- the system already runs jobs in-process
- cancellation can be modeled with `context.CancelFunc`
- the existing event and snapshot model can already record state transitions

Do not add storage-level abort behavior in this slice. That belongs to a later, broader driver-control effort.

## State Model

Add a distinct `cancelled` status.

### Why not reuse `interrupted`

`interrupted` already has a repository-specific meaning:

- the process restarted while a job was running

User-initiated cancellation is a different operator-visible event and should not be conflated with restart recovery.

### Target status behavior

- `created` → `running` via start
- `running` → `cancelled` via user cancel
- `running` → `interrupted` only via restart recovery
- `running` → `failed` only for actual failure
- `running` → `succeeded` only for completed execution

## Control Plane Design

The control plane should own cancellation.

### Required additions

- keep a per-job `context.CancelFunc` for currently running jobs
- add `CancelJob(jobID string) error`
- reject cancellation for jobs that are not currently `running`
- append a cancellation event when the request is accepted

### Run-loop behavior

`runJob` should treat `context.Canceled` differently from ordinary execution errors:

- if cancellation occurs while a stage/work is running, mark the active stage `cancelled`
- mark the job `cancelled`
- preserve already collected stage/work metrics and partial totals
- persist snapshots and events before returning

This keeps cancellation observable rather than collapsing it into a generic failure.

## Execution-Layer Behavior

The execution engine currently treats `context.Canceled` as a normal non-error completion. That prevents the control plane from distinguishing external cancellation from runtime-based stopping.

The target behavior is:

- runtime deadline expiry remains a normal completion path
- external context cancellation returns `context.Canceled`

This allows the control plane to classify cancellation accurately.

## Web Behavior

Extend the current job detail page with one minimal action:

- when job status is `running`, show `Cancel Job`

The handler should expose:

- `POST /jobs/{id}/cancel`

On success:

- redirect back to `/jobs/{id}`

On error:

- redirect back with the current error mechanism

No new pages are needed.

## Partial Results Policy

Cancellation should preserve what the system already knows.

That means:

- completed stages remain completed
- the current stage becomes `cancelled`
- already collected work summaries from completed work in the active stage should be kept if available
- job-level metrics should reflect whatever was actually collected before cancellation

The policy is "stop cleanly, keep evidence", not "rewind the job".

## Tests To Add

### Control-plane tests

Add tests proving:

- a running job can be cancelled
- cancelled jobs end in `cancelled`, not `failed` or `interrupted`
- cancellation writes events
- partial metrics/results remain readable after cancellation

### Web tests

Add tests proving:

- running jobs render a cancel action
- `POST /jobs/{id}/cancel` redirects correctly
- cancelled jobs remain viewable in job detail

## Success Criteria

This slice is complete when:

1. running jobs can be cancelled through the control plane
2. running jobs can be cancelled from the job detail page
3. cancelled jobs use a dedicated `cancelled` state
4. cancellation preserves partial evidence rather than erasing it
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
