# Storage Config Fallback Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make storage-level `part_size` and `restore_days` act as execution-time defaults while preserving explicit operation-level overrides.

**Architecture:** Keep the fallback logic inside the execution-config path. Merge storage config KV pairs first, overlay operation config KV pairs second, and use the merged view consistently in runtime execution and preflight validation.

**Tech Stack:** Go 1.26, package-local tests in `internal/domain/execution` and `internal/controlplane`, existing storage adapter port

---

### Task 1: Add Failing Execution Tests

**Files:**
- Modify: `internal/domain/execution/engine_test.go`
- Modify: `internal/domain/execution/opconfig_test.go`

- [ ] **Step 1: Write failing tests for storage-default and op-override behavior**

Add tests proving:
- storage-level `part_size` drives `MultipartPut` when op config omits it
- op-level `part_size` overrides the storage-level value
- storage-level `restore_days` drives `RestoreObject` when op config omits it
- op-level `restore_days` overrides the storage-level value

- [ ] **Step 2: Run execution tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution
```

Expected:
- failures because runtime execution currently only uses operation-config values

### Task 2: Implement Merged Execution Config

**Files:**
- Modify: `internal/domain/execution/opconfig.go`
- Modify: `internal/domain/execution/engine.go`
- Modify: `internal/executor/executor.go`

- [ ] **Step 1: Add a merge helper for storage and operation config**

Implement a helper that:
- parses storage config first
- overlays operation config second
- returns the same parsed shape as current op config handling

- [ ] **Step 2: Switch runtime execution to the merged config path**

Use the merged parse for:
- `mwrite`
- `mfilewrite`
- `restore`

- [ ] **Step 3: Keep preflight aligned**

Ensure validation paths use the same merged view rather than reintroducing the old behavior in preflight.

- [ ] **Step 4: Re-run execution tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution
```

Expected:
- execution package passes with fallback and override behavior covered

### Task 3: Add Preflight Regression Coverage

**Files:**
- Modify: `internal/controlplane/manager_test.go`

- [ ] **Step 1: Write a focused preflight regression test**

Add a test showing a storage-level-only `part_size` or `restore_days` config passes preflight and aligns with runtime expectations.

- [ ] **Step 2: Run control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- control-plane package passes with preflight still aligned

### Task 4: Final Verification

**Files:**
- Review only: `internal/domain/execution/opconfig.go`
- Review only: `internal/domain/execution/engine.go`
- Review only: `internal/executor/executor.go`
- Review only: `internal/domain/execution/engine_test.go`
- Review only: `internal/controlplane/manager_test.go`

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

- [ ] **Step 3: Review scope**

Run:
```bash
git diff -- internal/domain/execution internal/executor internal/controlplane docs
```

Expected:
- the slice stays focused on config fallback and aligned tests
