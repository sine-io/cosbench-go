# Remote Multistage Smoke Design

## Goal

Extend the existing local multi-process remote smoke path so it can also prove multi-stage remote job progression against a local MinIO-backed endpoint.

This slice is intended to close the gap between:

- unit and integration evidence that remote jobs can advance from one stage to the next
- system-level smoke evidence that the same behavior works across real controller-only and driver-only processes

## Problem

The repository already has:

- remote multi-process local smoke with one controller-only process and two driver-only processes
- backend parity for `s3` and `sio`
- remote work-unit scheduling across multiple drivers
- multi-stage remote progression covered in control-plane and app integration tests

What is still missing is a repeatable local smoke path that proves multi-stage progression in the real multi-process topology.

Without that, the strongest evidence for multi-stage remote progression remains below the system-smoke layer.

## Scope

### In Scope

- keep the current single-stage remote smoke path intact
- add one dedicated multistage remote smoke scenario for local MinIO-backed execution
- add one minimal S3 multistage fixture
- add scenario-aware summary fields and multistage checks
- verify multistage progression with one controller-only process and two driver-only processes

### Out Of Scope

- changing the remote protocol or scheduler behavior itself
- adding multistage support to the GitHub workflow in this slice
- adding SIO multistage smoke in the same change
- performance or stress testing
- browser automation

## Recommended Approach

Keep the current smoke helper as the single entrypoint, but add a scenario parameter:

- `single`
- `multistage`

### Why scenario parameterization

- it preserves the current minimal remote smoke path as the default
- it avoids duplicating process orchestration logic
- it keeps artifacts and command shape stable
- it creates a clean path for later workflow or backend expansion without multiplying scripts

## Fixture Strategy

Retain the current fixture:

- `testdata/workloads/remote-smoke-s3-two-driver.xml`

Add one new multistage fixture:

- `testdata/workloads/remote-smoke-s3-multistage-two-driver.xml`

Recommended properties:

- `storage type="s3"`
- two stages in strict sequence
- each stage contains one write-only work
- each work uses `workers="2"`
- small `totalOps`
- object ranges do not overlap across stages

The first version should stay on `s3` only because this slice is about stage progression, not backend parity. `sio` multistage smoke can build on the same scenario machinery later.

## Multistage Checks

The helper should preserve existing checks and add multistage-specific checks when `scenario=multistage`.

### Existing checks retained

- process readiness
- driver registration and health
- unit distribution across drivers
- job success
- basic controller and driver visibility

### New checks

#### 1. Stage presence

- the job exposes at least two stages
- both stages reach `succeeded`

#### 2. Stage mission coverage

- controller mission snapshots include attempts for both stage names

#### 3. Stage barrier ordering

- the second stage must not start before the first stage finishes
- the simplest proof is:
  - stage A has `finished_at`
  - stage B has `started_at`
  - `stage_a.finished_at <= stage_b.started_at`

#### 4. Aggregation across stages

- result stage totals contain at least two stage entries
- total operation count and byte count remain non-zero

## Summary Artifacts

Keep the current artifact location:

- `.artifacts/remote-smoke/`

Extend the summary shape with multistage context:

- `scenario`
- `stage_names`
- `stages_seen`

These should appear in both:

- `summary.json`
- `summary.md`

The existing fields should remain stable so current consumers do not break.

## Implementation Shape

The smallest coherent implementation is:

1. add the multistage fixture
2. add failing helper tests for scenario selection and summary shape
3. parameterize the helper by scenario
4. add multistage-specific checks using existing controller API payloads and mission snapshots
5. verify both the default single-stage smoke and the new multistage smoke locally against MinIO

## Success Criteria

This slice is complete when:

1. `make --no-print-directory smoke-remote-local` still validates the current single-stage path
2. `SMOKE_REMOTE_LOCAL_SCENARIO=multistage make --no-print-directory smoke-remote-local` succeeds locally
3. multistage summaries explicitly record `scenario=multistage`
4. the multistage smoke proves two stages succeeded in order across real controller-only and driver-only processes

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
