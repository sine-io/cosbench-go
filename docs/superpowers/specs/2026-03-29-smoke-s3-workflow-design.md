# Smoke S3 Workflow Design

## Goal

Add a manual GitHub Actions workflow that runs the existing `make smoke-s3` path against a real external S3/SIO-compatible endpoint using repository secrets.

This slice is intended to provide a stable remote entrypoint for real-endpoint validation without changing the default CI posture.

## Problem

The repository already has:

- a local `make smoke-s3` path for real endpoints
- local MinIO-based smoke workflows
- remote happy-path and recovery workflows backed by temporary local MinIO

What is still missing is a dedicated GitHub-hosted entrypoint for the real external endpoint smoke path.

That leaves live-endpoint validation dependent on local shell access even though the repository already has secrets-aware workflow patterns for other tasks.

## Scope

### In Scope

- add one manual `Smoke S3` workflow
- read required endpoint credentials from repository secrets
- support a small set of useful workflow inputs:
  - `backend`
  - `region`
  - `path_style`
  - `bucket_prefix`
- run `GO=go make smoke-s3`
- capture output to a file
- upload that file as an artifact
- append the output to the GitHub job summary

### Out Of Scope

- adding `Smoke S3` to default `CI`
- adding retries or matrix execution
- changing `internal/driver/s3` smoke logic

## Recommended Approach

Add one dedicated `workflow_dispatch` workflow rather than trying to overload existing local MinIO workflows.

### Why a dedicated workflow

- the external-endpoint smoke path has a different trust boundary than local MinIO-based workflows
- it depends on repository secrets
- operators should be able to trigger it explicitly without threading through generic workflow inputs that were designed for different scenarios

## Inputs

Required secrets:

- `COSBENCH_SMOKE_ENDPOINT`
- `COSBENCH_SMOKE_ACCESS_KEY`
- `COSBENCH_SMOKE_SECRET_KEY`

Recommended workflow inputs:

- `backend`
  - default `s3`
- `region`
  - default empty, allowing the test default to apply
- `path_style`
  - default empty
- `bucket_prefix`
  - default empty

The workflow should map those inputs to the existing environment variables consumed by `internal/driver/s3/smoke_test.go`.

## Output Contract

The workflow should:

- tee the smoke output to `smoke-s3-output.txt`
- upload `smoke-s3-output.txt` as `smoke-s3-output`
- append that same file to `$GITHUB_STEP_SUMMARY`

This mirrors the existing `Smoke Local` pattern and keeps the output surface easy to understand.

## Success Criteria

This slice is complete when:

1. a manual `Smoke S3` workflow exists
2. it can run `make smoke-s3` using repository secrets plus optional workflow inputs
3. it uploads a `smoke-s3-output` artifact
4. it writes the smoke output to the job summary
5. default `CI` remains unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
