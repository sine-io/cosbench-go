# Smoke Ready Structured Remote Smoke Design

## Goal

Teach `smoke-ready` and `smoke-ready-json` to prefer structured remote smoke and remote recovery summaries over workflow-level conclusions.

## Problem

The remote smoke workflows already emit structured `summary.json` and matrix aggregates, but `scripts/smoke_ready.py` still derives:

- `remote_happy_latest_success`
- `remote_recovery_latest_success`

from workflow conclusion alone.

That is weaker than the new real-endpoint and legacy-live paths, which now consume structured evidence first.

## Desired Behavior

`smoke-ready` should consume remote structured summaries first:

- `Remote Smoke Local` / `Remote Smoke Matrix` should drive:
  - `remote_happy_latest_success`
  - `remote_happy_latest_result`
- `Remote Smoke Recovery` / `Remote Smoke Recovery Matrix` should drive:
  - `remote_recovery_latest_success`
  - `remote_recovery_latest_result`

Recommended result values:

- `executed`
- `failed`
- `partial`
- `pending`
- `none`

For the happy-path pair:

- prefer `Remote Smoke Matrix` aggregate summary when it is the newer run
- otherwise use `Remote Smoke Local` summary

For the recovery pair:

- prefer `Remote Smoke Recovery Matrix` aggregate summary when it is the newer run
- otherwise use `Remote Smoke Recovery` summary

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing remote smoke workflows
- changing legacy or real-endpoint smoke semantics

## Design

Add remote structured detail loaders:

- `Remote Smoke Local` downloads `remote-smoke-output`
- `Remote Smoke Matrix` downloads `remote-smoke-matrix-aggregate`
- `Remote Smoke Recovery` downloads `remote-smoke-recovery-summary`
- `Remote Smoke Recovery Matrix` downloads `remote-smoke-recovery-matrix-aggregate`

Read `summary.json` from those artifacts and use:

- single-run remote summaries: `overall == "pass"` => `executed`, else `failed`
- matrix summaries:
  - `overall == "pass"` => `executed`
  - `overall == "partial"` => `partial`
  - otherwise `failed`

Then derive:

- `remote_happy_latest_success = (remote_happy_latest_result == "executed")`
- `remote_recovery_latest_success = (remote_recovery_latest_result == "executed")`

## Acceptance Criteria

- structured remote summaries drive `remote_happy_latest_success`
- structured remote recovery summaries drive `remote_recovery_latest_success`
- both new `remote_*_latest_result` fields are present
- old workflow presence fields remain unchanged
