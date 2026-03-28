# Smoke Ready Source Fields Design

## Goal

Extend `smoke-ready` and `smoke-ready-json` so the aggregated remote result fields also report which workflow produced the current latest result.

## Problem

`smoke-ready` now exposes:

- `remote_happy_latest_result`
- `remote_recovery_latest_result`

but it does not expose which workflow those values came from.

That makes it harder to tell whether the current source of truth is:

- `Remote Smoke Local` or `Remote Smoke Matrix`
- `Remote Smoke Recovery` or `Remote Smoke Recovery Matrix`

without manually comparing timestamps in the full workflow table.

## Desired Behavior

Add two source fields:

- `remote_happy_latest_source`
- `remote_recovery_latest_source`

Values should be the workflow display names:

- `Remote Smoke Local`
- `Remote Smoke Matrix`
- `Remote Smoke Recovery`
- `Remote Smoke Recovery Matrix`
- `none`

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing workflow behavior
- changing how latest result itself is computed
- adding source fields for non-aggregated categories

## Design

The helper already chooses the newer run between the single-run and matrix variants.
Capture that chosen workflow name in two new summary fields:

- `remote_happy_latest_source`
- `remote_recovery_latest_source`

Then print them in text mode under `## Summary`.

## Acceptance Criteria

- JSON output includes both new source fields
- text output includes both new source lines
- source values match the workflow that supplied the current latest result
