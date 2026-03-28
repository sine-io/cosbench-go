# Remote Reliability Hardening Design

## Goal

Harden the current remote controller/driver protocol without expanding scope beyond the existing unified Go service.

This slice focuses on three reliability gaps that remain after the initial remote split landed:

- lease expiry and mission requeue
- idempotent mission reporting for events, samples, and completion
- heartbeat timeout transitioning drivers to `unhealthy`

## Context

The current remote split is functionally present:

- drivers can register and heartbeat
- controller can schedule and claim missions
- driver agents can execute work and report results
- `combined` mode can run the same HTTP loopback path in-process

What is still weak is fault handling:

- an expired mission lease is not reclaimed automatically
- repeated event/sample uploads can double-apply state
- drivers remain healthy indefinitely unless they explicitly heartbeat again

These are reliability issues, not feature gaps.

## Scope

### In Scope

- on-demand lease expiry cleanup and requeue during:
  - mission claim
  - driver heartbeat
  - combined loopback processing
- idempotent handling for:
  - mission event batches
  - mission sample batches
  - mission completion calls
- heartbeat timeout policy and read/scheduling-time health refresh
- tests for the above behavior

### Out Of Scope

- background sweep loops
- persistent retry queues
- cryptographic request signing
- non-S3 backend mission execution
- advanced distributed scheduling strategy

## Recommended Approach

### 1. Lease expiry should be reclaimed on demand

Do not add a background sweeper yet.

Instead, introduce one controller helper that scans mission state against `now` and:

- finds missions in `claimed` or `running`
- checks `mission.Lease.ExpiresAt`
- if expired:
  - mark mission `expired`
  - clear lease
  - persist the state change

Run this helper from:

- `ClaimMission`
- `RecordDriverHeartbeat`
- `ProcessCombinedMission`

This keeps the implementation deterministic and easy to test.

### 2. Mission reporting should be mission-local idempotent

Add mission-local batch identity:

- `batch_id`
- `mission_id`
- payload kind (`events` or `samples`)

Controller stores received batch ids per mission, separately for events and samples.

If the same batch is replayed:

- return success
- do not re-append events
- do not re-accumulate samples

Completion should also be idempotent:

- if mission already reached a terminal state, return success with no duplicate aggregation

### 3. Driver health should refresh at read and claim time

Do not create a background heartbeat watcher yet.

Instead, add one helper that compares `LastHeartbeatAt` against a timeout window and marks the driver `unhealthy` when stale.

Call this helper from:

- `ListDriverNodes`
- `GetDriverOverview`
- `ClaimMission`

This keeps health state fresh where it matters operationally.

### 4. Keep heartbeat and lease semantics separate

Heartbeat timeout means:

- the driver should not receive new work

Lease expiry means:

- the existing work can be reassigned

Do not automatically revoke active leases just because a driver timed out.
That coupling would make the model harder to reason about and harder to tune later.

## Data Model Additions

Expected additions:

- per-mission received event batch ids
- per-mission received sample batch ids
- optional completion marker or terminal completion timestamp
- driver health timeout constant or config default

These can remain file-backed with the existing snapshot store.

## Testing Strategy

### Lease expiry

- expired claimed mission becomes `expired` and is available for re-claim
- unexpired mission is not stolen

### Idempotence

- duplicate event batch does not double-append logs
- duplicate sample batch does not double-count metrics
- duplicate completion call does not re-run aggregation

### Heartbeat timeout

- stale driver becomes `unhealthy`
- `unhealthy` driver cannot claim new work
- healthy driver still can

## Success Criteria

This slice is complete when:

1. expired leases are reclaimed on demand and can be re-claimed
2. repeated event/sample uploads are idempotent per mission batch
3. repeated completion calls do not duplicate mission/job aggregation
4. stale drivers transition to `unhealthy`
5. unhealthy drivers cannot claim new missions

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
