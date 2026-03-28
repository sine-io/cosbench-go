# Smoke S3 Matrix Design

## Goal

Add a manual matrix workflow that runs the real-endpoint smoke path for both supported backends:

- `s3`
- `sio`

This slice is intended to provide a single remote entrypoint for real-endpoint parity checks without changing the default CI workflow.

## Problem

The repository already has:

- a manual `Smoke S3` workflow
- support for `backend=s3` and `backend=sio`

What is still missing is a one-shot workflow that validates both real-endpoint variants together and produces an aggregate summary.

Without that, operators must trigger `Smoke S3` twice and manually compare the results.

## Scope

### In Scope

- add a manual `Smoke S3 Matrix` workflow
- run two matrix rows:
  - `backend=s3`
  - `backend=sio`
- reuse the existing `make smoke-s3` path
- use the same repository secrets as `Smoke S3`
- upload per-row artifacts
- aggregate row outputs into a combined summary and aggregate artifact

### Out Of Scope

- scheduled live-endpoint runs
- default CI integration
- changes to the `internal/driver/s3` smoke test implementation

## Recommended Approach

Mirror the current `Remote Smoke Matrix` structure, but use the real-endpoint `make smoke-s3` command instead of the local MinIO helper.

### Why a separate workflow

- `Smoke S3` remains useful as the simplest single-backend manual entrypoint
- the matrix workflow serves a different purpose: backend parity verification against the real external endpoint
- separating them keeps both workflows easy to understand

## Workflow Shape

The workflow should contain:

1. one `workflow_dispatch` matrix job with:
   - `backend=s3`
   - `backend=sio`
2. one aggregate job that:
   - downloads row artifacts
   - aggregates the row output text files
   - writes a combined Markdown summary
   - uploads an aggregate artifact

It should not add `schedule` at this stage because live-endpoint runs may have cost or operational sensitivity.

## Inputs

Use the same optional inputs as `Smoke S3` where practical:

- `region`
- `path_style`
- `bucket_prefix`

The matrix should own `backend`, so the user does not need to pass it.

## Output Contract

Per-row artifacts:

- `smoke-s3-s3`
- `smoke-s3-sio`

Aggregate artifact:

- `smoke-s3-matrix-aggregate`

The aggregate summary should include at least:

- backend
- row status
- whether smoke output exists

## Success Criteria

This slice is complete when:

1. a manual `Smoke S3 Matrix` workflow exists
2. it runs both `s3` and `sio` rows
3. it uploads per-row artifacts plus one aggregate artifact
4. the default `CI` workflow remains unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
