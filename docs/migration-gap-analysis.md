# Migration Gap Analysis

This document tracks the remaining gaps against the **local-only v1 migration boundary** for `cosbench-go`.

## Already Landed

- workload XML upload, parsing, and normalization
- endpoint persistence and reuse
- local job creation / start lifecycle
- file snapshots for jobs, results, events, and endpoints
- browser-visible dashboard, history, and job detail pages
- controller-facing matrix, config, advanced-config, stage-detail, timeline, Prometheus, and artifact surfaces
- unified driver-facing pages and APIs for overview, missions, workers, and logs
- S3/SIO-compatible configuration path and adapter wiring
- local execution for the active operation subset used by repository fixtures
- JSON and CSV export
- stage-level and operation-level result summaries
- remote driver registration, heartbeat, scheduling, mission claim, mission reporting, and combined-mode loopback execution
- local multi-process MinIO smoke for one controller-only and two driver-only processes
- local and remote happy-path smoke parity across both `s3` and `sio`
- local recovery smoke parity across both `s3` and `sio`
- dedicated manual GitHub workflows for remote happy-path and recovery smoke
- non-blocking GitHub workflows for remote happy-path and recovery matrices, including aggregate summaries and aggregate artifacts
- structured remote smoke and remote recovery summary artifacts, plus aggregate result reporting that `smoke-ready` now consumes directly
- structured `Smoke S3` and `Smoke S3 Matrix` summary artifacts for real-endpoint evidence, plus aggregate row status reporting beyond simple artifact presence
- a manual `Legacy Live Compare` workflow for rendering representative legacy XML against a real S3/SIO-compatible endpoint, with clean `skipped` behavior when repository live secrets are absent
- a manual `Legacy Live Compare Matrix` workflow for the representative `s3` and `sio` legacy samples, with per-row artifacts plus an aggregate summary/artifact and the same clean `skipped` behavior when repository live secrets are absent
- `smoke-ready` / `smoke-ready-json` reporting for local readiness, real-endpoint smoke readiness, real-endpoint smoke result states, legacy live compare readiness, legacy live compare matrix readiness, remote happy-path readiness, remote recovery readiness, the latest known workflow run state across those surfaces, and legacy live result states such as `executed` versus `skipped`, with the real-endpoint path now consuming structured smoke summaries first

## Closed for the Local v1 Boundary

- work-level summaries are persisted and exposed through job detail, JSON export, and CSV export
- `mfilewrite` uses real local file input for multipart upload
- `delay` performs real waiting behavior
- obvious adapter/config/file errors are rejected during start-time preflight

## Remaining Work Outside or Beyond the Local v1 Boundary

### Real-endpoint shakeout still remains

- S3/SIO paths are implemented, but broader live-environment validation is still pending
- auth, retry, and storage-specific edge behavior is not fully characterized against real systems
- `Smoke S3` now remains a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-28 (`23695743149`) produced a structured summary artifact with `result=skipped` because the live smoke tests were skipped for missing `COSBENCH_SMOKE_*` repository secrets; that proves workflow ergonomics and summary generation, not endpoint parity
- `Smoke S3 Matrix` now also remains a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-28 (`23695743153`) produced an aggregate summary with both `s3` and `sio` rows marked `skipped` for the same missing-secret reason; that proves matrix ergonomics, structured aggregation, and summary consumption, not endpoint parity
- `Legacy Live Compare` now has a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-28 (`23696226320`) produced a normalized `result.json` with `result=skipped` because `COSBENCH_SMOKE_*` repository secrets were not configured; that proves workflow ergonomics and summary generation, not endpoint parity
- `Legacy Live Compare Matrix` now also has a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-28 (`23696226307`) produced an aggregate summary with both `s3` and `sio` rows marked `skipped`; that proves matrix ergonomics, structured aggregation, and summary consumption, not endpoint parity
- the current comparison checklist, runbook, and seed findings now live in `docs/legacy-comparison-matrix.md`
- the latest repository-hosted `Remote Smoke Matrix` run on 2026-03-28 (`23696657083`) and `Remote Smoke Recovery Matrix` run on 2026-03-28 (`23696657085`) both completed successfully, confirming that the structured remote summary artifacts still align with `smoke-ready` consumption on the current mainline head; `smoke-ready` now also exposes those workflows as `remote_happy_latest_source` / `remote_recovery_latest_source` plus direct `remote_happy_latest_url` / `remote_recovery_latest_url` fields

## Deferred By Design

- non-S3 drivers
- DB-backed persistence
- full historical COSBench parity
- legacy UI/chart parity

## Main Remaining Risks

1. Real-world XML variance may exceed the active fixture subset.
2. The remote split is now real and backed by multi-process happy-path and recovery smoke, but broader production-grade scheduling and durability behavior still needs deeper validation.
3. SineIO-specific behavior may still diverge under larger or real-endpoint workloads even though the local and remote smoke paths now cover `sio`.
4. Behavior that looks correct under local mock runs may still diverge under real endpoint pressure.

The migration is considered closed for this phase because the local-only gaps are resolved and the deferred items remain explicitly out of scope.
