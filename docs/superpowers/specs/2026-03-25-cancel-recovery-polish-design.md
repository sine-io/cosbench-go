# Cancel Recovery Polish Design

## Goal

Polish restart recovery for user-cancelled jobs so the system distinguishes:

- jobs that were merely running when the process died
- jobs that had already received a user cancellation request before the process died

This slice should eliminate the misleading case where a user-initiated cancellation is later shown as a generic `interrupted` recovery.

## Problem

Current behavior:

- `CancelJob()` records a cancellation request and triggers the in-memory cancel function
- until the running goroutine exits, the persisted job status remains `running`
- on restart, `loadSnapshots()` rewrites any persisted `running` job to `interrupted`

That means a cancelled-in-flight job can be misclassified as `interrupted` after restart.

## Scope

### In Scope

- add an intermediate `cancelling` state
- persist `cancelling` immediately when a user requests cancellation
- recover `cancelling` jobs as `cancelled` on restart
- keep ordinary `running` jobs recovering as `interrupted`
- make the Web UI show `cancelling` distinctly
- add focused control-plane and Web tests

### Out of Scope

- remote worker cancellation semantics
- storage-level abort support
- cancel request retries or backoff
- CLI cancel commands

## Recommended Approach

Introduce an explicit `cancelling` state.

This is the smallest honest representation because:

- it records that the user asked for cancellation
- it avoids claiming the job is already fully cancelled before the goroutine exits
- it gives restart recovery enough information to distinguish “cancel requested” from “unexpected interruption”

## State Model

Target lifecycle:

- `created` → `running`
- `running` → `cancelling` when the user requests cancellation
- `cancelling` → `cancelled` when the run loop exits on `context.Canceled`
- `running` → `interrupted` only when the process dies unexpectedly
- `cancelling` → `cancelled` on restart recovery, because the cancellation request is already authoritative

## Control Plane Design

### Cancel request path

When `CancelJob(jobID)` is called:

- verify the job is currently `running`
- mark the job `cancelling`
- mark the active stage `cancelling` if one is active
- append a cancellation-requested event
- persist immediately
- invoke the stored cancel function

### Restart recovery

In `loadSnapshots()`:

- persisted `running` jobs still become `interrupted`
- persisted `cancelling` jobs become `cancelled`
- recovered `cancelled` jobs should get a clear recovery note such as “job cancellation was in progress before restart”

This keeps restart recovery consistent with the intent already recorded before shutdown.

## Web Behavior

The job detail page should:

- show `cancelling` as a visible status
- hide the `Cancel Job` button once the job is already `cancelling`

No new routes are needed.

## Tests To Add

### Control-plane tests

Add tests proving:

- `CancelJob()` moves a running job into `cancelling` before the goroutine fully exits
- a persisted `cancelling` job is recovered as `cancelled`
- a persisted plain `running` job is still recovered as `interrupted`

### Web tests

Add tests proving:

- a `cancelling` job detail page renders correctly
- a `cancelling` job does not still show the cancel action

## Success Criteria

This slice is complete when:

1. user cancellation is immediately persisted as `cancelling`
2. restart recovery distinguishes `running` from `cancelling`
3. `cancelling` jobs recover as `cancelled`, not `interrupted`
4. the Web UI exposes `cancelling` clearly
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
