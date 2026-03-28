# Remote Reliability Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Harden the remote controller/driver protocol with on-demand lease expiry cleanup, idempotent mission reporting, and heartbeat-timeout health degradation.

**Architecture:** Keep the protocol shape intact and harden behavior at controller decision points. Reclaim expired leases during claim and heartbeat flows, track per-mission batch ids for event/sample idempotence, and refresh driver health lazily at claim/read time instead of adding background sweepers.

**Tech Stack:** Go 1.26, existing `internal/controlplane`, `internal/web`, `internal/driver/agent`, file-backed snapshot persistence

---

### Task 1: Add Failing Lease-Expiry Tests

**Files:**
- Modify: `internal/controlplane/mission_scheduler_test.go`
- Modify: `internal/driver/agent/agent_test.go`

- [ ] **Step 1: Write failing tests for expired lease requeue**

Add tests proving:
- an expired claimed mission is marked `expired` and can be claimed by a different driver
- an unexpired claimed mission is not reassigned
- combined loopback processing triggers the same expiry cleanup path

- [ ] **Step 2: Run focused control-plane and agent tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/driver/agent
```

Expected:
- failures because expired leases are not reclaimed today

### Task 2: Implement On-Demand Lease Expiry Cleanup

**Files:**
- Modify: `internal/controlplane/mission_scheduler.go`
- Modify: `internal/controlplane/driver_registry.go`
- Modify: `internal/driver/agent/agent.go`

- [ ] **Step 1: Add a controller helper that expires stale mission leases**

Reclaim missions in `claimed` or `running` when `Lease.ExpiresAt` is older than `now`.

- [ ] **Step 2: Call that helper from claim, heartbeat, and combined loopback entry points**

Keep the behavior deterministic and synchronous.

- [ ] **Step 3: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/driver/agent
```

Expected:
- lease-expiry tests pass

- [ ] **Step 4: Commit the lease-expiry slice**

Run:
```bash
git add internal/controlplane/mission_scheduler.go internal/controlplane/driver_registry.go internal/driver/agent/agent.go internal/controlplane/mission_scheduler_test.go internal/driver/agent/agent_test.go
git commit -m "feat: reclaim expired mission leases"
```

### Task 3: Add Failing Idempotence Tests

**Files:**
- Modify: `internal/web/driver_api_test.go`
- Modify: `internal/controlplane/mission_scheduler_test.go`

- [ ] **Step 1: Write failing tests for event/sample/complete replay**

Add tests proving:
- the same event batch is accepted once
- the same sample batch is accepted once
- repeating completion does not duplicate aggregation

- [ ] **Step 2: Run focused control-plane and web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- failures because duplicate reporting currently double-applies

### Task 4: Implement Mission-Local Idempotence

**Files:**
- Modify: `internal/domain/remote.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/controlplane/mission_scheduler.go`
- Modify: `internal/web/driver_api.go`
- Modify: `internal/driver/agent/http_client.go`
- Modify: `internal/snapshot/store.go`

- [ ] **Step 1: Add mission-local batch tracking**

Track event and sample batch ids per mission, plus a terminal completion marker.

- [ ] **Step 2: Extend reporting payloads with batch ids**

Keep the existing endpoints but add `batch_id` to event and sample upload payloads.

- [ ] **Step 3: Make event/sample/complete handlers idempotent**

Duplicate payloads should return success without double-applying state.

- [ ] **Step 4: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- replay/idempotence coverage passes

- [ ] **Step 5: Commit the idempotence slice**

Run:
```bash
git add internal/domain/remote.go internal/controlplane/manager.go internal/controlplane/mission_scheduler.go internal/web/driver_api.go internal/driver/agent/http_client.go internal/snapshot/store.go internal/web/driver_api_test.go internal/controlplane/mission_scheduler_test.go
git commit -m "feat: make mission reporting idempotent"
```

### Task 5: Add Failing Heartbeat Timeout Tests

**Files:**
- Modify: `internal/controlplane/mission_scheduler_test.go`
- Modify: `internal/controlplane/driver_readmodel_test.go`

- [ ] **Step 1: Write failing tests for stale-driver health degradation**

Add tests proving:
- a driver with stale heartbeat becomes `unhealthy`
- unhealthy drivers cannot claim new missions
- healthy drivers still can

- [ ] **Step 2: Run focused control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failures because driver health is not refreshed from heartbeat age

### Task 6: Implement Heartbeat Timeout Refresh

**Files:**
- Modify: `internal/controlplane/driver_registry.go`
- Modify: `internal/controlplane/driver_readmodel.go`
- Modify: `internal/controlplane/mission_scheduler.go`

- [ ] **Step 1: Add a stale-heartbeat health refresh helper**

Mark drivers `unhealthy` when `LastHeartbeatAt` is older than the timeout window.

- [ ] **Step 2: Invoke health refresh from read and claim paths**

At minimum:
- `ListDriverNodes`
- `GetDriverOverview`
- `ClaimMission`

- [ ] **Step 3: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- heartbeat-timeout coverage passes

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

- [ ] **Step 6: Commit the heartbeat-hardening slice**

Run:
```bash
git add internal/controlplane/driver_registry.go internal/controlplane/driver_readmodel.go internal/controlplane/mission_scheduler.go internal/controlplane/driver_readmodel_test.go internal/controlplane/mission_scheduler_test.go
git commit -m "feat: degrade stale drivers to unhealthy"
```
