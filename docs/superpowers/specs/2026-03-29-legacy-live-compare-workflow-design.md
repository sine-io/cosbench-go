# Legacy Live Compare Workflow Design

## Goal

Add a manual GitHub Actions workflow that renders a representative legacy sample into a real runnable workload using repository secrets, then runs the existing CLI against a real endpoint and captures the result as a structured artifact.

This slice is intended to turn the current live-compare runbook into a stable remote entrypoint without changing default CI.

## Problem

The repository already has:

- local runbooks for live endpoint comparison
- representative legacy fixtures:
  - `testdata/legacy/s3-config-sample.xml`
  - `testdata/legacy/sio-config-sample.xml`
- a stable CLI path:
  - `go run ./cmd/cosbench-go ... -json -quiet -summary-file`

What is still missing is a repeatable GitHub-hosted path for running one of those legacy fixtures against a real endpoint.

The main blocker is that the legacy sample XML files contain placeholder storage config values and cannot be run directly without rendering in real credentials and endpoint values.

## Scope

### In Scope

- add a small render script that materializes a runnable XML from a legacy fixture plus `COSBENCH_SMOKE_*`
- add a manual `Legacy Live Compare` workflow
- run one selected legacy fixture at a time
- upload:
  - rendered workload XML
  - CLI summary JSON
  - stdout/stderr log
- write a short job summary

### Out Of Scope

- default CI integration
- matrix execution
- automatic comparison against legacy Java output
- editing the canonical legacy fixture files in `testdata/legacy/`

## Recommended Approach

Use a render step plus the existing CLI.

### Why render instead of changing the CLI

- the legacy sample files are intentionally stored as portable references with placeholder values
- a small render helper isolates the transformation cleanly
- the existing CLI remains unchanged
- the same helper can be reused locally if needed

## Workflow Shape

The workflow should be `workflow_dispatch` only.

Recommended inputs:

- `fixture`
  - default `testdata/legacy/sio-config-sample.xml`
- `backend`
  - default `sio`
- `region`
  - optional
- `path_style`
  - optional

Required secrets:

- `COSBENCH_SMOKE_ENDPOINT`
- `COSBENCH_SMOKE_ACCESS_KEY`
- `COSBENCH_SMOKE_SECRET_KEY`

Workflow flow:

1. checkout
2. setup Go
3. render selected legacy fixture into a temporary runnable XML
4. run:
   - `go run ./cmd/cosbench-go <rendered> -backend <backend> -json -quiet -summary-file <path>`
5. upload artifacts
6. write summary

## Render Script Behavior

The render helper should:

- read the input legacy XML
- replace storage config placeholders:
  - `<accesskey>`
  - `<scretkey>`
  - `<endpoint>`
  - optional `path_style_access`
  - optional `region`
- preserve the original storage type in the XML

It should not:

- rewrite stage/work/op structure
- normalize or validate beyond template substitution

## Output Contract

The workflow should produce an artifact directory containing:

- `rendered-workload.xml`
- `summary.json`
- `run.log`

The job summary should include:

- fixture path
- backend
- rendered workload path
- summary JSON path

## Success Criteria

This slice is complete when:

1. a render helper can materialize both legacy sample fixtures into runnable XML
2. a manual `Legacy Live Compare` workflow exists
3. it can run one selected legacy fixture against a real endpoint using repository secrets
4. it uploads rendered workload, JSON summary, and log artifacts
5. default CI remains unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
