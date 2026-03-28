# Remote Smoke Workflow Scenario Design

## Goal

Extend the existing manual `Remote Smoke Local` GitHub workflow so it can also run the new multistage remote smoke scenario.

This slice is intended to close the gap between:

- local MinIO evidence for `SMOKE_REMOTE_LOCAL_SCENARIO=multistage`
- GitHub runner evidence for the same remote multistage smoke path

## Problem

The repository already has:

- a local remote smoke helper with `backend=s3|sio`
- local multistage smoke support via `SMOKE_REMOTE_LOCAL_SCENARIO=multistage`
- a manual GitHub workflow that can run the helper with a `backend` input

What is still missing is workflow-level support for the `scenario` dimension.

Without that, the strongest multistage remote smoke evidence remains local-only.

## Scope

### In Scope

- extend the existing manual workflow with a `scenario` input
- pass `backend` and `scenario` through the existing `smoke-remote-local` helper path
- update README invocation examples
- add a small automated test or sanity check that locks the workflow contract

### Out Of Scope

- creating a separate workflow file
- making remote smoke part of default `CI`
- adding `sio + multistage` support in the workflow if the helper does not support it
- changing the helper’s functional behavior in this slice
- protocol, scheduler, or controller changes

## Recommended Approach

Keep the current `Remote Smoke Local` workflow as the single manual entrypoint and add one more input:

- `backend`
- `scenario`

### Why extend the existing workflow

- it preserves one stable artifact and summary location
- it avoids splitting a small manual verification surface into multiple workflows
- it matches the existing helper model, which already treats scenario as an input
- it keeps future additions composable without creating workflow sprawl

## Workflow Shape

The workflow should remain `workflow_dispatch` only.

It should accept:

- `backend`, default `s3`
- `scenario`, default `single`

The run step should call the same helper path as today, with both variables threaded through:

- `SMOKE_REMOTE_LOCAL_BACKEND`
- `SMOKE_REMOTE_LOCAL_SCENARIO`

The artifact and summary behavior should stay unchanged:

- upload `.artifacts/remote-smoke/`
- append `.artifacts/remote-smoke/summary.md` to the GitHub step summary when present

## Supported Combinations

This slice should explicitly support:

- `backend=s3, scenario=single`
- `backend=s3, scenario=multistage`
- `backend=sio, scenario=single`

This slice should not promise:

- `backend=sio, scenario=multistage`

If needed later, that can be added in a separate parity slice after the local helper supports it.

## Testing Strategy

Because this change is mostly YAML and documentation, the smallest stable verification is:

1. add a lightweight test that checks the workflow file contains a `scenario` input and passes it into the run command
2. run the local helper manually for:
   - default single-stage
   - multistage
3. verify the current manual workflow YAML still parses

This avoids trying to simulate GitHub Actions locally while still locking the intended contract.

## Repository Touch Points

### `.github/workflows/remote-smoke-local.yml`

Modify:

- add `workflow_dispatch.inputs.scenario`
- thread the scenario env var into the smoke command

### `README.md`

Update:

- workflow invocation examples
- explain that `scenario=multistage` is now available for the manual workflow

### `scripts/`

Optionally add one lightweight test file that checks workflow text shape. Do not add runtime orchestration logic here.

## Success Criteria

This slice is complete when:

1. the manual workflow accepts `scenario`
2. the workflow passes both `backend` and `scenario` to the helper
3. README shows how to trigger the multistage workflow run
4. local verification proves the helper still passes for default single-stage and multistage scenarios

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
