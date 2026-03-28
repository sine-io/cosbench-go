# Smoke Ready Legacy Live Matrix Design

## Goal

Extend `smoke-ready` and `smoke-ready-json` so they also report the new `Legacy Live Compare Matrix` workflow as a distinct readiness and latest-run signal.

## Current Problem

The helper now reports:

- `Smoke Local`
- `Smoke S3`
- `Legacy Live Compare`
- remote happy-path workflows
- remote recovery workflows

But it does not yet include `Legacy Live Compare Matrix`, so the repository readiness surface is missing the newest legacy live-validation entrypoint.

## Desired Behavior

Add `Legacy Live Compare Matrix` to the workflow and latest-run surfaces, with summary fields that remain distinct from the single-run legacy workflow:

- `legacy_live_ready`
- `legacy_live_latest_success`
- `legacy_live_matrix_ready`
- `legacy_live_matrix_latest_success`

This keeps the single-run and matrix entrypoints observable separately.

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing any workflow behavior
- collapsing single-run and matrix legacy signals into one field
- changing real-endpoint or remote smoke semantics

## Design

Add `Legacy Live Compare Matrix` to `WORKFLOW_NAMES`.

Then extend `build_payload()` with:

- `legacy_live_matrix_ready`
- `legacy_live_matrix_latest_success`

And update text output so the new workflow appears under:

- `## Workflows`
- `## Latest Runs`
- `## Summary`

## Acceptance Criteria

- `workflows.present["Legacy Live Compare Matrix"]` exists
- `workflows.latest["Legacy Live Compare Matrix"]` exists
- `summary.legacy_live_matrix_ready` exists
- `summary.legacy_live_matrix_latest_success` exists
- text output contains `Legacy Live Compare Matrix`
- existing legacy single-run fields remain unchanged
