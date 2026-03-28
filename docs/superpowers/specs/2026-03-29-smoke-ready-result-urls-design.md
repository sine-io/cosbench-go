# Smoke Ready Result URLs Design

## Goal

Extend `smoke-ready` and `smoke-ready-json` so the summary block includes direct URLs for the latest evidence runs.

## Problem

Today the summary block exposes:

- latest results
- latest success booleans
- source workflow names for aggregated remote categories

But it does not expose direct run URLs in the summary section. The URLs are only present in the lower-level `workflows.latest` map.

That means consumers who only read the summary still need extra logic to jump to the underlying workflow run.

## Desired Behavior

Add summary URL fields for the current evidence categories:

- `real_endpoint_latest_url`
- `real_endpoint_matrix_latest_url`
- `legacy_live_latest_url`
- `legacy_live_matrix_latest_url`
- `remote_happy_latest_url`
- `remote_recovery_latest_url`

For aggregated categories:

- use the URL of whichever workflow currently supplies the chosen latest result

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing workflow behavior
- changing result computation logic
- adding nested provenance structures

## Design

Reuse existing latest-run data and chosen source workflow names.

Single-source categories map directly:

- `real_endpoint_latest_url` => `workflows.latest["Smoke S3"].url`
- `real_endpoint_matrix_latest_url` => `workflows.latest["Smoke S3 Matrix"].url`
- `legacy_live_latest_url` => `workflows.latest["Legacy Live Compare"].url`
- `legacy_live_matrix_latest_url` => `workflows.latest["Legacy Live Compare Matrix"].url`

Aggregated remote categories map indirectly:

- `remote_happy_latest_url` => URL for `remote_happy_latest_source`
- `remote_recovery_latest_url` => URL for `remote_recovery_latest_source`

Then print those URL fields in the text summary after their paired source/result lines.

## Acceptance Criteria

- JSON output includes all six new URL fields
- text output includes readable URL lines
- remote aggregated URLs match the selected source workflow
