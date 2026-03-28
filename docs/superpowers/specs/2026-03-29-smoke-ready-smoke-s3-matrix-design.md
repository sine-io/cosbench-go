# Smoke Ready Smoke S3 Matrix Design

## Goal

Extend `smoke-ready` and `smoke-ready-json` so they also report the `Smoke S3 Matrix` workflow as a distinct real-endpoint matrix signal.

## Current Problem

`smoke-ready` currently tracks:

- `Smoke S3`
- `Legacy Live Compare`
- `Legacy Live Compare Matrix`
- remote smoke and recovery workflows

But it does not include `Smoke S3 Matrix`, even though that workflow is now a first-class manual entrypoint for dual-backend real-endpoint smoke validation.

That means the readiness view still only exposes the single-run real-endpoint smoke path.

## Desired Behavior

Keep the existing `real_endpoint_ready` and `real_endpoint_latest_success` semantics tied to `Smoke S3`, and add separate matrix fields:

- `real_endpoint_matrix_ready`
- `real_endpoint_matrix_latest_success`

Also include `Smoke S3 Matrix` in:

- workflow presence
- latest runs
- text summary

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing `Smoke S3` or `Smoke S3 Matrix` workflow behavior
- collapsing single-run and matrix real-endpoint signals into one field
- changing legacy or remote smoke semantics

## Design

Add `Smoke S3 Matrix` to `WORKFLOW_NAMES`.

Then extend `build_payload()` with:

- `real_endpoint_matrix_ready`
- `real_endpoint_matrix_latest_success`

and update text rendering so the new workflow appears in:

- `## Workflows`
- `## Latest Runs`
- `## Summary`

## Acceptance Criteria

- `workflows.present["Smoke S3 Matrix"]` exists
- `workflows.latest["Smoke S3 Matrix"]` exists
- `summary.real_endpoint_matrix_ready` exists
- `summary.real_endpoint_matrix_latest_success` exists
- text output contains `Smoke S3 Matrix`
- existing `real_endpoint_ready` / `real_endpoint_latest_success` fields remain unchanged
