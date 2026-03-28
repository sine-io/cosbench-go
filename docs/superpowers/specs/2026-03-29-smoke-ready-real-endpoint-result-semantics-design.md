# Smoke Ready Real Endpoint Result Semantics Design

## Goal

Correct `smoke-ready` and `smoke-ready-json` so real-endpoint smoke signals stop treating workflow-level success as proof of executed live smoke.

## Problem

`Smoke S3` and `Smoke S3 Matrix` currently finish with overall workflow conclusion `success` even when the live smoke tests are skipped because `COSBENCH_SMOKE_*` secrets are absent.

Today `scripts/smoke_ready.py` maps that workflow conclusion directly to:

- `real_endpoint_latest_success = true`
- `real_endpoint_matrix_latest_success = true`

That is misleading, because the latest run artifacts show both paths can report only skipped tests rather than executed endpoint validation.

## Desired Behavior

Keep the workflow availability fields, but base real-endpoint latest success on actual smoke execution results:

- `real_endpoint_latest_success` should be true only when the latest `Smoke S3` output shows executed smoke coverage rather than all-skipped coverage
- `real_endpoint_matrix_latest_success` should be true only when the latest `Smoke S3 Matrix` aggregate shows executed rows rather than all-skipped rows

Also add:

- `real_endpoint_latest_result`
- `real_endpoint_matrix_latest_result`

Recommended values:

- single-run: `executed | skipped | failed | pending | none`
- matrix: `executed | skipped | partial | failed | pending | none`

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing `Smoke S3` or `Smoke S3 Matrix` workflow behavior
- changing legacy live or remote smoke semantics
- changing migration or comparison docs in this round

## Design

For `Smoke S3`, inspect the uploaded `smoke-s3-output` artifact and classify:

- output containing only skipped smoke tests => `skipped`
- output containing real smoke test pass lines without skip markers => `executed`
- missing or failed retrieval => `failed` / `none` / `pending` as appropriate

For `Smoke S3 Matrix`, inspect the aggregate artifact `smoke-s3-matrix-aggregate/summary.json` and classify:

- all rows show skipped test output => `skipped`
- all rows show executed smoke output => `executed`
- mixed row outcomes => `partial`

Then derive:

- `real_endpoint_latest_success = (real_endpoint_latest_result == "executed")`
- `real_endpoint_matrix_latest_success = (real_endpoint_matrix_latest_result == "executed")`

## Acceptance Criteria

- skipped `Smoke S3` runs report `real_endpoint_latest_success = false`
- skipped `Smoke S3 Matrix` runs report `real_endpoint_matrix_latest_success = false`
- both new result fields are present in JSON output
- text output includes both new result lines
- legacy live and remote smoke reporting remain unchanged
