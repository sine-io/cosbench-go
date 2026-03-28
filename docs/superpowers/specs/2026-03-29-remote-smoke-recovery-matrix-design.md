# Remote Smoke Recovery Matrix Design

## Goal

Add a non-blocking `Remote Smoke Recovery Matrix` workflow that validates recovery behavior across the two supported storage targets:

- `s3 + recovery`
- `sio + recovery`

This slice is intended to make recovery verification recurrent and remote, just like the existing happy-path remote smoke matrix.

## Problem

The repository now has:

- local MinIO-backed recovery smoke for `s3`
- local MinIO-backed recovery smoke for `sio`
- a dedicated `Remote Smoke Recovery` workflow that can run either backend on demand

What is still missing is a recurring remote workflow that exercises both recovery variants together and publishes a combined aggregate view.

Without that, recovery remains triggerable on demand but is not automatically re-validated over time.

## Scope

### In Scope

- add one new workflow:
  - `Remote Smoke Recovery Matrix`
- support:
  - `workflow_dispatch`
  - `schedule`
- run a two-row matrix:
  - `s3 + recovery`
  - `sio + recovery`
- upload per-row artifacts
- aggregate row summaries into one combined summary and artifact

### Out Of Scope

- changing helper behavior
- changing the existing `Remote Smoke Recovery` workflow
- adding recovery to default `CI`
- adding non-S3 backends

## Recommended Approach

Mirror the structure of `Remote Smoke Matrix`, but keep recovery separate.

### Why a separate workflow

- happy-path smoke and failure-recovery smoke are different verification classes
- recovery runs are slower and operationally noisier
- a separate workflow keeps the existing happy-path matrix simple
- it allows independent cadence or debugging when needed

## Workflow Shape

The workflow should contain:

1. one matrix job with:
   - `backend=s3`, `scenario=recovery`
   - `backend=sio`, `scenario=recovery`
2. one aggregate job that:
   - downloads row artifacts
   - aggregates `summary.json`
   - writes a combined Markdown summary
   - uploads an aggregate artifact

The matrix job should use:

- `fail-fast: false`

The workflow should remain:

- non-blocking
- outside the default `CI` path

## Artifact Contract

Per-row artifacts:

- `remote-smoke-recovery-s3`
- `remote-smoke-recovery-sio`

Aggregate artifact:

- `remote-smoke-recovery-matrix-aggregate`

The aggregate artifact should contain:

- `summary.json`
- `summary.md`

## Aggregation

Add a dedicated aggregation script for recovery rows, rather than overloading the happy-path matrix script.

Expected row keys:

- `backend`
- `scenario`
- `overall`
- `job_status`
- `recovery_observed`
- `reclaimed_units`
- `drivers_participated`

## Success Criteria

This slice is complete when:

1. `Remote Smoke Recovery Matrix` exists as a manual + scheduled workflow
2. it runs both `s3` and `sio` recovery rows
3. it publishes an aggregate summary and aggregate artifact
4. the default `CI` workflow remains unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
