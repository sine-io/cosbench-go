# PROJECT.md — cosbench-go

## Project Definition
cosbench-go is a Go-native, Web-first object storage benchmarking platform designed to migrate the practical core of `cosbench-sineio` without inheriting its legacy OSGi/plugin runtime.

Primary target:
- Compatible with core COSBench workload XML input
- First-class support for S3 / SineIO-compatible object storage
- Single-binary deployment in early stages
- Web control plane first, with future remote worker expansion reserved

## Product Goal
Deliver a usable Go replacement for the main operational path of COSBench/SineIO benchmarking:
- upload/parse workload XML
- configure endpoint
- execute jobs
- inspect progress and results in a web UI
- preserve run history via file snapshots

## Non-Goals (Phase 1)
- Reimplement OSGi / Eclipse Equinox
- Full parity with every COSBench plugin
- Multi-driver support beyond S3/SineIO
- Full distributed worker cluster on day one
- SPA-first frontend architecture
- Multi-tenant auth/permission system

## Architecture Stance
Rebuild-first, not line-by-line port.

Chosen direction:
- Go-native single service first
- Internal module boundaries that allow later split into controller/worker roles
- Compatibility at the workload level, not runtime-implementation level

## Technical Stack
- Language: Go
- Frontend: server-rendered HTML/templates
- Storage: in-memory state + file snapshots
- Input: COSBench workload XML
- Driver: S3 / SineIO first
- Deployment: single binary

## Core Modules
- `cmd/server`: boot entrypoint
- `internal/domain`: core entities and contracts
- `internal/workloadxml`: XML parsing + normalization + validation
- `internal/controlplane`: job lifecycle + scheduling state machine
- `internal/executor`: stage execution + worker goroutines
- `internal/driver/s3`: S3/SineIO driver implementation
- `internal/reporting`: metrics aggregation and summaries
- `internal/snapshot`: file persistence for jobs/events/results
- `internal/web`: handlers + view models
- `web/templates`: UI templates

## Storage Strategy
Phase 1 uses:
- runtime state in memory
- snapshots persisted to files

Suggested directories:
- `data/jobs/`
- `data/results/`
- `data/events/`
- `data/endpoints/`

This is intentionally chosen over DB-first to reduce startup complexity while preserving recoverability.

## Execution Model
Phase 1:
- controller and executor run in one process
- jobs execute locally
- module seams remain clean enough to later introduce remote workers

Future direction:
- add remote worker protocol after local execution is stable

## Input Compatibility Contract
The system should prioritize compatibility with practical COSBench workload XML semantics that are actively used in `cosbench-sineio`.

Compatibility target order:
1. common workload XML structures used in current deployment
2. S3/SineIO operation mapping
3. stage execution semantics that affect benchmarking correctness
4. nice-to-have XML features later

## Web UX Contract
The first usable web UI must include:
- dashboard
- upload workload XML page
- endpoint configuration page
- job detail page
- history page

The UI does not need to be beautiful. It must be:
- stable
- inspectable
- operationally useful

## Quality Bar
The project is acceptable when:
- a real XML workload can be uploaded and parsed
- a real S3/SineIO benchmark can be started
- progress can be viewed in browser
- results can be inspected after completion
- runs survive process restart through snapshot persistence

## Key Risks
1. Underestimating XML compatibility complexity
2. SineIO-specific S3 behavior mismatches
3. Letting governance/UI features outrun execution correctness
4. Scope creep into “full COSBench parity” too early

## Scope Defense
When in doubt, prioritize:
1. benchmark correctness
2. XML compatibility for real workloads
3. operational observability
4. incremental extensibility

Do not prioritize cosmetic parity over execution integrity.
