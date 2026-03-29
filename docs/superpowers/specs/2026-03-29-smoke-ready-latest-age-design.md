# Smoke Ready Latest Age Design

## Goal

Expose `*_latest_age_seconds` for each latest-evidence surface in `smoke-ready` and `smoke-ready-json`.

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

But machine consumers still need to compute one common operational question themselves:

> How old is this latest evidence right now?

That requires parsing timestamps externally even though `smoke-ready` already knows both `generated_at` and each `*_latest_created_at`.

## Desired Behavior

Add `*_latest_age_seconds` for:

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

Compute age from:

- `generated_at`
- each `*_latest_created_at`

Formula:

- `age_seconds = generated_at - latest_created_at`

Expose `null` when `created_at` is missing or unparsable.

## Acceptance Criteria

- summary exposes `*_latest_age_seconds` for every latest-evidence surface
- schema contract includes those fields
- existing runtime behavior remains unchanged
