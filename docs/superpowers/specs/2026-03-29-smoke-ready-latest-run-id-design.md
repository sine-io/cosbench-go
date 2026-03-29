# Smoke Ready Latest Run ID Design

## Goal

Expose the latest GitHub Actions run id for each latest-evidence surface in `smoke-ready` and `smoke-ready-json`.

## Problem

The summary block already exposes:

- result
- source
- event
- url
- artifact
- created_at

But it still does not surface the numeric run id directly. Operators can extract it from the URL, but machine consumers and triage tooling should not have to parse URLs for a stable identifier.

## Desired Behavior

Add `*_latest_run_id` fields for all current latest-evidence surfaces:

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
- `scripts/test_smoke_ready.py`
- `scripts/test_smoke_ready_schema.py`
- `scripts/test_validate_smoke_ready_schema.py`
- `docs/smoke-ready.schema.json`
- short README and migration-gap notes

Out of scope:

- changing workflow behavior
- changing evidence selection rules
- changing runtime logic

## Design

Reuse the existing normalized `workflows.latest[*].database_id` values and thread them into the summary block as `*_latest_run_id`.

For aggregated remote surfaces, use the run id from whichever workflow currently wins latest-evidence selection, matching the existing `source/url/artifact/created_at` behavior.

## Acceptance Criteria

- `smoke-ready-json` summary exposes `*_latest_run_id` for every latest-evidence surface
- schema contract includes those fields
- existing runtime behavior remains unchanged
