# Work-Unit Scheduling Design

## Goal

Turn the current remote split from “one work maps to one remote mission” into “one work can be decomposed into multiple work-units distributed across multiple drivers”.

The design should preserve the existing stage-serial model while making stage-internal parallelism truly remote-capable.

## Problem

The current remote split has real protocol flows, but its scheduling model is still coarse:

- a work is effectively treated as one remotely executable unit
- multiple remote drivers do not yet share execution of a single work
- retry and reassignment are mission-centric rather than slice-centric

This limits the value of the remote split because a multi-worker work cannot yet fan out across multiple drivers in a principled way.

## Scope

### In Scope

- `WorkUnit` as a stable scheduling object
- `MissionAttempt` as an execution attempt over one `WorkUnit`
- stage-internal splitting of a work into `N` units based on `workers=N`
- multi-driver claim of units from the same stage
- reassignment of expired or failed unit attempts
- aggregation of unit results back into work/stage/job summaries
- tests for multi-driver distribution and retry

### Out Of Scope

- non-S3 backend execution
- cross-stage parallelism
- dynamic load-based rebalancing
- advanced object-range partition strategies
- cluster-wide HA scheduling

## Recommended Approach

### 1. Introduce two layers of remote execution identity

Use:

- `WorkUnit` for the stable slice of work
- `MissionAttempt` for a particular claim/lease/driver execution attempt

This separation is important because:

- a unit can be retried multiple times
- the retry history should not mutate the identity of the underlying slice
- lease expiry and driver reassignment belong to attempts, not to the unit itself

### 2. Use worker-slice decomposition first

The first split strategy should be based on the current `workers` field:

- if a work has `workers=N`
- controller creates `N` `WorkUnit`s
- each `WorkUnit` carries:
  - `unit_index`
  - `unit_count`
  - a slice descriptor with `worker_index` and `worker_count`

This allows the current execution targeting logic to be reused because it already computes targets from the `(idx, all)` pair.

### 3. Keep stage ordering unchanged

The scheduler should continue to treat stages as serial barriers:

- all units from the active stage must finish
- only then may the next stage be scheduled

That preserves current semantics while allowing remote parallelism within one stage.

### 4. Make retries attempt-based with a fixed ceiling

For the first version, use a simple retry model:

- each `WorkUnit` may create a new `MissionAttempt` when the previous one expires or fails
- maximum attempts per unit: `3`

If a unit exceeds that limit:

- mark the unit failed
- fail its parent work/stage/job

This avoids infinite retry loops while keeping behavior easy to reason about.

## Data Model

### `WorkUnit`

Suggested fields:

- `id`
- `job_id`
- `stage_name`
- `work_name`
- `unit_index`
- `unit_count`
- `work`
- `storage`
- `slice`
- `status`
- `created_at`
- `updated_at`

### `slice`

Suggested first-version fields:

- `worker_index`
- `worker_count`

Do not add container/object range math yet unless tests prove it is required.

### `MissionAttempt`

Suggested fields:

- `id`
- `work_unit_id`
- `driver_id`
- `attempt`
- `status`
- `lease`
- `error_message`
- `created_at`
- `updated_at`

### Status layering

- unit status describes the stable slice lifecycle
- attempt status describes the currently running or past execution attempt

## Scheduling Flow

1. controller selects the active stage
2. controller expands each work into `workers` count of `WorkUnit`s
3. controller creates one initial `MissionAttempt` per unit
4. healthy drivers claim pending attempts
5. driver executes the claimed unit using `worker_index` / `worker_count`
6. controller ingests events/samples/completion into the unit and parent aggregates
7. expired or failed attempts create a new attempt unless attempt limit is reached
8. the stage completes when all units are terminal and successful, or fails when any unit becomes terminal-failed

## Failure And Recovery Rules

### Lease expiry

- expired attempt becomes terminal-expired
- unit returns to requeueable state if attempt limit not reached

### Driver failure report

- failed attempt becomes terminal-failed
- unit retries if below attempt ceiling

### Heartbeat timeout

- stale driver cannot claim new attempts
- in-flight reassignment still depends on lease expiry

This keeps heartbeat and lease responsibilities separate.

## Aggregation Strategy

Keep unit-level detail internal for now.

Public summaries should still roll up to:

- work
- stage
- job

The controller should merge:

- samples from all attempts of the successful unit lineage
- then all units of a work
- then all works of a stage

This lets existing reporting/UI surfaces remain useful while the scheduler becomes more capable underneath.

## Testing Strategy

### Controller tests

- one work with `workers=3` creates `3` units
- different healthy drivers claim different attempts from the same stage
- expired attempt causes a new claimable attempt
- attempt ceiling causes unit failure and stage/job failure

### Integration tests

- combined-mode loopback with unit decomposition
- multiple drivers completing one stage
- driver timeout followed by lease expiry and reassignment

## Success Criteria

This slice is complete when:

1. one work is decomposed into `workers`-count `WorkUnit`s
2. multiple drivers can execute units from the same stage concurrently
3. retries are tracked as attempts over stable units
4. failed or expired units are retried up to the configured ceiling
5. final results still aggregate correctly to work/stage/job

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
