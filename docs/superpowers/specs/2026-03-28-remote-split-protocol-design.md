# Remote Split Protocol Design

## Goal

Turn the preserved local-only seam into a real controller/driver protocol while keeping one Go codebase and preserving `combined` mode for local runs and CI.

## Problems To Solve

The repository already documents conceptual controller/worker seams, but critical pieces are still missing:

- explicit protocol DTOs
- driver registration and liveness
- mission assignment
- sample/event/result upload
- controller scheduling across multiple drivers
- retry and replay rules
- minimal authentication between controller and drivers

Without these, remote execution remains conceptual only.

## In Scope

- runtime modes `controller-only`, `driver-only`, `combined`
- shared remote DTOs and status model
- driver registration and heartbeat
- controller mission planning and assignment
- driver mission polling/claiming
- batched sample/event upload
- final mission completion and failure reporting
- controller scheduling across multiple registered drivers
- file-backed persistence for drivers, missions, and timeline ingestion
- minimal token-based trust model

## Out Of Scope

- non-S3 backend mission execution
- cluster-wide HA controller behavior
- external message buses
- advanced security hardening beyond a practical shared-secret baseline

## Recommended Protocol Shape

Use a driver-pull protocol rather than controller-push.

Reasons:

- simpler for NAT and container deployment
- easier to make testable inside one process
- cleaner retry semantics for drivers that reconnect

### Core flow

1. driver starts and registers with controller
2. driver sends heartbeats periodically
3. driver polls for mission leases
4. controller assigns a mission or work-unit batch
5. driver executes locally
6. driver uploads events and timeline/sample batches
7. driver reports final mission completion or fatal failure

### Required APIs

Suggested controller-facing driver APIs:

- `POST /api/driver/register`
- `POST /api/driver/heartbeat`
- `POST /api/driver/missions/claim`
- `POST /api/driver/missions/:id/events`
- `POST /api/driver/missions/:id/samples`
- `POST /api/driver/missions/:id/complete`

Controller-visible management APIs remain under `/api/controller/...`.

## Scheduling Model

The first remote-capable scheduler should stay pragmatic:

- controller decomposes a runnable job into stage-scoped `WorkUnit`s
- only one stage is active at a time
- controller assigns available work-units to healthy drivers
- controller waits for all work-units in the stage to complete before advancing

This preserves current stage ordering semantics while enabling multiple drivers within a stage.

## Shared State Model

Introduce explicit remote entities:

- `DriverNode`: identity, mode, health, capabilities, last heartbeat, token fingerprint
- `Mission`: controller-owned assignment record
- `WorkUnit`: executable work fragment and limits
- `MissionLease`: claim/expiry information
- `SampleEnvelope`: batch of time-series points or samples
- `EventEnvelope`: batch of lifecycle events

Controller is authoritative for all mission lifecycle state.

## Authentication Model

Use a simple shared-secret token model for the first phase:

- controller is configured with accepted driver token(s)
- driver presents token on register and heartbeat
- controller returns a driver identity and capability contract

This is intentionally minimal, but must be wired so stronger auth can replace it later.

## Combined Mode Requirement

`combined` mode must use the same remote abstractions internally where reasonable.

That means:

- local execution in combined mode should flow through the mission planner/scheduler boundary
- one in-process driver agent can register against the in-process controller

This reduces the risk of remote-only bugs and keeps local verification useful.

## Persistence And Replay

The controller should persist:

- registered drivers
- outstanding missions
- mission outcomes
- uploaded sample/event batches or controller-aggregated timeline state

Drivers should persist enough local state to resume safe mission reporting after transient connection loss.

## Failure Rules

Minimum required rules:

- missed heartbeats mark drivers unhealthy
- leased but uncompleted missions can be requeued after lease expiry
- duplicate sample/event batches must be idempotent or deduplicated
- fatal mission failure should not silently stall stage progression

## Testing Strategy

### Protocol tests

- registration
- heartbeat state transitions
- mission claim semantics
- lease expiry
- batch ingestion idempotence

### Integration tests

- controller-only plus driver-only processes
- combined mode remote-loopback execution
- multi-driver stage completion
- driver disconnect and mission reassignment

### Human verification

- start controller-only server
- start one or more driver-only servers
- submit a workload
- observe mission assignment, completion, and controller-visible aggregation

## Success Criteria

This slice is complete when:

1. controller-only, driver-only, and combined modes all run
2. drivers can register, heartbeat, claim missions, and report completion
3. controller can schedule a stage across multiple drivers
4. uploaded samples/events feed the same controller result/timeline model used by the web layer
5. failure and lease-expiry paths are covered by tests

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
