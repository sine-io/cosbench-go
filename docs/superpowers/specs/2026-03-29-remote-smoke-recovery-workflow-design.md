# Remote Smoke Recovery Workflow Design

## Goal

Add a dedicated manual GitHub Actions workflow that runs the existing `scenario=recovery` remote smoke path on a GitHub-hosted runner.

This slice is intended to turn the newly landed local/manual recovery scenario into a stable, named remote verification entrypoint.

## Problem

The repository now supports recovery verification through:

- local MinIO-backed `SMOKE_REMOTE_LOCAL_SCENARIO=recovery`
- manual `Remote Smoke Local` workflow with explicit inputs

What is still missing is a dedicated workflow entrypoint for recovery.

Today the recovery path is available, but only through a parameterized workflow invocation:

- `gh workflow run "Remote Smoke Local" -f backend=s3 -f scenario=recovery`

That works, but it is more fragile and less discoverable than the existing named workflows for compare, local smoke, and remote smoke matrix.

## Scope

### In Scope

- add one new `workflow_dispatch` workflow:
  - `Remote Smoke Recovery`
- fix its behavior to:
  - `backend=s3`
  - `scenario=recovery`
- upload `.artifacts/remote-smoke/`
- write `.artifacts/remote-smoke/summary.md` into the job summary
- add a lightweight workflow contract test

### Out Of Scope

- modifying helper behavior
- adding recovery to default `CI`
- changing `Remote Smoke Local`
- adding SIO recovery in the same slice

## Recommended Approach

Add a separate workflow file rather than overloading `Remote Smoke Local`.

### Why a dedicated workflow

- recovery is operationally different from the normal happy-path smoke
- naming it explicitly makes it easy to trigger and reason about
- it keeps the current parameterized workflow intact for ad hoc exploration
- it matches the current pattern of named manual verification workflows

## Workflow Shape

The workflow should:

1. be `workflow_dispatch` only
2. check out the repository
3. set up Go
4. run:
   - `SMOKE_REMOTE_LOCAL_BACKEND=s3`
   - `SMOKE_REMOTE_LOCAL_SCENARIO=recovery`
   - `GO=go make --no-print-directory smoke-remote-local`
5. upload `.artifacts/remote-smoke/`
6. append `summary.md` to the GitHub job summary when present

## Testing Strategy

The smallest stable verification is:

1. add a workflow-contract test for:
   - file existence
   - `workflow_dispatch`
   - fixed `backend=s3`
   - fixed `scenario=recovery`
   - artifact upload and summary export
2. run full `go test ./...`
3. run full `go build ./...`

This slice does not require local MinIO smoke reruns because the helper and runtime logic remain unchanged.

## Success Criteria

This slice is complete when:

1. a manual `Remote Smoke Recovery` workflow exists
2. it always runs `backend=s3`, `scenario=recovery`
3. it uploads remote smoke artifacts and writes the summary
4. the repository still builds and tests cleanly

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
