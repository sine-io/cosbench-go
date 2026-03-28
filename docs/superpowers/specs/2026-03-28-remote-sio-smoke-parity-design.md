# Remote SIO Smoke Parity Design

## Goal

Extend the existing remote multi-process MinIO smoke path so it can validate both supported storage targets of the project:

- `s3`
- `sio`

The implementation should reuse the existing remote smoke helper and workflow rather than introducing a second parallel smoke system.

## Problem

The repository already has a real remote smoke path that proves:

- one `controller-only`
- two `driver-only`
- local MinIO
- remote work-unit distribution
- end-to-end aggregation

But that path currently only covers the `s3` backend.

Since the project’s functional target remains **S3 + SIO**, remote smoke evidence is asymmetric:

- S3 remote smoke: present
- SIO remote smoke: missing

That leaves one of the two primary backends without comparable remote multi-process verification.

## Scope

### In Scope

- add one SIO remote smoke fixture
- parameterize the existing local helper with `backend=s3|sio`
- parameterize the existing manual workflow with the same backend choice
- include backend identity in remote smoke summaries
- preserve the existing S3 remote smoke path

### Out Of Scope

- adding a second dedicated workflow
- matrix-running both backends in the same workflow invocation
- expanding the smoke path to cover all SIO-only operations
- quantitative legacy parity checking

## Recommended Approach

Use a single parameterized remote smoke entrypoint.

### Why parameterization is better than duplication

- one helper remains the single source of truth
- one artifact layout remains stable
- one workflow remains easier to operate
- future backends or remote smoke variants can follow the same pattern

### Configuration shape

For the local helper:

- `SMOKE_REMOTE_LOCAL_BACKEND=s3|sio`

For GitHub Actions:

- `workflow_dispatch.inputs.backend`

Default should remain:

- `s3`

## Fixture Strategy

Keep two dedicated fixtures:

- `testdata/workloads/remote-smoke-s3-two-driver.xml`
- `testdata/workloads/remote-smoke-sio-two-driver.xml`

Both fixtures should remain intentionally small and structurally similar so failures are attributable to backend path differences rather than workload complexity.

### SIO fixture guidance

The SIO fixture should:

- set `storage type="sio"`
- remain write-only for the first version
- use `workers="2"`
- use a small `totalOps`
- keep MinIO-compatible config, especially path-style access

The point is to validate the SIO remote driver path, not every SIO extension in one smoke.

## Summary And Artifact Shape

Keep the existing artifact contract but add:

- `backend`

to both:

- `summary.json`
- `summary.md`

This keeps downstream interpretation explicit.

## Checks

The SIO path should reuse the existing checks:

- process readiness
- drivers healthy
- units distributed
- job succeeded
- visibility

And add one backend-specific assertion:

- the selected backend is `sio`

The helper should also verify that the chosen fixture actually encodes the expected storage type for the selected backend.

## Workflow Changes

The manual `Remote Smoke Local` workflow should:

- accept `backend`
- pass that backend to the helper
- keep the same artifact upload and summary behavior

It should not:

- add a matrix
- change the trigger model
- become part of the default `CI` workflow yet

## Success Criteria

This slice is complete when:

1. the local helper supports `backend=s3|sio`
2. a dedicated `remote-smoke-sio-two-driver.xml` fixture exists
3. local `backend=sio` remote smoke passes
4. the manual GitHub workflow accepts `backend=sio`
5. the existing `backend=s3` path remains green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
