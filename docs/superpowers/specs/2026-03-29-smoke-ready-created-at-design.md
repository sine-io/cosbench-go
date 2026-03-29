# Smoke Ready Created-At Design

## Goal

Add `*_latest_created_at` fields to the `smoke-ready` summary block so operators can read result, source, URL, and timestamp from one place.

## Problem

`smoke-ready` already exposes:

- result fields
- source fields for aggregated remote categories
- direct run URLs

But the timestamp still lives only in `workflows.latest`. Summary-only consumers must drill into the lower-level structure to answer “how old is this evidence?”

## Desired Behavior

Add summary timestamp fields:

- `real_endpoint_latest_created_at`
- `real_endpoint_matrix_latest_created_at`
- `legacy_live_latest_created_at`
- `legacy_live_matrix_latest_created_at`
- `remote_happy_latest_created_at`
- `remote_recovery_latest_created_at`

For aggregated remote categories, use the timestamp from whichever workflow currently supplies the chosen latest result.

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- freshness scoring or staleness thresholds
- changing workflow behavior
- changing any existing result logic

## Design

Reuse `workflows.latest` and the already-selected aggregated source workflow names.

Map summary timestamps directly:

- single-source categories => their corresponding workflow `created_at`
- aggregated remote categories => `created_at` of `remote_*_latest_source`

Then print the new lines in text mode near the matching URL fields.

## Acceptance Criteria

- JSON output includes all six `*_latest_created_at` fields
- text output includes readable timestamp lines
- aggregated timestamps match the selected source workflow
