# Smoke Ready Smoke S3 Design

## Goal

Extend `smoke-ready` and `smoke-ready-json` so they include the new `Smoke S3` workflow in both workflow presence reporting and latest-run status reporting.

## Problem

The repository now has a dedicated manual `Smoke S3` workflow for real external endpoint validation, but `smoke-ready` still treats only these as smoke workflows:

- `Smoke Local`
- `Remote Smoke Local`
- `Remote Smoke Matrix`
- `Remote Smoke Recovery`
- `Remote Smoke Recovery Matrix`

That means the helper still under-reports the full smoke surface and cannot summarize the latest known real-endpoint smoke evidence.

## Scope

### In Scope

- add `Smoke S3` to the smoke workflow set
- include its latest run metadata in `smoke-ready-json`
- include its latest status in text mode
- add one summary boolean for latest real-endpoint smoke success

### Out Of Scope

- changing `Smoke S3` workflow behavior
- querying artifacts
- changing helper semantics for local env or remote happy/recovery groups

## Recommended Approach

Keep the current helper structure and add one more workflow category:

- `real_endpoint_latest_success`

This keeps the existing summary model intact while surfacing the newly added real-endpoint remote proof.

## Success Criteria

This slice is complete when:

1. `Smoke S3` appears in `workflows.present`
2. `Smoke S3` appears in `workflows.latest`
3. `summary.real_endpoint_latest_success` exists
4. text output shows the new workflow and summary line

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
