# Smoke Ready Latest Event Design

## Goal

Teach `smoke-ready` and `smoke-ready-json` to report which GitHub event produced the latest evidence for each surface.

## Problem

Several workflows now support multiple triggers:

- `workflow_dispatch`
- `schedule`
- `push` for some surfaces

But `smoke-ready` currently exposes only:

- source workflow name
- URL
- artifact
- timestamp

That leaves one useful question unanswered:

> Was this latest evidence produced manually, by schedule, or by push?

## Desired Behavior

Expose the latest triggering event in two places:

- per-workflow `workflows.latest[*].event`
- per-surface summary `*_latest_event`

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- `scripts/test_smoke_ready_schema.py`
- `scripts/test_validate_smoke_ready_schema.py`
- `docs/smoke-ready.schema.json`
- short README and migration-gap wording

Out of scope:

- changing workflow behavior
- renaming existing summary fields
- adding new workflows

## Design

Extend `load_workflow_latest_runs()` to request and normalize `event`.

Then add summary event fields for:

- `real_endpoint_latest_event`
- `real_endpoint_matrix_latest_event`
- `schema_validation_latest_event`
- `legacy_live_latest_event`
- `legacy_live_matrix_latest_event`
- `remote_happy_latest_event`
- `remote_recovery_latest_event`

For single-workflow surfaces, the event is taken directly from that workflow’s latest run.
For aggregated remote surfaces, the event is taken from whichever workflow currently wins the latest-result selection.

## Acceptance Criteria

- `workflows.latest[*]` includes `event`
- summary includes `*_latest_event` for every latest-evidence surface
- schema contract is updated
- existing runtime behavior remains unchanged
