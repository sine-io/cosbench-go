# COSBench-Go Kanban

Last updated: 2026-03-25
Owner: Ross

## Doing
- Live-endpoint behavior comparison against legacy workloads
- Storage-driver comparison notes against legacy Java behavior
- Release/branch hygiene once the current migration slice is stable
- Structured compare-local result outputs for local and manual workflow reuse
- Compare-local fixture manifest cleanup
- Compare-local output pruning cleanup
- Compare-local output safety cleanup
- Compare-local artifact index
- Compare-local workflow summary
- Compare-local artifact metrics
- Compare-local fixture filter

## Next
- Local workflow polish only if comparison ergonomics need more than `make compare-local`
- Environment provisioning for actual legacy live checks
- Broader smoke/release automation only if the repository actually needs it

## Done
- Migration strategy chosen
- Local-only v1 scope defined
- Migration spec written
- Workload domain model created
- Config inheritance implemented
- Special work normalization implemented
- Validation rules implemented
- XML workload parser implemented
- KV config parser implemented
- Legacy S3/SIO sample testdata added
- Parser tests added and passing
- Storage adapter port defined
- S3 adapter skeleton created
- SIO adapter skeleton created
- Storage factory created
- Execution engine skeleton created
- Real AWS SDK v2 S3/SIO wiring landed
- Web control plane and snapshot persistence landed
- JSON/CSV result export landed
- Migration docs aligned to the local-only v1 boundary
- Work-level result summaries landed in snapshots, exports, and job detail pages
- Start-time preflight validation landed
- Real local `mfilewrite` / `delay` semantics landed
- Sequential cleanup / list scanning landed
- CI-friendly `make build` / `make validate` targets landed
- Normalization-focused unit tests landed
- Storage-adapter focused request-level tests landed
- Multipart upload now preserves `storage_class`
- Opt-in S3/SIO smoke-test workflow landed
- High-value XML fixture coverage landed for inheritance, attributes, and special-op shapes
- Legacy comparison matrix and runbook landed
- Mock-override evidence captured for legacy S3/SIO sample XML
- Storage-driver comparison notes landed from legacy Java code review
- Storage-level `part_size` / `restore_days` fallback landed for execution parity
- Local cancel flow landed for running jobs
- Restart/recovery polish landed for cancelling jobs
- Stage-aware mock realism landed for local runs
- Representative edge XML fixtures landed for delay-stage, splitrw, and reuse-data shapes
- Parser-facing coverage landed for deferred compatibility aliases and range/prefetch config shapes
- Parser-tolerated coverage landed for deferred auth-bearing XML shapes
- Minimal CI workflow landed for `make validate`
- Manual compare-local workflow landed
- Compare-local workflow artifact upload landed
- Structured compare-local JSON output directory landed
- Compare-local fixture manifest landed
- Compare-local output pruning landed
- Compare-local output safety guard landed
- Compare-local artifact index landed
- Compare-local workflow summary landed
- Compare-local artifact metrics landed
- Compare-local fixture filter landed
- Local CLI ergonomics landed (`-f`, positional path, pure JSON stdout)
- Local comparison command landed for curated mock-backed fixtures
- Legacy live-run checklist landed
- Representative S3/SIO workload fixtures added
- `go vet` passing
- Current unit tests passing
- In-repo checklist / board established

## Deferred
- Remote controller/driver HTTP split
- Legacy web UI parity
- Legacy-style charts
- Non-S3/SIO storages
