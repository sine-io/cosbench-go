# Legacy Comparison Design

## Goal

Create a repeatable way to compare `cosbench-go` behavior against `cosbench-sineio` on the repository's high-value workload subset without turning this slice into a full benchmark program.

The output of this slice should be:

- a comparison checklist that names the specific fixtures to compare
- a consistent set of comparison dimensions
- a repository-local place to record findings and known deltas
- a thin workflow that contributors can run again later

## Why This Slice

The repository already has:

- active-subset XML fixtures
- local execution and export paths
- opt-in live smoke tests for S3/SIO connectivity

What it still lacks is a concrete, reusable process for answering a simple migration question:

> "On which representative workloads does `cosbench-go` currently match, differ, or remain unverified relative to `cosbench-sineio`?"

Right now that knowledge is implicit. This slice makes it explicit and repeatable.

## Scope

### In Scope

- a comparison matrix/checklist document under `docs/`
- a minimal runbook for collecting comparison evidence
- one place to record current status for representative fixtures
- light command helpers if they reduce repetition

### Out of Scope

- full historical parity across every COSBench workload
- always-on integration tests against legacy Java code
- large-scale performance benchmarking
- automatic orchestration of both systems end-to-end
- remote worker comparisons

## Recommended Approach

Treat this slice as **comparison infrastructure**, not as benchmark feature work.

That means:

- define the representative workload set
- define what counts as “same enough” for this phase
- document how to run both sides
- record findings in a stable table

Do not attempt to automate the entire legacy environment in this slice.

## Representative Comparison Set

Start with the fixtures already curated in this repository:

- `testdata/legacy/s3-config-sample.xml`
- `testdata/legacy/sio-config-sample.xml`
- `testdata/workloads/s3-active-subset.xml`
- `testdata/workloads/sio-multipart-subset.xml`
- `testdata/workloads/xml-inheritance-subset.xml`
- `testdata/workloads/xml-attribute-subset.xml`
- `testdata/workloads/xml-special-ops-subset.xml`

Not every fixture needs to be runnable on both systems immediately. The comparison document should explicitly mark:

- runnable now
- runnable with live endpoint setup
- parser-only comparison
- deferred

## Comparison Dimensions

For each representative fixture, capture only the highest-value dimensions:

- XML parse outcome
- normalized workload shape
- accepted storage/backend path
- execution outcome category
  - succeeded
  - succeeded with operation errors
  - failed preflight
  - failed during execution
- result surface availability
  - CLI summary
  - JSON export
  - CSV export
- notable semantic differences

For this phase, “match” does **not** require byte-identical report output.
It requires behaviorally consistent execution and comparable surfaced results.

## Deliverables

### 1. Comparison Document

Create a document such as `docs/legacy-comparison-matrix.md` that includes:

- workload / fixture name
- legacy reference status
- `cosbench-go` status
- comparison result
  - match
  - acceptable delta
  - mismatch
  - not yet run
- notes / follow-up

### 2. Runbook Section

Include short command recipes for:

- running `cosbench-go` locally against a fixture
- running opt-in smoke coverage
- locating the legacy sample/config references in `../cosbench-sineio`

If a thin `make` helper improves this, it should stay minimal and only wrap existing commands.

### 3. Seed Findings

Populate the matrix with current known information rather than leaving it empty.

It is acceptable for some rows to say “not yet run on live endpoint” as long as the reason is explicit.

## Non-Goals

This slice must not devolve into:

- building a benchmark harness for legacy Java
- inventing synthetic pass/fail thresholds for throughput
- trying to prove absolute performance equivalence
- automating unavailable environments

The goal is clarity and repeatability, not maximum automation.

## Success Criteria

This slice is complete when:

1. a repository-local legacy comparison matrix exists
2. the matrix covers the current representative fixture set
3. comparison dimensions are explicit and reused across rows
4. contributors have a short runbook for reproducing the comparison steps
5. current known deltas are recorded instead of living only in conversation

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
