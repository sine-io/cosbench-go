# Smoke Ready Schema Validation Reporting Design

## Goal

Expose the new `Smoke Ready Validate` workflow through `smoke-ready` / `smoke-ready-json` so the unified readiness view includes the latest schema-validation evidence.

## Problem

The repository now has:

- `Smoke Ready Validate` as a manual GitHub workflow
- `smoke-ready-validate` and `smoke-ready-validate-json` as local entrypoints

But `scripts/smoke_ready.py` does not track that workflow anywhere. The readiness surface still reports only:

- local smoke
- real-endpoint smoke
- legacy live compare
- remote happy-path and recovery

That leaves the schema-validation evidence outside the repository’s main status view.

## Desired Behavior

Add `Smoke Ready Validate` to the `smoke-ready` surface with:

- workflow presence and latest run metadata
- summary readiness
- latest result
- latest source
- latest URL
- latest artifact
- latest created_at

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- `scripts/test_smoke_ready_schema.py`
- `docs/smoke-ready.schema.json`
- short doc updates

Out of scope:

- changing the validator helper itself
- changing workflow semantics
- adding default CI gating

## Design

Treat `Smoke Ready Validate` as a separate contract-surface signal, not as part of remote or real-endpoint smoke.

Add one new workflow constant and include it in:

- `WORKFLOW_NAMES`
- workflow presence/latest reporting
- artifact-name mapping

Download `smoke-ready-validate-output` and read `validation.json`.

Classify its latest result as:

- `validated` when `validation.json.valid == true`
- `failed` when the workflow completed but validation says false or the artifact is unreadable
- `pending` while the workflow run is still in progress
- `none` when no run exists

Expose these summary fields:

- `schema_validation_ready`
- `schema_validation_latest_success`
- `schema_validation_latest_result`
- `schema_validation_latest_source`
- `schema_validation_latest_url`
- `schema_validation_latest_artifact`
- `schema_validation_latest_created_at`

## Acceptance Criteria

- `smoke-ready-json` includes `Smoke Ready Validate` in workflow presence/latest reporting
- summary includes the new schema-validation fields
- schema contract is updated accordingly
- existing smoke semantics remain unchanged
