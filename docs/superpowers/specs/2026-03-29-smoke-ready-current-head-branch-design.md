# Smoke Ready Current Head Branch Design

## Goal

Expose the current checkout branch alongside `current_head_sha` in `smoke-ready` and `smoke-ready-json`.

## Problem

The helper already exposes:

- top-level `current_head_sha`
- per-surface latest `head_branch`

But it still does not expose the current checkout branch itself. That means consumers can compare evidence branches to the current branch only by separately asking Git.

## Desired Behavior

Add top-level:

- `current_head_branch`

## Scope

In scope:

- `scripts/smoke_ready.py`
- smoke-ready tests and schema contract
- schema document
- short README / AGENTS / migration-gap notes

Out of scope:

- changing workflow behavior
- changing runtime logic

## Design

Resolve the current branch once:

- `git rev-parse --abbrev-ref HEAD`

Allow tests to override it with:

- `SMOKE_READY_MOCK_CURRENT_HEAD_BRANCH`

Expose it as top-level `current_head_branch`.

## Acceptance Criteria

- top-level `current_head_branch` is present
- schema contract includes the field
- existing runtime behavior remains unchanged
