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
- controller-side proactive lease sweep and automatic mission requeue for expired remote leases
- controller-side error-event emission for expired remote mission leases
- controller-side proactive stale-driver health sweep for remote execution
- controller-side error-event emission for claimed or running jobs affected by stale-driver heartbeat timeout
- local multi-process MinIO smoke for one controller-only and two driver-only processes
- local and remote happy-path smoke parity across both `s3` and `sio`
- local recovery smoke parity across both `s3` and `sio`
- local recovery smoke now also validates the controller-side `mission lease expired` and stale-driver heartbeat-timeout events, not just reassignment success
- dedicated manual GitHub workflows for remote happy-path and recovery smoke
- non-blocking GitHub workflows for remote happy-path and recovery matrices, including aggregate summaries and aggregate artifacts
- structured remote smoke and remote recovery summary artifacts, plus aggregate result reporting that `smoke-ready` now consumes directly
- structured `Smoke S3` and `Smoke S3 Matrix` summary artifacts for real-endpoint evidence, plus aggregate row status reporting beyond simple artifact presence
- a manual `Legacy Live Compare` workflow for rendering representative legacy XML against a real S3/SIO-compatible endpoint, with clean `skipped` behavior when repository live secrets are absent
- a manual `Legacy Live Compare Matrix` workflow for the representative `s3` and `sio` legacy samples, with per-row artifacts plus an aggregate summary/artifact and the same clean `skipped` behavior when repository live secrets are absent
- `smoke-ready` / `smoke-ready-json` reporting for local readiness, real-endpoint smoke readiness, real-endpoint smoke result states, legacy live compare readiness, legacy live compare matrix readiness, remote happy-path readiness, remote recovery readiness, the latest known workflow run state across those surfaces, and legacy live result states such as `executed` versus `skipped`, with the real-endpoint path now consuming structured smoke summaries first
- `smoke-ready` now also reports `Smoke Ready Validate` as a separate contract-surface signal, including its latest validation result plus direct run id, URL, artifact, duration, and timestamp fields
- `smoke-ready` now also exposes the latest trigger event for each evidence surface, so repository-hosted evidence can distinguish manual versus scheduled refreshes
- `smoke-ready` now also exposes the latest run id for each evidence surface, so machine consumers can join evidence back to GitHub runs without parsing URLs
- `smoke-ready` now also exposes the latest duration in seconds for each evidence surface, so machine consumers can reason about workflow timing without scraping GitHub
- `smoke-ready` now also exposes the latest age in seconds for each evidence surface, so machine consumers can detect evidence staleness directly instead of re-deriving it from timestamps
- `smoke-ready` now also exposes `*_latest_fresh` booleans and a shared `freshness_thresholds_seconds` block, so machine consumers can apply the repository's freshness policy directly instead of re-encoding it
- `smoke-ready` now also exposes `*_current` booleans, so machine consumers can ask the simpler question “is the latest evidence valid for this checkout?” without recombining success, freshness, and head-alignment themselves
- `smoke-ready` now also exposes `*_current_reason`, so machine consumers and operators can tell whether a non-current signal comes from `not_successful`, `stale`, `head_mismatch`, or a currently valid record
- `smoke-ready` now also exposes family-level `*_current_ready` booleans, so callers can answer “is current contract / real-endpoint / legacy-live / remote evidence available for this checkout?” without manually folding the lower-level per-surface fields
- `smoke-ready` now also exposes family-level `*_current_ready_reason` strings, so callers can tell whether those aggregate non-current states come from `not_successful`, `stale`, `head_mismatch`, or `mixed`
- the human-readable `smoke-ready` output now also prints those freshness/current signals and the threshold block directly, so operators do not need to switch to JSON mode just to inspect them
- `smoke-ready` now also exposes the latest head SHA for each evidence surface, so machine consumers can tell exactly which commit produced the latest evidence
- `smoke-ready` now also exposes the latest head branch for each evidence surface, so machine consumers can distinguish `main` evidence from non-main evidence without extra GitHub queries
- `smoke-ready` now also exposes whether each latest-evidence surface matches the current checkout HEAD, so machine consumers can detect stale evidence without re-implementing SHA comparison logic
- `smoke-ready` now also exposes the current checkout branch alongside `current_head_sha`, so machine consumers can read the active local ref without shelling out separately
- the latest repository-hosted `Smoke Ready Validate` run on 2026-03-30 (`23721782673`) completed successfully, uploaded both `smoke-ready-validate-output` and `smoke-ready-validate-summary`, and confirms that `smoke-ready` now consumes normalized schema-validation evidence directly while also emitting `*_latest_age_seconds` in the current summary payload so freshness is machine-readable without recomputing from timestamps
- `Smoke Ready Validate` now also finalizes its own captured `smoke-ready.json` with the current workflow run metadata before upload, so `schema_validation_current=true` and `schema_validation_current_reason=current` describe the run that produced the artifact instead of a previous completed validation run

## Closed for the Local v1 Boundary

- work-level summaries are persisted and exposed through job detail, JSON export, and CSV export
- `mfilewrite` uses real local file input for multipart upload
- `delay` performs real waiting behavior
- obvious adapter/config/file errors are rejected during start-time preflight

## Remaining Work Outside or Beyond the Local v1 Boundary

### Real-endpoint shakeout still remains

- S3/SIO paths are implemented, but broader live-environment validation is still pending
- auth, retry, and storage-specific edge behavior is not fully characterized against real systems
- `Smoke S3` now remains a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-30 (`23721784622`) produced a structured summary artifact with `result=skipped` because the live smoke tests were skipped for missing `COSBENCH_SMOKE_*` repository secrets; that proves workflow ergonomics and summary generation, not endpoint parity
- `Smoke S3 Matrix` now also remains a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-30 (`23721785194`) produced an aggregate summary with both `s3` and `sio` rows marked `skipped` for the same missing-secret reason; that proves matrix ergonomics, structured aggregation, and summary consumption, not endpoint parity
- `Legacy Live Compare` now has a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-30 (`23721785876`) produced a normalized `result.json` with `result=skipped` because `COSBENCH_SMOKE_*` repository secrets were not configured; that proves workflow ergonomics and summary generation, not endpoint parity
- `Legacy Live Compare Matrix` now also has a stable manual GitHub Actions entrypoint, and the latest repository-hosted run on 2026-03-30 (`23721786500`) produced an aggregate summary with both `s3` and `sio` rows marked `skipped`; that proves matrix ergonomics, structured aggregation, and summary consumption, not endpoint parity
- the current comparison checklist, runbook, and seed findings now live in `docs/legacy-comparison-matrix.md`
- the latest repository-hosted `Remote Smoke Matrix` run on 2026-03-30 (`23721783308`) and `Remote Smoke Recovery Matrix` run on 2026-03-30 (`23723374807`) both completed successfully, confirming that the structured remote summary artifacts still align with `smoke-ready` consumption on the current mainline head; `smoke-ready` now also exposes those workflows as `remote_happy_latest_source` / `remote_recovery_latest_source` plus direct `remote_happy_latest_url` / `remote_recovery_latest_url` fields

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
