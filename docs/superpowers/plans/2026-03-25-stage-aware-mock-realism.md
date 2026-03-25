# Stage-Aware Mock Realism Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make local `mock`-backed multi-stage workload runs preserve object state across works and stages within a single run or job.

**Architecture:** Reuse one `mock` adapter per local run context instead of constructing a fresh adapter per work. Apply that reuse in both the control-plane job runner and the CLI local runner, while leaving real S3/SIO adapter lifetime behavior unchanged.

**Tech Stack:** Go 1.26, package-local tests in `internal/controlplane` and `cmd/cosbench-go`, representative XML fixtures under `testdata/workloads`

---

### Task 1: Add a Representative Mock Stage-Aware Fixture

**Files:**
- Create: `testdata/workloads/mock-stage-aware.xml`

- [ ] **Step 1: Create the fixture**

The fixture should cover:
- `init`
- `prepare`
- `read`
- `list`
- `cleanup`
- `dispose`

- [ ] **Step 2: Verify fixture readability**

Run:
```bash
sed -n '1,220p' testdata/workloads/mock-stage-aware.xml
```

Expected:
- the XML clearly expresses stage-to-stage object continuity under the `mock` backend

### Task 2: Add Failing Control-Plane Coverage

**Files:**
- Modify: `internal/controlplane/manager_test.go`

- [ ] **Step 1: Write a failing test for multi-stage mock realism**

Add a test proving:
- a `mock` job using `prepare -> read -> list -> cleanup` succeeds under the control plane
- object state survives across stage boundaries within one job

- [ ] **Step 2: Run the control-plane tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failure because the current control-plane path creates a fresh mock adapter per work

### Task 3: Add Failing CLI Coverage

**Files:**
- Modify: `cmd/cosbench-go/main.go`
- Create: `cmd/cosbench-go/main_test.go`

- [ ] **Step 1: Write a failing CLI-oriented test**

Refactor target:
- a testable helper for local workload execution

Assertion:
- the CLI local runner succeeds on `mock-stage-aware.xml` with shared mock state

- [ ] **Step 2: Run the CLI package tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- failure because there is no helper or shared mock adapter reuse yet

### Task 4: Implement Mock Adapter Reuse

**Files:**
- Modify: `internal/controlplane/manager.go`
- Modify: `cmd/cosbench-go/main.go`
- Modify: `internal/infrastructure/storage/factory.go` only if a small helper materially reduces duplication

- [ ] **Step 1: Implement per-job mock adapter reuse in the control plane**

Behavior:
- when resolved backend is `mock`, reuse one adapter instance for the whole job
- dispose it once the job completes, fails, or is cancelled

- [ ] **Step 2: Implement per-run mock adapter reuse in the CLI path**

Behavior:
- when effective backend is `mock`, reuse one adapter instance for the whole workload invocation
- dispose it once the CLI run finishes

- [ ] **Step 3: Keep real backends unchanged**

Do not change the per-work lifetime of `s3` / `sio`.

- [ ] **Step 4: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./cmd/cosbench-go
```

Expected:
- both packages pass with shared mock-state behavior covered

### Task 5: Final Verification and Status Update

**Files:**
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Update board/checklist state**

Reflect that stage-aware realism for local mock-backed runs is no longer open.

- [ ] **Step 2: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass

- [ ] **Step 3: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 4: Review scope**

Run:
```bash
git diff -- testdata/workloads/mock-stage-aware.xml internal/controlplane/manager.go internal/controlplane/manager_test.go cmd/cosbench-go/main.go cmd/cosbench-go/main_test.go BOARD.md TODO.md
```

Expected:
- the slice stays limited to local mock realism
