# Work-Unit Scheduling Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Introduce `work-unit` decomposition so one work can be distributed across multiple remote drivers within the same stage.

**Architecture:** Add stable `WorkUnit` records plus retryable `MissionAttempt` records, decompose each work by `workers=N`, keep stages serial, and aggregate unit-level results back into the existing work/stage/job summaries.

**Tech Stack:** Go 1.26, existing `internal/controlplane`, `internal/domain`, `internal/driver/agent`, `internal/web`, file-backed snapshots

---

### Task 1: Add Failing Work-Unit And Attempt Tests

**Files:**
- Modify: `internal/controlplane/mission_scheduler_test.go`
- Create: `internal/controlplane/work_unit_test.go`

- [ ] **Step 1: Write failing tests for work-unit decomposition**

Cover:
- one work with `workers=3` creates `3` stable units
- each unit has `unit_index` and `unit_count`
- an initial mission attempt exists for each unit

- [ ] **Step 2: Run focused control-plane tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failures because work units and mission attempts do not exist yet

### Task 2: Implement WorkUnit And MissionAttempt State

**Files:**
- Modify: `internal/domain/remote.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/controlplane/remote_state.go`
- Modify: `internal/snapshot/store.go`
- Create: `internal/controlplane/work_unit.go`

- [ ] **Step 1: Add stable work-unit and attempt domain types**

Include:
- `WorkUnit`
- `MissionAttempt`
- unit/attempt status fields

- [ ] **Step 2: Persist units and attempts through snapshot state**

The controller should be able to reload them on restart.

- [ ] **Step 3: Re-run focused control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- decomposition/state tests pass

- [ ] **Step 4: Commit the state-model slice**

Run:
```bash
git add internal/domain/remote.go internal/controlplane/manager.go internal/controlplane/remote_state.go internal/snapshot/store.go internal/controlplane/work_unit.go internal/controlplane/work_unit_test.go internal/controlplane/mission_scheduler_test.go
git commit -m "feat: add remote work unit state"
```

### Task 3: Decompose Works Into Units And Claim Attempts

**Files:**
- Modify: `internal/controlplane/mission_scheduler.go`
- Modify: `internal/controlplane/driver_registry.go`
- Modify: `internal/web/driver_api.go`
- Modify: `internal/web/driver_api_test.go`

- [ ] **Step 1: Write failing tests for multi-driver same-stage claims**

Cover:
- two drivers can claim different attempts from the same stage
- unit claim order is deterministic enough for tests

- [ ] **Step 2: Run focused control-plane and web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- failures because scheduling is still work-level, not unit-level

- [ ] **Step 3: Implement work decomposition in the scheduler**

For each work:
- create `workers` count of units
- create an initial attempt per unit
- claim returns an attempt carrying one unit slice

- [ ] **Step 4: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- multi-driver claim tests pass

- [ ] **Step 5: Commit the scheduler slice**

Run:
```bash
git add internal/controlplane/mission_scheduler.go internal/controlplane/driver_registry.go internal/web/driver_api.go internal/web/driver_api_test.go internal/controlplane/mission_scheduler_test.go
git commit -m "feat: schedule work units across drivers"
```

### Task 4: Execute Units Through The Driver Agent

**Files:**
- Modify: `internal/driver/agent/agent.go`
- Modify: `internal/driver/agent/http_client.go`
- Modify: `internal/driver/agent/agent_test.go`
- Modify: `internal/executor/executor.go`

- [ ] **Step 1: Write failing tests for unit-slice execution**

Cover:
- claimed unit passes `worker_index` / `worker_count` into execution
- multiple units from the same work aggregate back correctly

- [ ] **Step 2: Run focused agent tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/agent
```

Expected:
- failures because the driver agent still executes whole-work missions

- [ ] **Step 3: Implement unit-slice execution**

The driver should execute the claimed unit using the current execution engine with the unit’s worker slice metadata.

- [ ] **Step 4: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/agent
```

Expected:
- agent tests pass with work-unit execution

- [ ] **Step 5: Commit the driver execution slice**

Run:
```bash
git add internal/driver/agent/agent.go internal/driver/agent/http_client.go internal/driver/agent/agent_test.go internal/executor/executor.go
git commit -m "feat: execute remote work units"
```

### Task 5: Retry Ceiling And Aggregation

**Files:**
- Modify: `internal/controlplane/mission_scheduler.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/controlplane/mission_scheduler_test.go`

- [ ] **Step 1: Write failing tests for attempt ceiling and aggregation**

Cover:
- expired/failed attempts create a new attempt up to the ceiling
- exceeding the ceiling fails the parent unit and stage/job
- successful unit aggregation still rolls up correctly to work/stage/job

- [ ] **Step 2: Run focused control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failures because retry ceiling and unit-aware aggregation are not implemented yet

- [ ] **Step 3: Implement retry ceiling and final unit aggregation**

Use:
- stable unit identity
- retryable attempt identity
- existing reporting merge path for roll-up

- [ ] **Step 4: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass

- [ ] **Step 5: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 6: Commit the retry and aggregation slice**

Run:
```bash
git add internal/controlplane/mission_scheduler.go internal/controlplane/manager.go internal/controlplane/mission_scheduler_test.go
git commit -m "feat: retry and aggregate remote work units"
```
