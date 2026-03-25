# PLAN.md — cosbench-go

## Execution Rules
- Build rebuild-first, not line-by-line port.
- Protect scope: S3/SineIO + XML + Web-first + local execution only for initial waves.
- Keep module seams clean so remote workers can be added later.
- Every task must specify touched files and verification.

---

## Wave 1 — Project Skeleton & Core Domain

### <task type="auto">
**Goal:** establish the core project skeleton and domain model.

**Touched files:**
- `cmd/server/main.go`
- `internal/domain/*.go`
- `internal/app/*.go`
- `go.mod`

**Requirements:**
- define domain entities for workload, stage, operation, job, metrics, endpoint config
- create a minimal server bootstrap
- ensure package boundaries reflect long-term architecture

**Verification:**
- `go build ./...`
- minimal startup command runs without panic
</task>

### <task type="auto">
**Goal:** define snapshot storage layout and contracts.

**Touched files:**
- `internal/snapshot/*.go`
- `docs/storage-layout.md`

**Requirements:**
- define file layout for jobs/results/events/endpoints
- implement interfaces for save/load primitives
- no DB dependency

**Verification:**
- unit tests for round-trip save/load
</task>

---

## Wave 2 — Workload XML Compatibility Core

### <task type="auto">
**Goal:** implement COSBench workload XML parser for the active subset.

**Touched files:**
- `internal/workloadxml/*.go`
- `testdata/workloads/*.xml`

**Requirements:**
- parse practical XML subset used by current `cosbench-sineio`
- normalize XML into internal `domain` structures
- return actionable validation errors

**Verification:**
- parser tests using real or representative workload XML
- invalid XML cases return explicit errors
</task>

### <task type="auto">
**Goal:** create a compatibility matrix for XML features.

**Touched files:**
- `docs/xml-compat-matrix.md`

**Requirements:**
- list supported / unsupported / deferred XML elements
- tie support decisions to real workload needs

**Verification:**
- matrix reflects current parser behavior
</task>

---

## Wave 3 — S3/SineIO Driver and Local Executor

### <task type="auto">
**Goal:** implement S3/SineIO endpoint configuration and client factory.

**Touched files:**
- `internal/driver/s3/*.go`
- `internal/domain/endpoint.go`

**Requirements:**
- support endpoint, region, access key, secret key
- support path-style/custom endpoint options as needed
- make SineIO-compatible configuration explicit

**Verification:**
- config validation tests
- client factory test coverage where feasible
</task>

### <task type="auto">
**Goal:** implement local executor for core object operations.

**Touched files:**
- `internal/executor/*.go`
- `internal/driver/s3/*.go`

**Requirements:**
- support core benchmark operations required by real workloads
- support local stage execution and worker goroutines
- emit metrics/events back to reporting layer

**Verification:**
- targeted executor tests
- local dry-run/integration path succeeds on representative config
</task>

---

## Wave 4 — Control Plane & Web UI

### <task type="auto">
**Goal:** implement control-plane lifecycle and job state machine.

**Touched files:**
- `internal/controlplane/*.go`
- `internal/app/*.go`

**Requirements:**
- create job
- start job
- track stage/job status
- persist snapshots during lifecycle

**Verification:**
- unit tests for lifecycle transitions
- restart recovery preserves visible state
</task>

### <task type="auto">
**Goal:** ship the first usable web UI.

**Touched files:**
- `internal/web/*.go`
- `web/templates/*.html`
- `web/static/*`

**Requirements:**
- dashboard
- workload upload page
- endpoint config page
- job detail page
- history page

**Verification:**
- manual browser walkthrough
- upload → run → inspect flow works end-to-end
</task>

---

## Wave 5 — Reporting, Hardening, and M3 Readiness

### <task type="auto">
**Goal:** improve reporting to migration-grade usefulness.

**Touched files:**
- `internal/reporting/*.go`
- `web/templates/job_detail*.html`

**Requirements:**
- throughput
- latency
- error counts
- percentile summaries
- stage-level summaries

**Verification:**
- metrics snapshots render correctly in UI
- summary output matches executor events
</task>

### <task type="auto">
**Goal:** document migration gaps and next remote-worker seams.

**Touched files:**
- `docs/migration-gap-analysis.md`
- `docs/remote-worker-seams.md`

**Requirements:**
- identify remaining gaps vs `cosbench-sineio`
- define future split between controller and remote worker
- keep M3 scope honest

**Verification:**
- written docs reflect current implementation reality
</task>

---

## Immediate Next Build Order
1. Wave 1
2. Wave 2
3. Wave 3
4. Minimal end-to-end check
5. Wave 4
6. Wave 5
