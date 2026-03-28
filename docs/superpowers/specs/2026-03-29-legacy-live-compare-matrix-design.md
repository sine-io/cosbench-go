# Legacy Live Compare Matrix Design

## Goal

Add a manual `Legacy Live Compare Matrix` workflow that runs the two representative legacy fixtures through their matching backends:

- `testdata/legacy/s3-config-sample.xml` on `s3`
- `testdata/legacy/sio-config-sample.xml` on `sio`

and emits both per-row artifacts and one aggregate summary/artifact.

## Why This Matters

The repository now has a single-run `Legacy Live Compare` workflow, but operators still have to trigger it twice and manually compare the results if they want coverage across both representative legacy sample families.

The current remaining migration risk is broader real-endpoint shakeout. A small matrix workflow improves repeatability and operator visibility without changing any runtime code.

## Scope

In scope:

- one manual GitHub Actions matrix workflow
- one small aggregation script
- workflow-shape tests
- README usage note

Out of scope:

- changing `Legacy Live Compare` itself
- changing XML rendering or execution behavior
- adding this workflow to default `CI`

## Design

Create `.github/workflows/legacy-live-compare-matrix.yml` with:

- `workflow_dispatch`
- optional `region` input
- optional `path_style` input
- matrix rows:
  - `backend=s3`, `fixture=testdata/legacy/s3-config-sample.xml`
  - `backend=sio`, `fixture=testdata/legacy/sio-config-sample.xml`

Each row should inline the same preflight/render/run pattern already used by `Legacy Live Compare`, so missing secrets still produce a clean `skipped` artifact instead of a failed run.

Each row should upload a per-row artifact:

- `legacy-live-compare-s3`
- `legacy-live-compare-sio`

Add an `aggregate` job that downloads those artifacts, reads their `summary.json` files, and emits:

- `.artifacts/legacy-live-compare-matrix-aggregate/summary.json`
- `.artifacts/legacy-live-compare-matrix-aggregate/summary.md`

plus a GitHub job summary and one aggregate artifact.

## Acceptance Criteria

- the workflow exists as a manual matrix entrypoint
- both representative legacy fixtures are covered
- each row uploads its own artifact
- the aggregate job uploads a combined artifact and summary
- missing live secrets still yield a successful workflow with `skipped` row evidence rather than row failure
