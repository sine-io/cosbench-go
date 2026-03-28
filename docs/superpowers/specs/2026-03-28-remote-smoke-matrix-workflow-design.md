# Remote Smoke Matrix Workflow Design

## Goal

Add a non-blocking GitHub Actions workflow that runs the remote smoke matrix across the current supported `backend × scenario` combinations on a schedule and on manual demand.

This slice is intended to turn the already-landed local and manual one-off smoke evidence into a recurring remote verification path without making it part of the default merge gate.

## Problem

The repository already has:

- local MinIO-backed remote smoke for:
  - `s3 + single`
  - `s3 + multistage`
  - `sio + single`
  - `sio + multistage`
- a manual single-run workflow that can trigger one remote smoke combination at a time

What is still missing is a recurring, consolidated remote proof that all currently supported combinations still work on GitHub-hosted runners.

Without that, the matrix remains manually explorable but not automatically re-validated over time.

## Scope

### In Scope

- add one new workflow dedicated to remote smoke matrix execution
- run the four supported combinations:
  - `s3 + single`
  - `s3 + multistage`
  - `sio + single`
  - `sio + multistage`
- support both `schedule` and `workflow_dispatch`
- upload artifacts per matrix row
- reuse each job’s `.artifacts/remote-smoke/summary.md` in the job summary

### Out Of Scope

- changing the existing `Remote Smoke Local` single-run workflow
- adding the matrix workflow to default `CI`
- changing helper behavior
- changing controller or driver code

## Recommended Approach

Add a separate workflow:

- `Remote Smoke Matrix`

Keep the existing workflow:

- `Remote Smoke Local`

### Why a separate workflow

- the current single-run workflow is useful as a targeted manual tool
- the matrix workflow serves a different purpose: broad recurring validation
- separating them keeps both entrypoints simple
- it avoids overloading one YAML file with “single run” and “run all combinations” behavior

## Trigger Model

Use:

- `schedule`
- `workflow_dispatch`

The workflow should not be triggered by:

- `push`
- `pull_request`

This keeps the matrix non-blocking and aligned with the project’s current goal of opt-in smoke coverage outside the default CI gate.

## Matrix Shape

One job with `strategy.matrix.include` is sufficient.

Recommended rows:

- `backend=s3`, `scenario=single`
- `backend=s3`, `scenario=multistage`
- `backend=sio`, `scenario=single`
- `backend=sio`, `scenario=multistage`

Recommended job behavior:

- `fail-fast: false` so one row failing does not prevent the others from producing evidence
- per-row artifact names, for example:
  - `remote-smoke-s3-single`
  - `remote-smoke-s3-multistage`
  - `remote-smoke-sio-single`
  - `remote-smoke-sio-multistage`

## Runtime Behavior

Each matrix row should:

1. check out the repository
2. set up Go
3. run:
   - `SMOKE_REMOTE_LOCAL_BACKEND`
   - `SMOKE_REMOTE_LOCAL_SCENARIO`
   - `GO=go make --no-print-directory smoke-remote-local`
4. upload `.artifacts/remote-smoke/`
5. append `.artifacts/remote-smoke/summary.md` to the GitHub job summary when present

## Failure Strategy

The matrix workflow should preserve evidence aggressively:

- `fail-fast: false`
- artifact upload must still run with `if: always()`
- summary export must still run with `if: always()`

If a row fails before the helper emits `summary.md`, the workflow should write a short fallback note to the job summary for that row.

## Repository Touch Points

### `.github/workflows/remote-smoke-matrix.yml`

New workflow file containing:

- `schedule`
- `workflow_dispatch`
- matrix job
- artifact upload
- summary export

### `README.md`

Add:

- a short note that a recurring matrix workflow exists
- one `gh workflow run "Remote Smoke Matrix"` example

### `scripts/`

Add one small workflow-contract test so the matrix shape is locked by automation instead of manual inspection.

## Success Criteria

This slice is complete when:

1. a new non-blocking matrix workflow exists
2. it covers all four supported remote smoke combinations
3. every matrix row uploads its own artifact and summary
4. the default `CI` workflow remains unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
