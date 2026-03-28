# Legacy Live Structured Summary Design

## Goal

Upgrade `Legacy Live Compare` and `Legacy Live Compare Matrix` so they emit structured result summaries instead of relying only on the raw CLI summary plus workflow step metadata.

## Problem

Right now legacy live evidence is split across:

- raw CLI summary JSON from `cmd/cosbench-go`
- workflow step conclusions for executed vs skipped detection
- ad hoc skip JSON emitted by workflow preflight

That makes downstream consumers like `smoke-ready` depend on GitHub step metadata instead of a stable artifact contract.

## Desired Behavior

Each legacy live run should publish:

- the original CLI or skip summary JSON
- one normalized structured result summary

Recommended structured fields:

- `result`: `executed | skipped | failed`
- `fixture`
- `backend`

For the matrix workflow:

- each row publishes its structured result summary
- the aggregate script consumes those row summaries directly

## Scope

In scope:

- one new legacy live summary script
- `legacy-live-compare.yml`
- `legacy-live-compare-matrix.yml`
- `aggregate_legacy_live_compare_matrix.py`
- workflow/aggregate tests
- small README note

Out of scope:

- changing `cmd/cosbench-go`
- changing XML rendering
- changing `smoke-ready` in this round

## Design

Add `scripts/summarize_legacy_live_compare.py` that reads:

- the existing `.artifacts/legacy-live-compare/summary.json`
- optional `fixture`
- optional `backend`

and emits a normalized summary file, for example:

```json
{
  "result": "skipped",
  "fixture": "testdata/legacy/sio-config-sample.xml",
  "backend": "sio"
}
```

Rules:

- preflight skip JSON with `status=skipped` => `result=skipped`
- non-skip CLI summary JSON => `result=executed`
- missing/unreadable summary in completed workflow path => `result=failed`

Then:

- `Legacy Live Compare` uploads both raw and normalized summaries
- `Legacy Live Compare Matrix` rows do the same
- `aggregate_legacy_live_compare_matrix.py` prefers row normalized summaries and reports row statuses from them

## Acceptance Criteria

- single-run legacy workflow emits a normalized summary artifact
- matrix row artifacts emit normalized summaries
- matrix aggregate consumes normalized row summaries
- current skip path still works, but now through structured row status instead of step metadata alone
