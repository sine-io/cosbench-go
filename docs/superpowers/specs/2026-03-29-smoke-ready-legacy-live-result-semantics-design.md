# Smoke Ready Legacy Live Result Semantics Design

## Goal

Correct `smoke-ready` and `smoke-ready-json` so legacy live signals stop treating workflow-level `success` as proof of executed live validation.

## Problem

Both legacy live workflows currently return overall workflow `success` even when the actual execution step is `skipped` because repository live secrets are absent.

`scripts/smoke_ready.py` only looks at the latest run conclusion, so it currently reports:

- `legacy_live_latest_success: true`
- `legacy_live_matrix_latest_success: true`

for runs that never executed a live workload.

That conflicts with the repository docs, which already describe those runs as workflow ergonomics evidence rather than endpoint parity evidence.

## Desired Behavior

Keep the existing workflow availability booleans, but refine legacy live latest-run reporting:

- `legacy_live_latest_success` should be `true` only when the `Run legacy live compare` step actually succeeded
- `legacy_live_matrix_latest_success` should be `true` only when the matrix rows actually executed successfully
- add `legacy_live_latest_result`
- add `legacy_live_matrix_latest_result`

Recommended result values:

- single-run: `executed | skipped | failed | pending | none`
- matrix: `executed | skipped | partial | failed | pending | none`

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- README wording for the refined semantics

Out of scope:

- changing workflow behavior
- changing non-legacy smoke semantics
- changing docs outside the helper’s user-facing note

## Design

Keep the existing latest-run lookup, but include `databaseId` so the helper can inspect step-level details for the two legacy workflows.

For `Legacy Live Compare`:

- call `gh run view <databaseId> --json jobs`
- inspect the `Run legacy live compare` step
- map its conclusion to `executed`, `skipped`, `failed`, or `pending`

For `Legacy Live Compare Matrix`:

- inspect each row job whose name starts with `legacy-live-compare-matrix (`
- read the same execution step on each row
- summarize all row outcomes to one of:
  - `executed` when all rows executed
  - `skipped` when all rows skipped
  - `partial` when row outcomes differ
  - `failed` when completed rows fail without a mixed outcome
  - `pending` when the run is still in flight

## Acceptance Criteria

- skipped legacy runs produce `legacy_live_latest_success: false`
- skipped legacy matrix runs produce `legacy_live_matrix_latest_success: false`
- both new `*_latest_result` fields are present in JSON output
- text output prints both result fields
- non-legacy smoke reporting remains unchanged
