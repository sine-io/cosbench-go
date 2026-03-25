# COSBench-Go Migration Board

Last updated: 2026-03-25
Owner: Ross
Status: In Progress

## Legend
- [ ] Not started
- [~] In progress
- [x] Done
- [-] Deferred

---

## 0. Project visibility
- [x] Create a persistent migration checklist / board
- [x] Track completed phases in-repo so sine can inspect progress anytime
- [~] Keep this board updated at each meaningful milestone

---

## 1. Foundation / bootstrap
- [x] Confirm Go toolchain availability in runtime
- [x] Initialize module `github.com/sine-io/cosbench-go`
- [x] Establish baseline verification (`go vet`, `go test ./...`)
- [x] Add CI-friendly Make targets for validate / test / build
- [x] Add repository-local CI automation for `make validate`
- [x] Add manual workflow automation for `make compare-local`
- [x] Add manual workflow automation for `make smoke-s3`
- [x] Add a GitHub job summary and explicit backend choice for the manual smoke workflow
- [x] Constrain the smoke workflow `path_style` input to explicit choices
- [x] Show required smoke secret presence in the manual smoke workflow summary
- [x] Fail the manual smoke workflow early when required secrets are missing
- [x] Upload the raw smoke workflow output as an artifact
- [x] Keep the smoke summary visible even when preflight fails
- [x] Upload `compare-local` workflow output as an artifact
- [x] Add structured per-fixture JSON outputs for `make compare-local`
- [x] Define a single manifest for the curated `compare-local` fixture set
- [x] Recreate the `compare-local` output directory before regenerating results
- [x] Guard `COMPARE_LOCAL_OUTPUT_DIR` so compare-local only refreshes dedicated directories
- [x] Add a top-level `index.json` for compare-local artifact discovery
- [x] Add a GitHub job summary for the manual compare-local workflow
- [x] Include per-fixture metrics in the compare-local index and workflow summary
- [x] Add a single-fixture filter for local and manual compare-local runs
- [x] Add a non-destructive `worktree-audit` helper for local cleanup planning
- [x] Fail fast when `COMPARE_LOCAL_FILTER` does not match a curated fixture
- [x] Add a machine-readable `worktree-audit-json` helper
- [x] Add structured `ahead` / `behind` fields to the JSON worktree audit output
- [x] Add top-level summary counts to the JSON worktree audit output
- [x] Add an explicit current-worktree marker to the JSON audit output
- [x] Add a configurable `WORKTREE_AUDIT_BASE_REF` override for audit and prune helpers
- [x] Sort worktree audit outputs by usefulness
- [x] Add a merged-only `worktree-audit-merged` helper
- [x] Add a machine-readable `worktree-audit-merged-json` helper
- [x] Add a stale-only `worktree-audit-stale` helper
- [x] Add a non-destructive `worktree-prune-plan` helper
- [x] Add a machine-readable `worktree-prune-plan-json` helper
- [x] Add a Markdown `worktree-cleanup-report` helper
- [x] Add a machine-readable `worktree-cleanup-report-json` helper
- [x] Add `make compare-local-list` for fixture-name discovery
- [x] Add a local `summary.md` artifact for compare-local outputs
- [x] Allow comma-separated compare-local fixture subsets
- [x] Add `make compare-local-list-json` for machine-readable fixture discovery
- [x] Make compare-local listing targets respect `COMPARE_LOCAL_FILTER`

---

## 2. Migration boundary
- [x] Inspect legacy `cosbench-sineio` module boundaries
- [x] Define migration strategy: behavior-compatible Go rewrite, not line-by-line translation
- [x] Define local-only v1 scope (XML + S3/SIO + execution + reporting + web control plane)
- [x] Write migration spec document (`docs/migration-spec-v1.md`)
- [x] Align all migration-facing docs to the same local-only boundary

---

## 3. Workload domain model
- [x] Define workload domain model
- [x] Define workflow / stage / work / operation / storage structs
- [x] Implement config inheritance chain
- [x] Implement special work normalization:
  - [x] `init`
  - [x] `prepare`
  - [x] `mprepare`
  - [x] `cleanup`
  - [x] `dispose`
  - [x] `delay`
- [x] Implement validation rules:
  - [x] required names / workers / storage
  - [x] runtime / totalOps / totalBytes limit validation
  - [x] op ratio sum = 100
  - [x] `workers <= totalOps` when `totalOps` is set
  - [x] SIO-only operation restrictions

---

## 4. Parsing layer
- [x] Implement XML -> domain mapping
- [x] Implement workload file parser
- [x] Implement storage KV parser (`k=v;k2=v2`)
- [x] Add legacy sample testdata for S3
- [x] Add legacy sample testdata for SIO
- [x] Add parser unit test for S3 sample
- [x] Add parser unit test for SIO sample
- [ ] Define and implement escaping rules beyond legacy baseline if truly needed
- [-] Model explicit XML `<auth>` nodes in the local-only closure

---

## 5. Storage abstraction
- [x] Define `StorageAdapter` port
- [x] Define storage metadata DTOs (`ObjectMeta`, `ObjectEntry`)
- [x] Implement adapter factory by storage type
- [x] Implement real AWS SDK v2 client wiring for S3
- [x] Implement real AWS SDK v2 client wiring for SIO
- [x] Support endpoint override, path-style access, proxy, `no_verify_ssl`, `aws_region`, `storage_class`, `restore_days`, and multipart `part_size`

---

## 6. Execution engine
- [x] Create driver execution engine skeleton
- [x] Add worker-pool concurrency structure
- [x] Add runtime-based stop condition
- [x] Add totalOps-based stop condition
- [x] Add sample collection
- [x] Implement weighted operation picker matching legacy ratio behavior
- [x] Implement config expression parsing for container / object / size patterns
- [x] Implement bucket / object naming generators
- [x] Implement `executeOp()` dispatch to storage adapter
- [x] Implement basic metrics aggregation
- [x] Honor storage-level `part_size` / `restore_days` as execution defaults with op-level override
- [x] Preserve stage-aware `mock` state across one local run/job
- [x] Implement real `mfilewrite` and `delay` semantics
- [x] Implement sequential scanning for cleanup / list flows
- [x] Implement cancel / abort path

---

## 7. Operation support
### Core COSBench ops
- [x] `init`
- [x] `prepare`
- [x] `write`
- [x] `read`
- [x] `delete`
- [x] `cleanup`
- [x] `dispose`
- [x] `list`
- [x] `delay`

### SineIO-specific ops
- [x] `mprepare`
- [x] `mwrite`
- [x] `mfilewrite`
- [x] `localwrite`
- [x] `head`
- [x] `restore`

---

## 8. Runnable entrypoints
- [x] Add CLI entrypoint (`cmd/cosbench-go`)
- [x] Print normalized workload summary before run
- [x] Run single-process / single-driver benchmark locally
- [x] Output JSON result summary
- [x] Output human-readable console summary
- [x] Improve local CLI ergonomics (`-f`, positional workload path, pure JSON stdout)
- [x] Add `make compare-local` for repeatable mock-backed comparison runs
- [x] Add CLI `-summary-file` support for reusable local comparison artifacts

---

## 9. Local control plane and reporting
- [x] Implement control-plane lifecycle and snapshots
- [x] Implement dashboard, history, endpoint, upload, and job detail pages
- [x] Define benchmark report model
- [x] Stage-level summary
- [x] Throughput / bandwidth metrics
- [x] Success / failure ratios
- [x] JSON export
- [x] CSV export
- [x] Work-level summary
- [x] Start-time preflight validation
- [x] Polish restart recovery for `cancelling` vs `interrupted` jobs

---

## 10. Deferred remote split
- [-] Define controller / worker HTTP transport
- [-] Define mission / workload / sample DTOs
- [-] Implement controller skeleton
- [-] Implement driver skeleton
- [-] Implement driver registration / heartbeat
- [-] Implement mission assignment
- [-] Implement sample upload / final report upload
- [-] Support multi-driver execution

---

## 11. Validation / compatibility
- [x] `go vet` clean on current code
- [x] Current parser tests passing
- [x] Add operation-picker tests
- [x] Add expression parser tests
- [x] Add local mock-storage integration tests
- [x] Add normalization-focused unit tests
- [x] Add high-value XML fixture coverage for inheritance, attributes, and special-op shapes
- [x] Add representative edge XML fixtures for delay-stage, splitrw, and reuse-data shapes
- [x] Add parser-facing coverage for deferred compatibility storage aliases and range/prefetch config shapes
- [x] Add parser-tolerated coverage for deferred auth-bearing XML shapes
- [x] Add storage adapter tests
- [x] Add real S3/SIO smoke-test workflow
- [~] Compare benchmark behavior against legacy workloads (matrix seeded; local `compare-local` evidence collected; live-run checklist documented; live environment still pending)
- [x] Add storage-driver comparison notes from legacy Java code review

---

## 12. Current closure slice
- [x] Write migration-closure design spec
- [x] Write migration-closure implementation plan
- [x] Converge migration docs
- [x] Implement work-level result visibility
- [x] Implement `mfilewrite` semantics
- [x] Implement real `delay` semantics
- [x] Add preflight validation
