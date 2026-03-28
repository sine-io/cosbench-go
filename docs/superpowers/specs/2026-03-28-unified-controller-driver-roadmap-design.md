# Unified Controller/Driver Migration Roadmap Design

## Goal

Define the shared architecture and sequencing for the next migration phase so `cosbench-go` can grow from a local-only v1 tool into a unified Go service that covers:

- compatibility gaps still marked partial or deferred within the S3/SIO migration boundary
- controller-facing web/API capability beyond the current dashboard/history subset
- a real controller/driver remote split
- driver-facing pages inside the same Go service

This document is the contract for the four sub-project specs that follow it.

## Context

Current `cosbench-go` explicitly closes only the local-only v1 boundary:

- one process hosts control-plane and execution
- XML compatibility is limited to the active S3/SIO subset
- web coverage is intentionally smaller than legacy controller-web / driver-web
- remote worker behavior is a seam, not a finished protocol

The user requested the remaining partial and deferred items to be migrated in Go, while also simplifying the legacy deployment model:

- use a unified API surface
- allow one Go service to expose both controller and driver capabilities
- do not spend this phase on non-S3 backends

## Architectural Decision

Use one Go server architecture with role-aware behavior instead of recreating the legacy Java multi-webapp packaging.

The service must support three runtime modes:

- `controller-only`
- `driver-only`
- `combined`

`combined` remains the easiest local and CI path. `controller-only` and `driver-only` make the eventual remote split real instead of simulated.

This is a packaging simplification, not a responsibility collapse:

- controller remains the authoritative owner of jobs, scheduling, events, summaries, and persisted state
- driver remains the owner of local worker execution, mission runtime, and health/capability reporting

## Route And API Model

The HTTP surface is unified, but role names stay explicit.

### HTML routes

- controller pages: `/controller/...`
- driver pages: `/driver/...`
- shared operational pages if needed: `/system/...`

### API routes

- controller resources: `/api/controller/...`
- driver resources: `/api/driver/...`
- shared operational resources: `/api/system/...`

This keeps the user-visible model simple while avoiding ambiguous ownership.

## Shared Domain Additions

The sub-projects must converge on a shared model before implementation details diverge.

### Existing entities that remain authoritative

- `Job`
- `Stage`
- `Work`
- `Endpoint`
- `JobResult`
- `JobEvent`

### New entities required by this roadmap

- `AuthSpec`
- `DriverNode`
- `Mission`
- `WorkUnit`
- `MissionLease`
- `SampleEnvelope`
- `EventEnvelope`
- `ArtifactRef`
- `TimeSeriesPoint`
- `StageTimeline`

### Ownership rules

- controller persists authoritative `Job`, `Mission`, `DriverNode`, `JobResult`, `ArtifactRef`, and timeline state
- driver persists only its local runtime cache and replay buffers
- any state shown on driver pages must be derivable from driver runtime and controller-visible mission state

## Persistence Direction

The roadmap stays on file-backed snapshots for now.

The snapshot layout should expand incrementally rather than being replaced wholesale. New sub-projects may add directories such as:

- `data/drivers/`
- `data/missions/`
- `data/artifacts/`
- `data/timelines/`

Database adoption remains out of scope for this roadmap.

## Sub-Project Boundaries

### 1. Compatibility closure

Close the current parser-only or partially aligned runtime gaps inside the S3/SIO family:

- `filewrite`
- explicit `<auth>` modeling and runtime resolution
- `siov1` / `gdas` runtime behavior beyond alias-only parsing
- `prefetch` / `range-read` execution semantics
- the concrete S3/SIO delta list already captured in repository docs

### 2. Controller web/API closure

Add controller-facing parity for the deferred capability set, but through the new unified route model instead of legacy page duplication:

- matrix
- config / advanced config
- stage detail
- timeline + timeline CSV
- Prometheus export
- config/log downloads

### 3. Remote split protocol

Turn the current seam into an actual controller/driver protocol:

- registration
- heartbeat
- mission assignment
- sample/event upload
- final result upload
- multi-driver scheduling

### 4. Driver pages inside the unified service

Add driver-facing pages using the same Go service and shared APIs. No separate driver webapp is required.

### Excluded from this roadmap

- non-S3 backends
- legacy packaging parity
- DB-backed persistence
- byte-identical Java-era report surfaces

## Delivery Order

The implementation order is fixed:

1. Compatibility closure
2. Controller web/API closure
3. Remote split protocol
4. Driver pages in unified service

Reasons:

- compatibility work stabilizes runtime semantics and shared config behavior first
- controller pages need better query surfaces and timeline artifacts before remote execution lands
- remote split must be designed against a stable controller data model
- driver pages depend on real driver node and mission state from the remote split

## Testing Policy

Every sub-project must define three verification layers:

- unit tests for parser/domain/protocol behavior
- integration tests for HTTP routes, persistence, and runtime interactions
- repeatable CLI or HTTP verification steps for humans and CI

No sub-project should rely on UI-only verification to prove correctness.

## Success Criteria

This roadmap is ready for implementation when:

1. all four sub-project specs use the same route model and entity vocabulary
2. all four plans assume the same delivery order
3. no plan depends on non-S3 backend work
4. controller and driver responsibilities remain distinct even when hosted in one process

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
