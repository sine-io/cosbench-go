# Migration Gap Analysis

This document tracks the remaining gaps against the **local-only v1 migration boundary** for `cosbench-go`.

## Already Landed

- workload XML upload, parsing, and normalization
- endpoint persistence and reuse
- local job creation / start lifecycle
- file snapshots for jobs, results, events, and endpoints
- browser-visible dashboard, history, and job detail pages
- S3/SIO-compatible configuration path and adapter wiring
- local execution for the active operation subset used by repository fixtures
- JSON and CSV export
- stage-level and operation-level result summaries

## Closed for the Local v1 Boundary

- work-level summaries are persisted and exposed through job detail, JSON export, and CSV export
- `mfilewrite` uses real local file input for multipart upload
- `delay` performs real waiting behavior
- obvious adapter/config/file errors are rejected during start-time preflight

## Remaining Work Outside or Beyond the Local v1 Boundary

### Real-endpoint shakeout still remains

- S3/SIO paths are implemented, but broader live-environment validation is still pending
- auth, retry, and storage-specific edge behavior is not fully characterized against real systems
- the current comparison checklist, runbook, and seed findings now live in `docs/legacy-comparison-matrix.md`

## Deferred By Design

- remote controller/worker split
- worker registration, heartbeats, and mission dispatch
- non-S3 drivers
- DB-backed persistence
- full historical COSBench parity
- legacy UI/chart parity

## Main Remaining Risks

1. Real-world XML variance may exceed the active fixture subset.
2. SineIO-specific behavior may still diverge under larger or real-endpoint workloads.
3. Behavior that looks correct under local mock runs may still diverge under real endpoint pressure.

The migration is considered closed for this phase because the local-only gaps are resolved and the deferred items remain explicitly out of scope.
