# cosbench-go

Go re-implementation of COSBench with behavioral compatibility focused on the active `cosbench-sineio` workload subset.

## Current Scope
- Local-only v1 closure: single-process control plane plus executor
- XML workload compatibility for the active S3/SIO subset tracked in `testdata/`
- Storage focus: **S3** and **SIO**
- Web inspection, snapshot persistence, JSON export, and CSV export
- Remote controller/driver split is now available as an early protocol skeleton with `controller-only`, `driver-only`, and `combined` runtime modes

## Current Status
- XML parsing, normalization, storage config, local execution, snapshots, and web flows are landed
- The local-only v1 closure now includes stronger start-time preflight checks, modeled auth inheritance, real file-backed `filewrite` / `mfilewrite` behavior, prefetch/range read request shaping, and work-level reporting in exports and job detail views
- The controller-facing closure now includes matrix, config, advanced-config, stage-detail, timeline, timeline CSV, Prometheus, and controller artifact endpoints under the unified Go service
- The remote split now includes persisted driver/mission state, registration/heartbeat/claim endpoints, a driver agent, combined-mode loopback execution, and a shared bearer token on driver write endpoints
- The unified service now also includes driver-facing overview, missions, mission detail, workers, and logs pages under `/driver/...`
- The repository now also includes a local `smoke-remote-local` path that runs one controller-only process, two driver-only processes, and local MinIO to validate remote work-unit execution end-to-end

## Local CLI
- Run with any of these equivalent workload forms:
  - `go run ./cmd/cosbench-go -workload testdata/workloads/s3-active-subset.xml -backend mock`
  - `go run ./cmd/cosbench-go -f testdata/workloads/s3-active-subset.xml -backend mock`
  - `go run ./cmd/cosbench-go testdata/workloads/s3-active-subset.xml -backend mock`
- `-json` now emits machine-readable JSON to stdout without progress text mixed in
- `-quiet` suppresses progress output entirely; `make compare-local` now uses it
- `-summary-file <path>` writes the same summary JSON to a file for later reuse

## Technical Direction
- Go-native rewrite, not line-by-line Java translation
- Clean boundaries across `internal/controlplane`, `internal/executor`, `internal/driver`, `internal/reporting`, and `internal/snapshot`
- Core dependencies are the Go standard library and AWS SDK v2 for S3/SIO-compatible access
- Remote driver write endpoints use `Authorization: Bearer <token>` with `COSBENCH_DRIVER_SHARED_TOKEN` as the current shared-token source

## Smoke Tests
- Live smoke coverage is opt-in and does not run by default in `go test ./...`
- Run `make --no-print-directory smoke-local` to start a temporary local MinIO endpoint and verify both the S3 smoke path and the SIO multipart smoke path end-to-end
- Run `make --no-print-directory smoke-remote-local` to validate the remote controller/driver split against local MinIO with one controller-only and two driver-only processes
- Run `make --no-print-directory smoke-ready` for a human-readable local/repo readiness summary
- Run `make --no-print-directory smoke-ready-json` for the same readiness view as JSON
- Run `GO=$(which go || echo /snap/bin/go) make smoke-s3`
- Required env:
  - `COSBENCH_SMOKE_ENDPOINT`
  - `COSBENCH_SMOKE_ACCESS_KEY`
  - `COSBENCH_SMOKE_SECRET_KEY`
- Optional env:
  - `COSBENCH_SMOKE_BACKEND` (`s3` by default, `sio` enables multipart smoke coverage)
  - `COSBENCH_SMOKE_REGION`
  - `COSBENCH_SMOKE_PATH_STYLE`
- `COSBENCH_SMOKE_BUCKET_PREFIX`
- If required env vars are missing, the smoke tests skip cleanly
- `smoke-remote-local` writes controller, driver, and MinIO artifacts under `.artifacts/remote-smoke/`
- Real `make smoke-s3` is now a local or private-network-only path. GitHub-hosted runners do not execute it because the repository does not have a public S3-compatible endpoint.
- To trigger the manual GitHub smoke workflow with GitHub CLI:
  - `gh workflow run "Smoke Local" --repo sine-io/cosbench-go`

## CI
- Repository CI runs `make validate` on `push` and `pull_request`
- The default CI path does not run `make smoke-s3`; live endpoint checks remain opt-in
- A manual GitHub Actions workflow can run `make compare-local` on demand
- A separate manual GitHub Actions workflow can run `make smoke-local` without external secrets to verify the local live-endpoint path on GitHub-hosted runners
- A manual GitHub Actions workflow can run the remote multi-process MinIO smoke path on demand
- GitHub-hosted runners no longer attempt real `make smoke-s3`; that path remains for local or private-network execution only
- The manual `compare-local` workflow uploads `.artifacts/compare-local/` as a downloadable artifact
- That manual workflow also writes a GitHub job summary from `.artifacts/compare-local/index.json`
- A manual GitHub Actions workflow can run `make smoke-remote-local` and upload `.artifacts/remote-smoke/`

To trigger the manual remote smoke workflow with GitHub CLI:

```bash
gh workflow run "Remote Smoke Local" --repo sine-io/cosbench-go
```

## Legacy Comparison
- The current comparison checklist and runbook live in `docs/legacy-comparison-matrix.md`
- Code-level S3/SIO delta notes live in `docs/storage-driver-comparison-notes.md`
- Live endpoint prerequisites and execution order live in `docs/legacy-live-run-checklist.md`
- Use it to track which representative fixtures are parser-only, runnable now, runnable with live endpoint setup, or still unverified against `cosbench-sineio`
- Run `GO=$(which go || echo /snap/bin/go) make compare-local` to refresh the safe mock-backed local comparison set
- Run `GO=$(which go || echo /snap/bin/go) make compare-local-list` to print the valid curated fixture names
- Run `GO=$(which go || echo /snap/bin/go) make compare-local-list-json` to print the curated fixture names and workload paths as JSON
- Run `make --no-print-directory worktree-audit` to list local worktrees and their status relative to `origin/main`, with generation-time, base-ref, and current-worktree header metadata
- Run `make --no-print-directory worktree-audit-json` to get the same worktree audit data as JSON, including a consistent top-level `meta` object and a `views.audit` envelope
- The JSON audit output now includes structured `ahead` / `behind` counts for each row
- The JSON audit output now also includes a top-level `summary` section with total and per-state counts
- That summary also includes `stale` and `prune_candidates` counts
- The JSON audit output now also includes a `current` flag for the current worktree row
- Branches whose patches are already present in the base ref via squash merge are now classified as `integrated`
- Audit outputs now sort merged rows first, then active rows by descending `behind` count
- Set `WORKTREE_AUDIT_BASE_REF=<ref>` when you want the audit and prune helpers to compare against something other than `origin/main`
- Run `make --no-print-directory worktree-audit-merged` to list only worktrees already merged into `origin/main`, with generation-time, base-ref, and current-worktree header metadata
- Run `make --no-print-directory worktree-audit-merged-json` to get the merged-only worktree audit data as JSON, including a consistent top-level `meta` object and a `views.audit` envelope
- Run `make --no-print-directory worktree-audit-integrated` to list only worktrees whose patches are already present in the base ref via squash merge or equivalent history, with generation-time, base-ref, and current-worktree header metadata
- Run `make --no-print-directory worktree-audit-integrated-json` to get the integrated-only worktree audit data as JSON, including a consistent top-level `meta` object and a `views.audit` envelope
- Run `make --no-print-directory worktree-audit-prune` to list only worktrees that are prune candidates under the current cleanup rules, with generation-time, base-ref, and current-worktree header metadata
- Run `make --no-print-directory worktree-audit-prune-json` to get the prune-candidates view as JSON, including a consistent top-level `meta` object and a `views.audit` envelope
- Run `make --no-print-directory worktree-audit-stale` to list only active worktrees that are behind `origin/main`, with generation-time, base-ref, and current-worktree header metadata
- Run `make --no-print-directory worktree-prune-plan` to print suggested cleanup commands for merged worktrees without executing them, along with generation-time, base-ref, and current-worktree header metadata
- Run `make --no-print-directory worktree-prune-plan-json` to get the same non-destructive cleanup plan as structured JSON with a consistent top-level `meta` object, plus `views.prune_plan`, `summary`, and `rows`
- Run `make --no-print-directory worktree-cleanup-report` to generate a single Markdown report combining the audit, integrated view, stale view, prune-candidates view, and prune plan, with summary counts plus generation-time and current-worktree metadata
- Run `make --no-print-directory worktree-cleanup-report-json` to get the same combined cleanup report in machine-readable form, with a consistent top-level `meta` object and a preferred `views` container for all section payloads
- Run `GO=$(which go || echo /snap/bin/go) make compare-local COMPARE_LOCAL_FILTER=mock-stage-aware` to refresh only one curated fixture
- Run `GO=$(which go || echo /snap/bin/go) make compare-local COMPARE_LOCAL_FILTER=mock-stage-aware,xml-splitrw-subset` to refresh a curated subset
- The list targets also respect `COMPARE_LOCAL_FILTER`, so you can preview the selected subset before running it
- That command refreshes `.artifacts/compare-local/` in place and rewrites its `*.json` results for the curated fixture set
- The curated fixture list for that command lives in `testdata/workloads/compare-local-fixtures.txt`
- `.artifacts/compare-local/index.json` is the top-level artifact index for those per-fixture summaries and their key metrics
- `.artifacts/compare-local/summary.md` is the local human-readable summary, and the manual workflow reuses that same file
- Filtered runs only accept fixture names from `testdata/workloads/compare-local-fixtures.txt`
- If you override `COMPARE_LOCAL_OUTPUT_DIR`, keep the basename as `compare-local`

## References
- Legacy project reference: `../cosbench-sineio`
- Representative XML inputs:
  - `testdata/legacy/s3-config-sample.xml`
  - `testdata/legacy/sio-config-sample.xml`
  - `testdata/workloads/s3-active-subset.xml`
  - `testdata/workloads/sio-multipart-subset.xml`
