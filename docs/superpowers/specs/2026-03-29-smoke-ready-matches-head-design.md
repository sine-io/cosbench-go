# Smoke Ready Latest Matches Head Design

## Goal

Expose whether each latest-evidence surface matches the current checkout HEAD in `smoke-ready` and `smoke-ready-json`.

## Problem

The readiness surface already exposes:

- result
- source
- event
- head SHA
- head branch
- run id
- URL
- artifact
- duration
- age
- created_at

That is enough raw information to compare evidence against the current checkout, but every consumer has to re-implement the comparison logic itself.

## Desired Behavior

Add:

- top-level `current_head_sha`
- summary `*_latest_matches_head` for:
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
- short README / AGENTS / migration-gap notes

Out of scope:

- changing workflow behavior
- changing runtime logic
- comparing against refs other than the current checkout HEAD

## Design

Resolve the current checkout head once:

- `git rev-parse HEAD`

Allow tests to override it with:

- `SMOKE_READY_MOCK_CURRENT_HEAD_SHA`

Expose `current_head_sha` at the top level and derive each `*_latest_matches_head` as:

- `true` when `*_latest_head_sha == current_head_sha`
- `false` otherwise

This keeps the field simple and machine-friendly while leaving the underlying `head_sha` available for deeper inspection.

## Acceptance Criteria

- top-level `current_head_sha` is present
- summary exposes `*_latest_matches_head` for every latest-evidence surface
- schema contract includes the new fields
- existing runtime behavior remains unchanged
