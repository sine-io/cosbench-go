# Remote SIO Multistage Parity Design

## Goal

Extend the existing remote smoke helper so it can validate the `sio + multistage` combination against local MinIO with the same controller-only and driver-only topology already used for other remote smoke paths.

This slice is intended to close the last obvious gap in the current remote smoke matrix:

- `s3 + single`
- `s3 + multistage`
- `sio + single`
- `sio + multistage`

## Problem

The repository already has:

- a local remote smoke helper with `backend` and `scenario` parameters
- remote smoke coverage for:
  - `s3 + single`
  - `s3 + multistage`
  - `sio + single`
- a manual GitHub workflow that can already pass both `backend` and `scenario`

What is still missing is the final fixture and helper selection path for:

- `backend=sio`
- `scenario=multistage`

Without that, the helper and workflow interfaces suggest a full matrix, but one quadrant is still unavailable.

## Scope

### In Scope

- add one `sio` multistage remote smoke fixture
- extend helper selection so `backend=sio, scenario=multistage` is valid
- add or update helper tests for the new combination
- document the new GitHub workflow invocation example
- validate the new path locally against MinIO

### Out Of Scope

- changing the remote protocol or scheduler
- changing the workflow YAML structure
- default CI integration
- adding new summary fields
- adding more than one new fixture

## Recommended Approach

Keep the current helper and workflow entrypoints unchanged. Add one new fixture and teach `fixture_for_selection(backend, scenario)` to map `sio + multistage` to it.

### Why this is the right boundary

- the existing helper interface is already the abstraction we want
- no controller or driver behavior needs to change if the SIO path is already functionally compatible
- the workflow already passes `backend` and `scenario`, so no new orchestration surface is needed
- this keeps the slice strictly about parity, not about redesign

## Fixture Strategy

Retain the current fixtures:

- `testdata/workloads/remote-smoke-s3-two-driver.xml`
- `testdata/workloads/remote-smoke-sio-two-driver.xml`
- `testdata/workloads/remote-smoke-s3-multistage-two-driver.xml`

Add:

- `testdata/workloads/remote-smoke-sio-multistage-two-driver.xml`

The new fixture should mirror the S3 multistage shape closely:

- `storage type="sio"`
- two `workstage` blocks
- one write-only work per stage
- `workers="2"` for each stage
- small `totalOps`
- distinct object ranges across the two stages

## Helper Behavior

The helper should continue to use:

- `SMOKE_REMOTE_LOCAL_BACKEND=s3|sio`
- `SMOKE_REMOTE_LOCAL_SCENARIO=single|multistage`

The only intended behavior change is that this combination becomes valid:

- `fixture_for_selection("sio", "multistage")`

All existing checks should continue to apply without backend-specific branching:

- `process_ready`
- `drivers_healthy`
- `units_distributed`
- `job_succeeded`
- `visibility`
- `stages_present`
- `stage_coverage`
- `stage_barrier`
- `stage_aggregation`

If the SIO multistage path requires custom per-backend check logic, that is a sign the slice is no longer simple parity and should stop for re-evaluation.

## Documentation

The workflow already accepts:

- `backend`
- `scenario`

So README only needs one new example showing the now-supported combination:

- `gh workflow run "Remote Smoke Local" --repo sine-io/cosbench-go -f backend=sio -f scenario=multistage`

## Success Criteria

This slice is complete when:

1. the repository contains a `remote-smoke-sio-multistage-two-driver.xml` fixture
2. `fixture_for_selection("sio", "multistage")` is valid
3. `SMOKE_REMOTE_LOCAL_BACKEND=sio SMOKE_REMOTE_LOCAL_SCENARIO=multistage make --no-print-directory smoke-remote-local` succeeds locally
4. the current workflow inputs are sufficient to trigger the new path without YAML changes

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
