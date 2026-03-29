# Smoke Ready Latest Head SHA Design

## Goal

Expose the latest evidence commit SHA for each surface in `smoke-ready` and `smoke-ready-json`.

## Problem

The readiness surface already exposes:

- result
- source
- event
- run id
- URL
- artifact
- duration
- created_at
- age_seconds

But it still does not expose which commit produced that evidence. Operators can infer it by opening the workflow page, but machine consumers should not need an extra GitHub lookup for that.

## Desired Behavior

Add `*_latest_head_sha` for:

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
- exposing branch names in this round
- changing runtime logic

## Design

Extend the latest-run normalization to capture `headSha`.

Expose it in two places:

- per-workflow `workflows.latest[*].head_sha`
- per-surface summary `*_latest_head_sha`

For aggregated remote surfaces, use the SHA from whichever workflow currently wins latest-evidence selection, matching the existing `source/event/run_id/url/artifact` behavior.

## Acceptance Criteria

- `workflows.latest[*]` includes `head_sha`
- summary exposes `*_latest_head_sha` for every latest-evidence surface
- schema contract includes those fields
- existing runtime behavior remains unchanged
