# Remote Smoke Matrix Summary Design

## Goal

Add a workflow-level aggregate summary for the `Remote Smoke Matrix` workflow so GitHub shows one combined view of all four `backend × scenario` rows in addition to the existing per-row summaries and artifacts.

This slice is intended to improve observability, not to change the helper, the matrix shape, or the default CI path.

## Problem

The repository already has:

- a non-blocking `Remote Smoke Matrix` workflow
- four matrix rows:
  - `s3 + single`
  - `s3 + multistage`
  - `sio + single`
  - `sio + multistage`
- per-row artifact upload
- per-row summary written into each job summary

What is still missing is a single combined summary that answers, at a glance:

- which combinations ran
- which combinations passed or failed
- which combinations produced artifacts and summaries

Without that, readers have to inspect each matrix row individually.

## Scope

### In Scope

- add one `aggregate` job to the matrix workflow
- download matrix row artifacts after the row jobs complete
- aggregate each row’s `summary.json` into one combined Markdown summary
- write that combined summary to the aggregate job summary
- optionally upload a combined aggregate artifact

### Out Of Scope

- changing helper output structure
- changing the matrix row definitions
- changing the default `CI` workflow
- adding new smoke scenarios or backends

## Recommended Approach

Add a lightweight Python aggregation script under `scripts/` and call it from a new `aggregate` workflow job.

### Why a script instead of inline shell

- parsing multiple JSON files is clearer and more reliable in Python than in YAML shell fragments
- the logic becomes testable
- the workflow file stays readable
- future matrix growth remains manageable

## Aggregate Job Shape

The new workflow job should:

1. run after the matrix job
2. use `if: always()` so it still runs when one or more rows fail
3. download all artifacts matching `remote-smoke-*`
4. run the aggregation script
5. append the combined Markdown summary to `$GITHUB_STEP_SUMMARY`

The aggregate job should not make the workflow blocking for the rest of the repository. It is still part of a non-blocking workflow.

## Aggregation Script Behavior

The script should:

- scan downloaded artifact directories
- find any available `summary.json`
- produce:
  - a combined `summary.json`
  - a combined `summary.md`

Each row in the combined summary should include at least:

- `backend`
- `scenario`
- `overall`
- `job_status`
- `drivers_seen`
- `units_claimed`
- `stages_seen`

If an expected artifact is missing or lacks `summary.json`, the script should mark that row as `missing` rather than silently ignore it.

## Failure Strategy

The aggregate step should be best-effort but truthful.

Recommended behavior:

- if some row artifacts are missing, include them as `missing` in the combined summary
- if at least one row summary exists, the aggregate script should succeed and emit the combined report
- only fail the aggregate job when no row summaries are available at all

This keeps the workflow useful when part of the matrix fails.

## Repository Touch Points

### `.github/workflows/remote-smoke-matrix.yml`

Modify:

- add `aggregate` job
- download artifacts
- call the aggregation script
- append combined summary

### `scripts/aggregate_remote_smoke_matrix.py`

New file:

- reads row summaries
- emits combined Markdown and JSON

### `scripts/test_remote_smoke_matrix_workflow.py`

Extend:

- assert aggregate job exists
- assert it downloads artifacts
- assert it calls the aggregation script

### `README.md`

Update:

- mention that the matrix workflow now emits a combined summary view

## Success Criteria

This slice is complete when:

1. the matrix workflow still runs the same four rows
2. a new `aggregate` job runs after the rows and still executes on partial failure
3. GitHub displays one combined matrix summary
4. the default `CI` workflow remains unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
