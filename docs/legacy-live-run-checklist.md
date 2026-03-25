# Legacy Live Run Checklist

This document is the step-by-step execution guide for future live comparisons against `cosbench-sineio`.

Use [legacy-comparison-matrix.md](/root/.openclaw/workspace/projects/cosbench-go/.worktrees/migration-closure/docs/legacy-comparison-matrix.md) as the system of record for findings.
Use this checklist to decide what to run, in what order, and what to record once a real endpoint is available.

## 1. Preconditions

Before claiming any live comparison result:

- ensure you have a reachable S3-compatible or SIO-compatible endpoint
- export:
  - `COSBENCH_SMOKE_ENDPOINT`
  - `COSBENCH_SMOKE_ACCESS_KEY`
  - `COSBENCH_SMOKE_SECRET_KEY`
- optionally export:
  - `COSBENCH_SMOKE_BACKEND`
  - `COSBENCH_SMOKE_REGION`
  - `COSBENCH_SMOKE_PATH_STYLE`
  - `COSBENCH_SMOKE_BUCKET_PREFIX`

If the environment is not available, keep matrix rows in their current pending/live-unverified state.

## 2. Smoke Precheck

First confirm that the current Go adapter path can talk to the target endpoint:

```bash
GO=$(which go || echo /snap/bin/go) make smoke-s3
```

Record:

- whether smoke passed or failed
- which backend was used
- any endpoint-specific setup quirks

If this fails, stop and fix credentials/connectivity before workload-level comparison.

## 3. Recommended Run Order

Run these in order:

1. `testdata/legacy/sio-config-sample.xml`
   Reason: strongest current candidate for `mprepare` + `mwrite`

2. `testdata/legacy/s3-config-sample.xml`
   Reason: strongest current S3 delta candidate because mock evidence already showed a high read/write error count

3. storage-level `part_size` / `restore_days`
   Reason: code parity is now in place and should be checked against a real endpoint

4. cleanup/list-sensitive scenarios
   Reason: delete tolerance and list-shape differences remain likely watchpoints

## 4. Recording Rules

After each live run, update the corresponding row in `docs/legacy-comparison-matrix.md` with:

- `Legacy Reference Status`
- `cosbench-go Status`
- `Result`
  - `match`
  - `acceptable delta`
  - `mismatch`
- `Notes`
  - exact outcome category:
    - succeeded
    - succeeded with operation errors
    - failed preflight
    - failed during execution
  - whether CLI / JSON / CSV outputs looked correct
  - any unexpected semantic differences

Use concrete dates, fixture names, and short factual notes.

## 5. Known Watchpoints

Review [storage-driver-comparison-notes.md](/root/.openclaw/workspace/projects/cosbench-go/.worktrees/migration-closure/docs/storage-driver-comparison-notes.md) before live runs.

Highest-value watchpoints:

- SIO `path_style_access` default handling
- delete tolerance for missing buckets/objects
- list output shape and downstream assumptions
- storage-level `part_size`
- storage-level `restore_days`
- slash-containing SIO bucket/container names

These are not guaranteed mismatches, but they are the most likely places for live differences to surface.

## 6. Environment Blocker Rule

If no live endpoint or credentials are available:

- do not invent live conclusions
- do not reclassify rows as `match`
- keep local mock evidence and parser evidence current

The process can be considered ready even when the environment is not.
