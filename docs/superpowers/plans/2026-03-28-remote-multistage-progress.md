# Remote Multistage Progress Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make remote execution progress automatically from one stage to the next when the current stage completes successfully.

**Architecture:** Introduce explicit active-stage tracking on the job, schedule only the current stage, and trigger next-stage scheduling from the controller’s existing unit/attempt aggregation path. Keep stage ordering strictly serial and reuse existing work-unit and attempt state.

**Tech Stack:** Go 1.26, existing `internal/controlplane`, `internal/domain`, `internal/app`, file-backed snapshots

---

### Task 1: Add Failing Multistage Control-Plane Tests

**Files:**
- Modify: `internal/controlplane/mission_scheduler_test.go`

- [ ] **Step 1: Write failing tests for automatic multistage progression**

Cover:
- a two-stage remote job schedules stage 1 automatically after stage 0 success
- a later stage is not claimable before the current stage succeeds
- a terminal failed unit prevents stage 2 scheduling

- [ ] **Step 2: Run focused control-plane tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failures because remote scheduling still assumes the first stage only

### Task 2: Add Explicit Active-Stage State

**Files:**
- Modify: `internal/domain/job.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/snapshot/store.go`

- [ ] **Step 1: Add persisted active-stage tracking to the job model**

Store which stage is currently eligible for remote scheduling.

- [ ] **Step 2: Initialize active-stage state when remote scheduling begins**

Make `StartJob()` set up the first remote stage deterministically.

- [ ] **Step 3: Re-run focused control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- still failing on progression logic, but no longer failing for missing active-stage state

### Task 3: Make Stage Scheduling Index-Aware

**Files:**
- Modify: `internal/controlplane/mission_scheduler.go`

- [ ] **Step 1: Replace hard-coded `Stages[0]` selection**

Schedule based on the job’s active-stage index.

- [ ] **Step 2: Ensure only the active stage creates units and attempts**

Do not pre-create future-stage attempts.

- [ ] **Step 3: Re-run focused control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- next-stage early-claim tests pass
- progression/failure tests may still be red

- [ ] **Step 4: Commit the index-aware scheduling slice**

Run:
```bash
git add internal/domain/job.go internal/controlplane/manager.go internal/snapshot/store.go internal/controlplane/mission_scheduler.go internal/controlplane/mission_scheduler_test.go
git commit -m "feat: track active remote stage"
```

### Task 4: Advance Or Stop From Aggregation

**Files:**
- Modify: `internal/controlplane/mission_scheduler.go`
- Modify: `internal/controlplane/manager.go`

- [ ] **Step 1: Trigger next-stage scheduling when all current-stage units succeed**

If the current stage finishes successfully:
- increment the active-stage index
- schedule the next stage immediately

- [ ] **Step 2: Stop progression when the current stage terminally fails**

If any current-stage unit reaches terminal failure:
- mark the job failed
- do not schedule later stages

- [ ] **Step 3: Re-run focused control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- multistage progression tests pass

### Task 5: Add Integration Coverage For Multistage Remote Flow

**Files:**
- Modify: `internal/app/remote_integration_test.go`

- [ ] **Step 1: Write failing integration tests for multistage remote execution**

Cover:
- combined-mode multistage progression
- controller-only plus driver-only multistage progression

- [ ] **Step 2: Run focused app tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/app
```

Expected:
- failures until the multistage remote flow is wired correctly

- [ ] **Step 3: Adjust integration wiring only if needed**

Keep changes focused on enabling the multistage progression path, not on adding new features.

- [ ] **Step 4: Re-run focused app tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/app
```

Expected:
- app integration tests pass for multistage remote progression

- [ ] **Step 5: Commit the integration slice**

Run:
```bash
git add internal/app/remote_integration_test.go internal/controlplane/mission_scheduler.go internal/controlplane/manager.go
git commit -m "feat: progress remote jobs across stages"
```

### Task 6: Final Verification

**Files:**
- Review only: `internal/controlplane/mission_scheduler.go`
- Review only: `internal/controlplane/manager.go`
- Review only: `internal/app/remote_integration_test.go`

- [ ] **Step 1: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass

- [ ] **Step 2: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 3: Review final scope**

Run:
```bash
git diff -- internal/domain/job.go internal/controlplane internal/app
```

Expected:
- the slice stays focused on multistage remote progression

- [ ] **Step 4: Commit any remaining focused cleanup**

Run:
```bash
git add internal/domain/job.go internal/controlplane internal/app
git commit -m "refactor: finalize remote multistage progression"
```

Skip this commit if no additional cleanup was needed.
