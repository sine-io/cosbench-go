# Smoke Ready Artifact Fields Design

## Goal

Extend the `smoke-ready` summary block so each latest evidence category also exposes the artifact name to download.

## Current Gap

The summary block now exposes, per category:

- `result`
- `source`
- `url`
- `created_at`

But it does not expose the artifact name. Operators still need to know the workflow internals to map:

- `Smoke S3` -> `smoke-s3-output`
- `Legacy Live Compare Matrix` -> `legacy-live-compare-matrix-aggregate`
- `Remote Smoke Recovery` -> `remote-smoke-recovery-summary`

## Desired Behavior

Add summary artifact fields:

- `real_endpoint_latest_artifact`
- `real_endpoint_matrix_latest_artifact`
- `legacy_live_latest_artifact`
- `legacy_live_matrix_latest_artifact`
- `remote_happy_latest_artifact`
- `remote_recovery_latest_artifact`

These should expose the artifact name a user would download from the corresponding workflow run.

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing workflow artifact names
- changing result/source/url/timestamp logic
- adding nested artifact metadata

## Design

Map fixed artifact names for the non-aggregated categories:

- `Smoke S3` -> `smoke-s3-output`
- `Smoke S3 Matrix` -> `smoke-s3-matrix-aggregate`
- `Legacy Live Compare` -> `legacy-live-compare-output`
- `Legacy Live Compare Matrix` -> `legacy-live-compare-matrix-aggregate`

For aggregated remote categories:

- if `remote_happy_latest_source == "Remote Smoke Local"` -> `remote-smoke-output`
- if `remote_happy_latest_source == "Remote Smoke Matrix"` -> `remote-smoke-matrix-aggregate`
- if `remote_recovery_latest_source == "Remote Smoke Recovery"` -> `remote-smoke-recovery-summary`
- if `remote_recovery_latest_source == "Remote Smoke Recovery Matrix"` -> `remote-smoke-recovery-matrix-aggregate`

## Acceptance Criteria

- JSON output includes all six artifact fields
- text output includes readable artifact lines
- remote aggregated artifact names match the selected source workflow
