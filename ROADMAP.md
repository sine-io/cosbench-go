# ROADMAP.md — cosbench-go

## Phase M1 — Minimum Usable Benchmark Loop
Goal: prove the Go system can run a real workload end-to-end.

### Deliverables
- single binary web server boots
- workload XML upload works
- XML parser supports the core subset used by current `cosbench-sineio` workloads
- S3/SineIO endpoint config can be supplied
- one benchmark job can be launched locally
- base operations supported: put / get / delete / list (as needed by current workloads)
- job state visible in UI
- basic metrics visible in UI
- snapshot persistence for jobs and results

### Exit Criteria
- one real benchmark can be executed from browser
- results remain inspectable after restart

---

## Phase M2 — Operationally Usable Replacement Candidate
Goal: make the platform feel practical, not demo-like.

### Deliverables
- more complete XML coverage for actively used workload structures
- richer stage/job lifecycle display
- error/event stream for failed jobs
- job history page
- endpoint config persistence and reuse
- better reporting: throughput, latency, percentiles, error counts
- more robust S3/SineIO configuration handling
- execution and reporting paths covered by testdata fixtures

### Exit Criteria
- daily-use scenarios no longer depend on old COSBench for the main path
- operators can diagnose common failures from UI/logs

---

## Phase M3 — Near-Mainline Migration Target
Goal: approach the functional completeness needed for practical migration from `cosbench-sineio`.

### Deliverables
- broad XML compatibility for high-value real workloads
- stable job lifecycle and recoverability
- web-first control plane good enough for regular team usage
- cleaner module boundaries for future remote worker split
- explicit driver abstraction for later non-S3 expansion
- stronger scenario validation before execution
- better result summaries/export

### Exit Criteria
- primary current benchmarking workflows can be performed in `cosbench-go`
- remaining gaps are explicit, narrow, and non-blocking for the main migration path

---

## Deferred Until After M3
- remote worker orchestration
- non-S3 driver expansion
- database-backed persistence
- SPA frontend
- auth / tenancy
- full parity with every historical COSBench feature

---

## Priority Order
1. correctness of execution
2. XML compatibility for real workloads
3. observability and reporting
4. recoverability via snapshots
5. UI polish
6. broader ecosystem extensibility
