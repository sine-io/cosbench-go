# Stage-Aware Mock Realism Design

## Goal

Make local `mock`-backed multi-stage workload runs behave more like a real storage target by preserving object state across stages and works within a single local run.

This slice is about local realism for smoke and representative runs. It is not a change to real S3/SIO behavior.

## Problem

Today both the Web control-plane path and the CLI local runner create a fresh `mock` adapter for each work.

That means:

- `prepare` writes do not naturally feed later `read` / `list` / `cleanup`
- multi-stage representative fixtures can produce fake local failures or unrealistic error counts
- `-backend mock` is less useful as a local stand-in for stage-ordered workloads

The problem is not the `mock` storage implementation itself. The problem is adapter lifetime.

## Scope

### In Scope

- reuse `mock` storage state across works and stages inside one local run
- apply that reuse to:
  - Web/control-plane job execution
  - CLI local runner
- add one representative stage-aware mock fixture and tests
- keep adapter cleanup bounded to the single run/job

### Out of Scope

- changing real S3/SIO adapter lifetime behavior
- introducing persistent mock state across multiple jobs
- changing storage port interfaces
- broad simulator behavior for every production edge case

## Recommended Approach

Implement per-run `mock` adapter reuse.

### Control plane

Inside one `runJob` execution:

- when the resolved backend is `mock`, reuse one adapter instance for the duration of the job
- completed stages should leave their created buckets/objects visible to later stages in the same job
- dispose the shared mock adapter once the job is finished, failed, or cancelled

### CLI

Inside one `cmd/cosbench-go` workload run:

- when the effective backend is `mock`, reuse one adapter instance across all works/stages in the workload
- dispose it once the CLI run finishes

This keeps the change local and avoids leaking state between independent runs.

## Why This Shape

This is the smallest change that fixes the realism problem:

- no new storage API
- no persistent datastore
- no special XML semantics
- no change to real backends

It only changes the lifetime of an existing in-memory backend under the local/mock path.

## Test Shape

Add a representative fixture such as `testdata/workloads/mock-stage-aware.xml` that proves stage-to-stage continuity:

1. `init`
2. `prepare`
3. `read`
4. `list`
5. `cleanup`
6. `dispose`

Expected behavior under local `mock`:

- the job succeeds
- `read` does not fail because objects vanished between stages
- `list` sees the prepared objects
- `cleanup` removes the prepared objects

## Success Criteria

This slice is complete when:

1. a multi-stage `mock` workload can pass through `prepare` → `read` / `list` / `cleanup` with shared state
2. Web/control-plane local runs preserve mock state for one job
3. CLI local runs preserve mock state for one workload invocation
4. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
