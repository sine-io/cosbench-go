# Remote Split Protocol Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn the current in-process seam into a real controller/driver protocol while preserving a useful `combined` mode for local execution and CI.

**Architecture:** Define shared remote DTOs and persisted remote state first, then implement controller-side registration and scheduling, then add a driver agent and unified HTTP endpoints, and finally route `combined` mode through the same abstractions for end-to-end verification.

**Tech Stack:** Go 1.26, existing `internal/controlplane`, `internal/snapshot`, `internal/web`, file-backed persistence, HTTP handlers in the existing server process

---

### Task 1: Add Shared Remote State And Protocol Coverage

**Files:**
- Create: `internal/domain/remote.go`
- Create: `internal/controlplane/remote_state.go`
- Create: `internal/controlplane/remote_state_test.go`
- Modify: `internal/snapshot/store.go`
- Modify: `internal/domain/job.go`

- [ ] **Step 1: Write failing tests for driver, mission, and lease state**

Cover:
- driver registration model
- mission lifecycle states
- mission lease expiry metadata
- snapshot persistence for remote state

- [ ] **Step 2: Run focused controller tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failures because remote entities and snapshot handling do not exist yet

- [ ] **Step 3: Implement shared remote DTOs and controller-owned remote state**

Add the domain types and persistence shape needed for later HTTP and scheduling work.

- [ ] **Step 4: Re-run focused controller tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- controller package passes with remote state coverage

- [ ] **Step 5: Commit the remote-model slice**

Run:
```bash
git add internal/domain/remote.go internal/controlplane/remote_state.go internal/controlplane/remote_state_test.go internal/snapshot/store.go internal/domain/job.go
git commit -m "feat: add remote mission state model"
```

### Task 2: Implement Controller Registration, Heartbeat, And Scheduling

**Files:**
- Create: `internal/controlplane/driver_registry.go`
- Create: `internal/controlplane/mission_scheduler.go`
- Create: `internal/controlplane/mission_scheduler_test.go`
- Create: `internal/web/driver_api.go`
- Create: `internal/web/driver_api_test.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/web/handler.go`

- [ ] **Step 1: Write failing tests for registration, heartbeat, claim, and scheduling**

Cover:
- driver registration
- heartbeat-driven health transitions
- mission claim semantics
- multi-driver stage scheduling

- [ ] **Step 2: Run focused controller and web tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- failures because registry, scheduler, and driver-facing HTTP routes are not implemented

- [ ] **Step 3: Implement controller registry and mission scheduler**

Make controller authoritative for:
- driver liveness
- work-unit creation
- mission assignment
- lease expiry and requeue rules

- [ ] **Step 4: Implement registration, heartbeat, and claim HTTP handlers**

Expose:
- `POST /api/driver/register`
- `POST /api/driver/heartbeat`
- `POST /api/driver/missions/claim`

- [ ] **Step 5: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- controller and web packages pass for registration/scheduling coverage

- [ ] **Step 6: Commit the controller protocol slice**

Run:
```bash
git add internal/controlplane/driver_registry.go internal/controlplane/mission_scheduler.go internal/controlplane/mission_scheduler_test.go internal/web/driver_api.go internal/web/driver_api_test.go internal/controlplane/manager.go internal/web/handler.go
git commit -m "feat: add driver registration and mission scheduling"
```

### Task 3: Add Driver Agent Runtime And Mission Reporting

**Files:**
- Create: `internal/driver/agent/agent.go`
- Create: `internal/driver/agent/http_client.go`
- Create: `internal/driver/agent/agent_test.go`
- Modify: `internal/executor/executor.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/web/driver_api.go`

- [ ] **Step 1: Write failing tests for driver mission execution and batch reporting**

Cover:
- mission polling
- local execution of claimed work
- batched event upload
- batched sample upload
- final completion upload

- [ ] **Step 2: Run focused agent/controller tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/agent ./internal/controlplane ./internal/web
```

Expected:
- failures because no driver agent exists and reporting endpoints are incomplete

- [ ] **Step 3: Implement the driver agent and reporting endpoints**

Add:
- driver-side polling loop
- event/sample batch upload
- mission completion/failure reporting
- controller-side ingestion into the existing result/timeline model

- [ ] **Step 4: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/agent ./internal/controlplane ./internal/web
```

Expected:
- agent, controller, and web packages pass

- [ ] **Step 5: Commit the driver runtime slice**

Run:
```bash
git add internal/driver/agent/agent.go internal/driver/agent/http_client.go internal/driver/agent/agent_test.go internal/executor/executor.go internal/controlplane/manager.go internal/web/driver_api.go
git commit -m "feat: add driver agent mission reporting"
```

### Task 4: Add Runtime Modes, Combined Loopback, And Final Verification

**Files:**
- Create: `internal/app/mode.go`
- Create: `internal/app/remote_integration_test.go`
- Modify: `internal/app/app.go`
- Modify: `cmd/server/main.go`
- Modify: `docs/remote-worker-seams.md`
- Modify: `README.md`

- [ ] **Step 1: Write failing integration tests for `controller-only`, `driver-only`, and `combined`**

Cover:
- controller-only server startup
- driver-only server startup
- combined mode in-process loopback mission execution

- [ ] **Step 2: Run integration tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/app
```

Expected:
- failures because mode wiring and loopback execution do not exist yet

- [ ] **Step 3: Implement runtime mode wiring and combined-mode loopback**

Make `combined` use the same scheduler/agent abstractions instead of a separate shortcut path.

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

- [ ] **Step 6: Commit the remote-split closure slice**

Run:
```bash
git add internal/app/mode.go internal/app/remote_integration_test.go internal/app/app.go cmd/server/main.go docs/remote-worker-seams.md README.md
git commit -m "feat: add controller driver runtime modes"
```
