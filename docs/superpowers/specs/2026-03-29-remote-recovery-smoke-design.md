# Remote Recovery Smoke Design

## Goal

Add a local MinIO-backed remote smoke scenario that validates controller/driver recovery when a driver claims work and then disappears before completing it.

This slice is intended to move the current recovery behavior from unit and integration evidence into a real multi-process smoke path.

## Problem

The repository already has:

- remote controller-only and driver-only runtime modes
- work-unit scheduling and retry/requeue behavior
- lease expiry and stale-driver health logic
- local multi-process remote smoke for:
  - `s3 + single`
  - `s3 + multistage`
  - `sio + single`
  - `sio + multistage`

What is still missing is a smoke path that proves those recovery behaviors work across real HTTP boundaries and separate OS processes.

Without that, the most important durability-related remote behavior remains validated only below the system-smoke layer.

## Scope

### In Scope

- add one new local remote smoke scenario:
  - `scenario=recovery`
- use:
  - one controller-only process
  - two driver-only processes
  - local MinIO
- force one driver to claim and then disappear
- validate lease expiry and reassignment to the surviving driver
- emit recovery-specific fields in smoke summary output

### Out Of Scope

- adding recovery into the existing `Remote Smoke Matrix`
- adding a new GitHub workflow in this slice
- chaos or randomized failure injection
- repeated failure cycles or stress testing
- SIO recovery parity in the same slice

## Recommended Approach

Add one purpose-built `recovery` scenario to the existing `scripts/smoke_remote_local.py` helper.

### Why keep it local/manual first

- recovery smoke is heavier and slower than the current happy-path matrix
- the helper will need to deliberately kill a driver process and wait for lease expiry
- it is better to validate this flow locally first before deciding whether it belongs in a recurring matrix workflow

## Fixture Strategy

Add one new fixture:

- `testdata/workloads/remote-smoke-s3-recovery-two-driver.xml`

Recommended properties:

- `storage type="s3"`
- one stage
- one work
- `workers="2"`
- small `totalOps`
- write-only workload shape

The first version should stay on `s3` because the goal is validating recovery mechanics, not adding another backend parity dimension at the same time.

## Recovery Flow

The helper should orchestrate this sequence:

1. start MinIO, controller-only, driver-1, and driver-2
2. submit and start the recovery workload
3. wait until controller snapshot data shows at least one mission claimed by driver-1
4. stop driver-1 intentionally
5. wait for lease expiry and mission reassignment
6. allow driver-2 to finish the reassigned work
7. verify the final job succeeds

The smoke should prove at least one unit was reclaimed after driver-1 disappeared.

## Summary Additions

Extend the existing summary output with:

- `recovery_observed`
- `reclaimed_units`

These should appear in both:

- `summary.json`
- `summary.md`

The existing summary fields should remain stable.

## Checks

In addition to the existing process, health, distribution, success, and visibility checks, add:

- `recovery_observed`
  - at least one work unit is seen with more than one mission attempt or is claimed by a different driver after driver-1 is stopped
- `recovery_job_succeeded`
  - final job status is `succeeded`

## Repository Touch Points

### `testdata/workloads/remote-smoke-s3-recovery-two-driver.xml`

New fixture for the recovery scenario.

### `scripts/smoke_remote_local.py`

Modify:

- add `scenario=recovery`
- orchestrate driver termination and reassignment observation
- extend summary payload

### `scripts/test_smoke_remote_local.py`

Modify:

- fixture selection assertions
- summary shape assertions for recovery fields

### `README.md`

Add:

- one local command example for `scenario=recovery`

## Success Criteria

This slice is complete when:

1. `SMOKE_REMOTE_LOCAL_SCENARIO=recovery make --no-print-directory smoke-remote-local` succeeds locally
2. the helper proves driver loss followed by reassignment
3. summary output records that recovery was observed
4. existing `single` and `multistage` smoke scenarios do not regress

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
