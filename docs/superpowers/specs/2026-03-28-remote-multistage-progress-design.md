# Remote Multistage Progress Design

## Goal

Extend the current remote controller/driver execution path so a multi-stage workload can progress automatically from one stage to the next when the current stage completes successfully.

The design should preserve the existing rule that stages remain serial barriers while allowing the already-landed remote `work-unit` scheduling to operate within each stage.

## Problem

The repository now has:

- remote controller/driver protocol
- shared-token auth on driver write endpoints
- work-unit decomposition within a stage
- multi-driver work-unit execution
- local and GitHub-backed remote smoke paths

But the core remote scheduling path still selects only the first stage:

- it creates work units from `Workflow.Stages[0]`
- it does not automatically advance to `Stages[1]`, `Stages[2]`, etc.

That means the remote path is still incomplete for realistic workloads containing more than one stage.

## Scope

### In Scope

- track one active remote stage per job
- schedule the active stage only
- automatically schedule the next stage when the current stage succeeds
- stop the job if any current-stage work unit reaches terminal failure
- tests for:
  - automatic progression
  - no early next-stage claims
  - failure stops progression

### Out Of Scope

- mixing local and remote stage execution in one job
- cross-stage parallelism
- conditional stage branching
- advanced scheduler policies beyond simple serial stage progression

## Recommended Approach

### 1. Add explicit active-stage tracking

The controller should track which stage is currently eligible for remote scheduling.

Suggested job-owned state:

- `active_stage_index`

This should be persisted with job-visible state rather than treated as ephemeral scheduler memory.

### 2. Make stage scheduling index-based

Replace the current hard-coded use of the first stage with an index-aware scheduler:

- schedule `Workflow.Stages[active_stage_index]`
- do not create units or attempts for later stages yet

This keeps the remote path consistent with the existing serial stage model.

### 3. Advance only after all current-stage units succeed

The advancement rule should be based on `WorkUnit` terminal status, not on attempt history.

Current stage succeeds when:

- every `WorkUnit` for the active stage has `status == succeeded`

At that point:

- if another stage exists:
  - increment `active_stage_index`
  - create units and attempts for the next stage
- else:
  - mark the job `succeeded`

### 4. Fail fast on terminal work-unit failure

If any unit in the active stage reaches terminal `failed` after exhausting its retry ceiling:

- mark the active stage `failed`
- mark the job `failed`
- do not schedule later stages

### 5. Keep claim scope limited to the active stage

Claiming should only see pending attempts from the active stage.

The safer implementation is:

- do not create later-stage attempts before they are eligible

That removes the need for claim-time filtering across future stages.

## State Model

### Job additions

Suggested additions:

- `active_stage_index`
- optional `remote_mode` or similar marker only if needed for clarity

### Existing entities reused

- `WorkUnit`
- `MissionAttempt`
- `JobResult`
- `JobTimeline`

No new remote identity type is needed in this slice.

## Scheduling Flow

1. `StartJob()` enters remote scheduling mode
2. controller initializes `active_stage_index = 0`
3. controller schedules stage `0`
4. drivers claim and execute work units for stage `0`
5. controller aggregates unit outcomes
6. if stage `0` succeeded:
   - advance to stage `1`
   - schedule stage `1`
7. repeat until:
   - last stage succeeds, or
   - any active-stage unit becomes terminal-failed

## Testing Strategy

### Control-plane tests

- two-stage job:
  - stage `0` success creates stage `1` units automatically
- next stage cannot be claimed early
- stage `0` terminal failure prevents stage `1` scheduling

### Integration tests

- `combined` mode processes a two-stage job end-to-end
- `controller-only + driver-only` processes a two-stage job end-to-end

### Smoke follow-up

This design should enable a later remote smoke fixture with more than one stage, but that fixture is not required to land in the same slice.

## Success Criteria

This slice is complete when:

1. remote scheduling no longer assumes `Stages[0]` forever
2. multi-stage jobs progress automatically stage by stage
3. later stages are not scheduled early
4. a terminal failed unit stops further stage progression
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
