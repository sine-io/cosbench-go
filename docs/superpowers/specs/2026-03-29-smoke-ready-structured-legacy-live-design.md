# Smoke Ready Structured Legacy Live Design

## Goal

Teach `smoke-ready` and `smoke-ready-json` to prefer the new structured `Legacy Live Compare` and `Legacy Live Compare Matrix` summary artifacts over workflow-step inspection.

## Problem

Legacy live workflows now emit normalized `result.json`, but `scripts/smoke_ready.py` still derives legacy live results by inspecting GitHub job step conclusions.

That means the repository has two competing sources of truth for legacy live result state:

- workflow-emitted structured summaries
- workflow step metadata

## Desired Behavior

`smoke-ready` should prefer structured legacy summaries when available:

- `Legacy Live Compare`:
  - use `result.json.result`
  - fall back to workflow step inspection only when the summary artifact is absent
- `Legacy Live Compare Matrix`:
  - use aggregate row statuses from `summary.json`
  - fall back to job-step inspection only for older artifacts

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- small README wording update

Out of scope:

- changing legacy live workflows
- changing real-endpoint or remote smoke semantics
- removing backward-compatible step inspection entirely

## Design

Extend `load_legacy_workflow_details()` so:

- single-run downloads `legacy-live-compare-output` and prefers `result.json`
- matrix downloads `legacy-live-compare-matrix-aggregate` and prefers aggregate `summary.json`

Then update result derivation:

- single-run:
  - prefer `detail["result"]["result"]`
  - fall back to current step-based logic
- matrix:
  - prefer aggregate row `status` when already normalized (`executed/skipped/failed`)
  - fall back to current step-based logic only when needed

## Acceptance Criteria

- structured `Legacy Live Compare` summaries drive `legacy_live_latest_result`
- structured `Legacy Live Compare Matrix` aggregate statuses drive `legacy_live_matrix_latest_result`
- old workflow metadata still works as fallback
- tests cover both structured and fallback paths
