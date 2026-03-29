# Smoke Ready Latest Head Branch Design

## Goal

Expose the latest evidence branch name for each surface in `smoke-ready` and `smoke-ready-json`.

## Problem

The readiness surface already exposes:

- result
- source
- event
- head SHA
- run id
- URL
- artifact
- duration
- created_at
- age_seconds

But it still does not expose which branch produced that evidence. The SHA alone is sufficient for strict correlation, but branch is the faster human cue for answering:

> Was this evidence produced from `main` or from some other branch?

## Desired Behavior

Add `*_latest_head_branch` for:

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
- changing evidence selection rules
- changing runtime logic

## Design

Extend latest-run normalization with `headBranch`.

Expose it in two places:

- per-workflow `workflows.latest[*].head_branch`
- per-surface summary `*_latest_head_branch`

For aggregated remote surfaces, use the branch from whichever workflow currently wins latest-evidence selection, matching the existing `source/event/head_sha/run_id/url` behavior.

## Acceptance Criteria

- `workflows.latest[*]` includes `head_branch`
- summary exposes `*_latest_head_branch` for every latest-evidence surface
- schema contract includes those fields
- existing runtime behavior remains unchanged
