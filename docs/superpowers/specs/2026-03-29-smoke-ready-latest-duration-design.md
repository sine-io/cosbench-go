# Smoke Ready Latest Duration Design

## Goal

Expose the latest workflow duration in seconds for each latest-evidence surface in `smoke-ready` and `smoke-ready-json`.

## Problem

The readiness surface already exposes:

- result
- source
- event
- run id
- URL
- artifact
- timestamp

But it does not expose how long the latest workflow run took. Operators can inspect that manually on GitHub, but machine consumers and dashboards should not have to scrape the UI.

## Desired Behavior

Add `*_latest_duration_seconds` for:

- real endpoint
- real endpoint matrix
- schema validation
- legacy live
- legacy live matrix
- remote happy
- remote recovery

## Scope

In scope:

- `scripts/smoke_ready.py`
- smoke-ready tests and schema contract
- schema document
- short README and migration-gap notes

Out of scope:

- changing workflow behavior
- exposing full started/updated timestamps in the public schema
- changing runtime logic

## Design

Extend the internal latest-run normalization to capture:

- `started_at`
- `updated_at`

Use those internal timestamps to derive integer `duration_seconds`.

For aggregated remote surfaces, duration comes from whichever workflow currently wins latest-evidence selection, matching the existing source/event/run-id behavior.

If a run is missing or timestamps are unavailable, expose `null`.

## Acceptance Criteria

- `smoke-ready-json` summary exposes `*_latest_duration_seconds` for every latest-evidence surface
- schema contract includes those fields
- existing runtime behavior remains unchanged
