# Smoke Ready Structured Smoke S3 Design

## Goal

Teach `smoke-ready` and `smoke-ready-json` to prefer the new structured `Smoke S3` and `Smoke S3 Matrix` summary artifacts over heuristic parsing of raw text output.

## Problem

The repository now emits structured smoke summaries:

- `Smoke S3` publishes `summary.json`
- `Smoke S3 Matrix` rows publish `summary.json`
- the matrix aggregate reflects structured row statuses

But `scripts/smoke_ready.py` still derives real-endpoint result states mainly by parsing raw text. That leaves duplicated parsing logic in two places:

- the workflow-side summary generator
- the readiness helper

## Desired Behavior

`smoke-ready` should prefer structured summaries when available:

- `Smoke S3`:
  - use row `summary.json.result`
  - fall back to raw text parsing only when summary is absent
- `Smoke S3 Matrix`:
  - use aggregate row `status` values directly when they are already normalized
  - fall back to raw text parsing only for older artifacts

This keeps backward compatibility while making the current path deterministic.

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing workflow behavior
- changing legacy live or remote smoke semantics
- removing backward-compatible text parsing entirely

## Design

For `Smoke S3`, allow `load_real_endpoint_details()` to load:

- raw text output
- structured summary JSON

Then compute `real_endpoint_latest_result` by:

1. using `summary["result"]` when present
2. otherwise falling back to text parsing

For `Smoke S3 Matrix`, let `smoke_matrix_result()`:

1. use row `status` directly when it is already one of:
   - `executed`
   - `skipped`
   - `failed`
2. only fall back to parsing `row["output"]` when the row still uses the old `present` representation

## Acceptance Criteria

- structured `Smoke S3` summaries drive `real_endpoint_latest_result`
- structured `Smoke S3 Matrix` row statuses drive `real_endpoint_matrix_latest_result`
- old artifact shapes still work as fallback
- tests cover both structured and fallback paths
