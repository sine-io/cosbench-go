# Remote Worker Seams

The current system runs controller and execution in one process. This document records the seams intentionally preserved for a future split.

## Current Separation Points

- `internal/controlplane`: owns job lifecycle, stage ordering, snapshot persistence, event emission
- `internal/executor`: owns stage/work execution against a storage adapter
- `internal/driver/s3`: owns storage-client configuration and object operations
- `internal/reporting`: owns metrics aggregation and summaries
- `internal/snapshot`: owns recoverable persisted state

## Future Controller Responsibilities

- accept workload submissions
- validate configuration before dispatch
- persist authoritative job state and events
- coordinate stage and worker assignment
- aggregate worker-reported metrics into job-visible summaries

## Future Remote Worker Responsibilities

- receive executable work units for a stage/work slice
- run local goroutines against the target storage endpoint
- stream samples/events back to controller
- expose health and capability information

## Candidate Split Contract

- controller sends: workload fragment, endpoint config reference/materialized config, execution limits, worker allocation
- worker returns: lifecycle events, sample stream, final work summary, fatal failure reason

## Why The Current Layout Helps

- `internal/controlplane` does not embed storage-driver details directly
- `internal/executor` can evolve from in-process call to remote dispatch boundary
- `internal/reporting` can aggregate local or remote-originated samples with the same summary model

## What Is Still Missing Before Remote Split

- explicit work-unit protocol
- authenticated worker registration and liveness model
- backpressure/retry rules for sample/event delivery
- controller-side scheduling strategy across multiple workers
