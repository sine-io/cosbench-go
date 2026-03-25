# cosbench-go

Go re-implementation of COSBench with behavioral compatibility focused on the active `cosbench-sineio` workload subset.

## Current Scope
- Local-only v1 closure: single-process control plane plus executor
- XML workload compatibility for the active S3/SIO subset tracked in `testdata/`
- Storage focus: **S3** and **SIO**
- Web inspection, snapshot persistence, JSON export, and CSV export
- Remote controller/driver split kept as a future seam, not a blocker for this phase

## Current Status
- XML parsing, normalization, storage config, local execution, snapshots, and web flows are landed
- The local-only v1 closure now includes stronger start-time preflight checks, real file-backed multipart behavior, and work-level reporting in exports and job detail views
- Remote controller/driver split remains intentionally deferred

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

## Smoke Tests
- Live smoke coverage is opt-in and does not run by default in `go test ./...`
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

## CI
- Repository CI runs `make validate` on `push` and `pull_request`
- The default CI path does not run `make smoke-s3`; live endpoint checks remain opt-in
- A manual GitHub Actions workflow can run `make compare-local` on demand
- The manual `compare-local` workflow uploads `.artifacts/compare-local/` as a downloadable artifact
- That manual workflow also writes a GitHub job summary from `.artifacts/compare-local/index.json`

## Legacy Comparison
- The current comparison checklist and runbook live in `docs/legacy-comparison-matrix.md`
- Code-level S3/SIO delta notes live in `docs/storage-driver-comparison-notes.md`
- Live endpoint prerequisites and execution order live in `docs/legacy-live-run-checklist.md`
- Use it to track which representative fixtures are parser-only, runnable now, runnable with live endpoint setup, or still unverified against `cosbench-sineio`
- Run `GO=$(which go || echo /snap/bin/go) make compare-local` to refresh the safe mock-backed local comparison set
- That command refreshes `.artifacts/compare-local/` in place and rewrites its `*.json` results for the curated fixture set
- The curated fixture list for that command lives in `testdata/workloads/compare-local-fixtures.txt`
- `.artifacts/compare-local/index.json` is the top-level artifact index for those per-fixture summaries and their key metrics
- If you override `COMPARE_LOCAL_OUTPUT_DIR`, keep the basename as `compare-local`

## References
- Legacy project reference: `../cosbench-sineio`
- Representative XML inputs:
  - `testdata/legacy/s3-config-sample.xml`
  - `testdata/legacy/sio-config-sample.xml`
  - `testdata/workloads/s3-active-subset.xml`
  - `testdata/workloads/sio-multipart-subset.xml`
