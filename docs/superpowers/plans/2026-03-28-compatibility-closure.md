# Compatibility Closure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Close the remaining S3/SIO-family compatibility gaps so parser-covered and alias-only behaviors become real runtime behavior in `cosbench-go`.

**Architecture:** Extend the workload model with explicit auth and backend profile concepts, route all execution through one effective-config resolver, and close the remaining operation and adapter deltas in the existing execution and S3/SIO adapter path.

**Tech Stack:** Go 1.26, existing `internal/domain`, `internal/infrastructure/xml`, `internal/domain/execution`, `internal/driver/s3`, file-backed fixtures under `testdata/`

---

### Task 1: Add Auth And Compatibility Parser Coverage

**Files:**
- Create: `internal/domain/auth.go`
- Modify: `internal/domain/workload.go`
- Modify: `internal/infrastructure/xml/workload_parser.go`
- Modify: `internal/infrastructure/xml/workload_parser_test.go`
- Modify: `internal/workloadxml/parser_test.go`
- Create: `testdata/workloads/xml-auth-inheritance-subset.xml`

- [ ] **Step 1: Write failing parser tests for explicit auth nodes and inheritance**

Add coverage proving:
- workload/stage/work auth nodes survive parsing
- inherited auth material is preserved in the normalized model
- existing auth-tolerated fixtures still parse

- [ ] **Step 2: Run XML parser tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml ./internal/workloadxml
```

Expected:
- failures because auth nodes are currently tolerated but not modeled

- [ ] **Step 3: Implement `AuthSpec` and XML/domain mapping**

Add the minimal domain and parser plumbing needed to preserve auth blocks through parsing and normalization.

- [ ] **Step 4: Re-run parser tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml ./internal/workloadxml
```

Expected:
- parser packages pass with explicit auth coverage

- [ ] **Step 5: Commit the parser/model slice**

Run:
```bash
git add internal/domain/auth.go internal/domain/workload.go internal/infrastructure/xml/workload_parser.go internal/infrastructure/xml/workload_parser_test.go internal/workloadxml/parser_test.go testdata/workloads/xml-auth-inheritance-subset.xml
git commit -m "feat: model workload auth specs"
```

### Task 2: Implement Runtime Config Resolution And `filewrite`

**Files:**
- Modify: `internal/domain/workload/normalize.go`
- Modify: `internal/domain/execution/opconfig.go`
- Modify: `internal/domain/execution/engine.go`
- Modify: `internal/domain/execution/engine_test.go`
- Modify: `internal/executor/executor.go`
- Modify: `internal/controlplane/manager_test.go`
- Create: `testdata/workloads/xml-filewrite-subset.xml`

- [ ] **Step 1: Write failing execution tests for resolved auth/config behavior and `filewrite`**

Add tests proving:
- auth and storage config resolve into one effective execution view
- `filewrite` reads a selected local file and uploads it as a single-part write
- unreadable file-backed inputs fail preflight and runtime consistently

- [ ] **Step 2: Run execution and control-plane tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution ./internal/controlplane
```

Expected:
- failures because there is no runtime `filewrite` path and no shared auth-aware config resolution

- [ ] **Step 3: Implement one effective runtime config resolver**

Resolve runtime config from:
- endpoint material
- storage config
- inherited auth
- operation config

Make preflight and execution use the same resolved view.

- [ ] **Step 4: Implement `filewrite` using the existing file-backed operation path**

Extend the execution layer so `filewrite`:
- selects a local file
- opens the file
- uploads with `PutObject`
- records actual byte count

- [ ] **Step 5: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution ./internal/controlplane
```

Expected:
- focused execution and control-plane coverage passes

- [ ] **Step 6: Commit the execution slice**

Run:
```bash
git add internal/domain/workload/normalize.go internal/domain/execution/opconfig.go internal/domain/execution/engine.go internal/domain/execution/engine_test.go internal/executor/executor.go internal/controlplane/manager_test.go testdata/workloads/xml-filewrite-subset.xml
git commit -m "feat: add filewrite and resolved auth config"
```

### Task 3: Add Backend Profile Policies And Read-Side Compatibility

**Files:**
- Modify: `internal/infrastructure/storage/factory.go`
- Modify: `internal/driver/s3/config.go`
- Modify: `internal/driver/s3/adapter.go`
- Modify: `internal/driver/s3/config_test.go`
- Modify: `internal/driver/s3/smoke_test.go`
- Modify: `internal/domain/execution/engine_test.go`
- Modify: `testdata/workloads/xml-range-prefetch-subset.xml`

- [ ] **Step 1: Write failing adapter/profile tests for `sio`, `siov1`, `gdas`, prefetch, and range reads**

Add tests proving:
- runtime profile selection is explicit for `sio`, `siov1`, and `gdas`
- prefetch config changes request shaping
- range-read config issues bounded reads
- documented defaulting and delete/list behaviors are either closed or deliberately codified

- [ ] **Step 2: Run driver tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/s3 ./internal/domain/execution
```

Expected:
- failures because the current adapter treats aliases loosely and only preserves prefetch/range config strings

- [ ] **Step 3: Implement explicit backend profiles and read-side request shaping**

Keep transport code shared where possible, but make behavior policy explicit by backend kind and read-mode flags.

- [ ] **Step 4: Close or codify the documented S3/SIO deltas**

Update tests and implementation so the known differences are either:
- fixed in code
- or converted into stable documented policy with tests

- [ ] **Step 5: Re-run focused driver tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/s3 ./internal/domain/execution
```

Expected:
- driver and execution tests pass with profile/read compatibility covered

- [ ] **Step 6: Commit the adapter slice**

Run:
```bash
git add internal/infrastructure/storage/factory.go internal/driver/s3/config.go internal/driver/s3/adapter.go internal/driver/s3/config_test.go internal/driver/s3/smoke_test.go internal/domain/execution/engine_test.go testdata/workloads/xml-range-prefetch-subset.xml
git commit -m "feat: close sio profile and read compatibility gaps"
```

### Task 4: Refresh Fixtures, Comparison Docs, And Full Verification

**Files:**
- Modify: `docs/xml-compat-matrix.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `docs/storage-driver-comparison-notes.md`
- Modify: `README.md`

- [ ] **Step 1: Update compatibility docs to reflect the new runtime coverage**

Mark which items moved from parser-only or deferred into supported runtime behavior, and note any remaining intentional policy differences.

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

- [ ] **Step 4: Review final scope**

Run:
```bash
git diff -- internal/domain internal/infrastructure/xml internal/domain/execution internal/executor internal/driver/s3 docs README.md testdata/workloads
```

Expected:
- the slice stays focused on compatibility closure inside the S3/SIO family

- [ ] **Step 5: Commit the documentation and verification updates**

Run:
```bash
git add docs/xml-compat-matrix.md docs/legacy-comparison-matrix.md docs/storage-driver-comparison-notes.md README.md
git commit -m "docs: record compatibility closure state"
```
