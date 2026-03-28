# Remote Smoke Workflow Design

## Goal

Add a GitHub Actions workflow that can run the local multi-process remote smoke path on demand against a GitHub-hosted runner.

This slice is not about changing the remote protocol itself. It is about proving that the already-landed remote smoke helper works in a clean runner environment and leaves useful artifacts behind.

## Problem

The repository now has a working local helper:

- `make --no-print-directory smoke-remote-local`

It verifies:

- one local MinIO
- one `controller-only`
- two `driver-only`
- shared token auth
- work-unit distribution
- remote execution and aggregation

What is missing is an automated runner-based path that can reproduce this outside the developer machine.

Without that, the strongest remote smoke path remains local-only evidence.

## Scope

### In Scope

- one manually triggered GitHub Actions workflow
- execution of `make --no-print-directory smoke-remote-local`
- upload of `.artifacts/remote-smoke/`
- reuse of `summary.md` in the GitHub job summary
- minimal workflow hardening so failures still preserve evidence when possible

### Out Of Scope

- making remote smoke part of the default `CI` workflow
- automatic retry logic
- matrix execution of multiple smoke fixtures
- protocol changes or controller/driver feature changes

## Recommended Approach

Add one dedicated workflow under `.github/workflows/` with `workflow_dispatch` only.

### Why manual-only first

- the remote smoke path is materially heavier than unit tests and local compare helpers
- it spans multiple OS processes and a local object store
- it should prove stability before it becomes a default merge gate

This makes the workflow useful immediately without turning it into a source of friction on every push.

## Workflow Shape

The workflow should:

1. check out the repository
2. set up Go
3. run `make --no-print-directory smoke-remote-local`
4. upload `.artifacts/remote-smoke/`
5. if present, append `.artifacts/remote-smoke/summary.md` to the GitHub job summary

## Failure Strategy

The workflow must preserve evidence aggressively.

Recommended behavior:

- if `smoke-remote-local` exits non-zero, the workflow still runs artifact upload
- if `.artifacts/remote-smoke/summary.md` exists, write it to the job summary regardless of success or failure
- if the helper fails before emitting a summary, add a short fallback note to the job summary stating that artifact emission did not complete

This keeps the workflow actionable when it fails.

## Repository Touch Points

### `.github/workflows/remote-smoke-local.yml`

New workflow file with:

- `workflow_dispatch`
- one Linux job
- artifact upload
- summary export

### `README.md`

Add:

- a short description of the workflow
- one example `gh workflow run ...` command

### `scripts/smoke_remote_local.py`

Modify only if GitHub runner behavior exposes a portability problem.

Do not expand the helper’s functional scope in this slice.

## Success Criteria

This slice is complete when:

1. a manual GitHub workflow exists for remote smoke
2. the workflow runs `smoke-remote-local` on a GitHub runner
3. `.artifacts/remote-smoke/` is uploaded as an artifact
4. `summary.md` appears in the job summary when available
5. failures still preserve useful evidence

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
