# Remote Multi-Process MinIO Smoke Design

## Goal

Add a repeatable local smoke path that validates the current remote controller/driver protocol in a real multi-process setup against a local MinIO endpoint.

This slice is intended to prove that the recently landed remote-split work is not just passing in-process loopback tests, but is also viable across real HTTP boundaries with separate controller and driver processes.

## Problem

The repository now contains:

- controller-only / driver-only / combined runtime modes
- remote registration, heartbeat, claim, report, and completion APIs
- shared-token auth for driver write paths
- work-unit scheduling across drivers

What is still missing is a concise, repeatable smoke run that verifies those capabilities together outside the single-process loopback path.

Without this, the remote split remains under-verified at the system boundary where it matters most:

- multiple OS processes
- real localhost networking
- real MinIO-backed object operations
- real work-unit distribution across multiple drivers

## Scope

### In Scope

- one local smoke helper for:
  - one `controller-only`
  - two `driver-only`
  - one local MinIO
- one minimal S3-focused remote smoke fixture
- deterministic health, scheduling, execution, aggregation, and visibility checks
- artifact capture under `.artifacts/remote-smoke/`
- machine-readable and human-readable summaries

### Out Of Scope

- GitHub-hosted remote smoke automation in this slice
- non-S3 backend smoke coverage
- browser automation
- broad UI testing beyond basic state visibility
- performance benchmarking

## Recommended Approach

Implement one purpose-built local smoke helper script, rather than a collection of manual commands.

### Why a single helper

- process orchestration is the hard part, not any one command
- the helper can generate ports, temp dirs, and tokens safely
- the helper can guarantee cleanup and write stable artifacts
- the helper can evolve into later CI automation if needed

## Runtime Topology

The helper should launch:

1. MinIO
2. controller-only server
3. driver-only server #1
4. driver-only server #2

All four processes run locally with isolated data directories.

The helper should generate and inject:

- random open ports
- one shared driver token
- one remote smoke output directory

## Workload Selection

Start with one minimal S3 fixture:

- `testdata/workloads/remote-smoke-s3-two-driver.xml`

Required properties:

- `storage type="s3"`
- one stage
- one work
- `workers="2"`
- small `totalOps`
- `write` only

The first version should avoid read-side dependencies so the smoke path isolates remote scheduling and reporting rather than workload preparation complexity.

SIO-specific remote smoke can be added later after the S3 path is stable.

## Checks

The helper should perform these hard checks, and fail the run if any of them do not hold.

### 1. Process readiness

- controller HTTP responds
- both driver HTTP surfaces respond
- MinIO endpoint responds

### 2. Driver registration and health

- controller sees exactly two drivers
- both are healthy

### 3. Work-unit distribution

- the submitted stage produces at least two units/attempts
- both drivers participate in claimed or completed attempts

### 4. Execution completion

- job reaches `succeeded`
- result metrics are non-zero

### 5. Aggregation integrity

- work/stage/job summaries exist
- operation and byte counts are not duplicated beyond the expected unit total

### 6. Control-plane visibility

- controller APIs expose job detail, timeline, and matrix data
- driver APIs expose overview and missions data

## Output Artifacts

Write artifacts under:

- `.artifacts/remote-smoke/`

Minimum recommended files:

- `summary.json`
- `summary.md`
- `controller.log`
- `driver1.log`
- `driver2.log`
- `minio.log`

### `summary.json`

Suggested keys:

- `controller_url`
- `driver_urls`
- `job_id`
- `job_status`
- `drivers_seen`
- `units_claimed`
- `drivers_participated`
- `operation_count`
- `byte_count`
- `checks`
- `overall`

### `summary.md`

A concise human-readable summary suitable for local inspection or later reuse in CI summaries.

## Failure Behavior

The helper should be strict:

- any failed check returns non-zero
- any spawned process failure returns non-zero
- all failures must still emit `summary.json` and `summary.md` with the reason captured

Cleanup must still happen on failure.

## Implementation Shape

The smallest coherent implementation is:

1. add the smoke fixture
2. add the orchestration helper
3. add result extraction against existing controller and driver APIs
4. add tests for helper internals where cheap and stable
5. document the command in README or adjacent docs if the helper proves stable

## Success Criteria

This slice is complete when:

1. a local helper can start MinIO, controller-only, and two driver-only processes
2. a minimal remote smoke fixture succeeds end-to-end
3. both drivers actually participate in unit execution
4. summary artifacts are written under `.artifacts/remote-smoke/`
5. the helper exits non-zero when any required check fails

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
