# Smoke Ready Legacy Live Design

## Goal

Extend `smoke-ready` and `smoke-ready-json` so they report `Legacy Live Compare` as a distinct workflow/readiness signal alongside the existing local, real-endpoint smoke, remote happy-path, and remote recovery views.

## Current Problem

`Legacy Live Compare` now exists as a manual real-endpoint workflow with stable skip semantics when secrets are absent, but `scripts/smoke_ready.py` does not include it anywhere:

- it is missing from the workflow presence list
- it is missing from the latest-run view
- there is no summary field that tells operators whether legacy live validation is available or recently successful

That leaves the repository readiness dashboard blind to the newest live-compare surface.

## Desired Behavior

`smoke-ready` should treat `Legacy Live Compare` as its own readiness category:

- include `Legacy Live Compare` in workflow presence
- include its latest run status/conclusion/time/url
- add `legacy_live_ready`
- add `legacy_live_latest_success`

This should remain distinct from:

- `real_endpoint_ready`, which continues to mean the `Smoke S3` live smoke path exists
- `real_endpoint_latest_success`, which continues to reflect only `Smoke S3`

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README note if needed for operator discoverability

Out of scope:

- changing `Legacy Live Compare` workflow behavior
- changing `Smoke S3` semantics
- changing remote smoke workflows or matrix aggregation

## Design

Add `Legacy Live Compare` to `WORKFLOW_NAMES` in `scripts/smoke_ready.py`.

Then extend the summary object with:

- `legacy_live_ready`: true when the workflow is present
- `legacy_live_latest_success`: true when the latest `Legacy Live Compare` run concluded `success`

Update text rendering so the new workflow appears under both:

- `## Workflows`
- `## Latest Runs`

And add the two new summary lines under `## Summary`.

## Acceptance Criteria

- `smoke-ready-json` includes `workflows.present["Legacy Live Compare"]`
- `smoke-ready-json` includes `workflows.latest["Legacy Live Compare"]`
- `summary.legacy_live_ready` exists
- `summary.legacy_live_latest_success` exists
- text output prints `Legacy Live Compare` in workflow and latest-run sections
- existing fields for `Smoke S3`, remote happy-path, and remote recovery remain unchanged
