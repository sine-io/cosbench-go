# Remote Smoke Matrix Aggregate Artifact Design

## Goal

Add a downloadable aggregate artifact to the `Remote Smoke Matrix` workflow so the combined matrix summary can be consumed outside the GitHub job summary UI.

This slice is intended to improve evidence portability, not to change helper behavior, matrix coverage, or default CI behavior.

## Problem

The repository already has:

- a non-blocking `Remote Smoke Matrix` workflow
- four matrix rows for all supported `backend × scenario` combinations
- per-row artifacts
- an `aggregate` job that emits combined `summary.md` and `summary.json`
- a combined GitHub job summary

What is still missing is a downloadable aggregate artifact containing those combined outputs.

Without that, the combined report is visible in GitHub but awkward to archive, compare, or consume programmatically.

## Scope

### In Scope

- upload `.artifacts/remote-smoke-matrix-aggregate/` as an artifact from the `aggregate` job
- lock the new upload behavior with the existing workflow-contract test
- document that the matrix workflow now emits a downloadable aggregate artifact

### Out Of Scope

- helper changes
- matrix row changes
- default CI changes
- new summary fields
- new workflows

## Recommended Approach

Add one artifact upload step to the existing `aggregate` job.

### Why this is the right boundary

- the aggregate job already produces the directory we want
- artifact upload is the smallest possible extension
- no helper, controller, or driver code needs to change
- the behavior remains non-blocking and workflow-scoped

## Upload Behavior

The new step should:

- live in the `aggregate` job
- upload `.artifacts/remote-smoke-matrix-aggregate/`
- use a stable artifact name such as `remote-smoke-matrix-aggregate`

Recommended behavior:

- run after aggregation
- not break the workflow if the aggregate output is absent unexpectedly
- preserve the current job summary behavior unchanged

## Testing Strategy

The minimal stable verification is:

1. extend the workflow-contract test to assert:
   - aggregate upload step exists
   - artifact name is stable
   - upload path matches the aggregate output directory
2. run full repository `go test ./...`
3. run full repository `go build ./...`

This slice does not need additional MinIO smoke coverage because the helper and matrix logic are unchanged.

## Repository Touch Points

### `.github/workflows/remote-smoke-matrix.yml`

Modify:

- add aggregate artifact upload step

### `scripts/test_remote_smoke_matrix_workflow.py`

Modify:

- assert aggregate upload step exists and uses the intended artifact contract

### `README.md`

Modify:

- mention the downloadable aggregate artifact

## Success Criteria

This slice is complete when:

1. the aggregate job uploads a combined aggregate artifact
2. the combined GitHub summary still works
3. the workflow-contract test covers the new upload behavior
4. the default `CI` workflow remains unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
