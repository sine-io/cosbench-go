# Smoke Ready Unified Source Schema Design

## Goal

Complete the `smoke-ready` summary schema so every latest evidence category exposes the same four concepts:

- `result`
- `source`
- `url`
- `created_at`

## Current Gap

The summary schema is still inconsistent:

- `remote_happy` and `remote_recovery` already expose `*_latest_source`
- `real_endpoint`, `real_endpoint_matrix`, `legacy_live`, and `legacy_live_matrix` do not

For those categories, the source is implicit, but not carried in the summary block itself.

## Desired Behavior

Add these fields:

- `real_endpoint_latest_source`
- `real_endpoint_matrix_latest_source`
- `legacy_live_latest_source`
- `legacy_live_matrix_latest_source`

Values should be the workflow display names:

- `Smoke S3`
- `Smoke S3 Matrix`
- `Legacy Live Compare`
- `Legacy Live Compare Matrix`

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing result logic
- changing URL/timestamp logic
- adding nested provenance objects

## Design

These four fields are direct mappings:

- `real_endpoint_latest_source = "Smoke S3"`
- `real_endpoint_matrix_latest_source = "Smoke S3 Matrix"`
- `legacy_live_latest_source = "Legacy Live Compare"`
- `legacy_live_matrix_latest_source = "Legacy Live Compare Matrix"`

Add them to JSON summary output and print them in text mode near the paired result/URL/timestamp lines.

## Acceptance Criteria

- JSON output includes the four new source fields
- text output includes readable source lines for those categories
- remote source fields remain unchanged
