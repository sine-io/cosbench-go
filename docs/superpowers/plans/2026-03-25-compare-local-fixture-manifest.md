# Compare Local Fixture Manifest Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the compare-local fixture set come from one manifest file instead of duplicated lists.

**Architecture:** Add a small manifest under `testdata/workloads/`, update the Makefile to iterate over it, and add a test that validates each listed workload path can still be parsed.

**Tech Stack:** Go tests, Makefile, plain-text manifest, Markdown docs

---

### Task 1: Add Failing Manifest Coverage

**Files:**
- Create: `cmd/cosbench-go/compare_local_manifest_test.go`
- Create: `testdata/workloads/compare-local-fixtures.txt`

- [ ] **Step 1: Write a failing test**

Add coverage that loads the compare-local manifest and validates:
- each row has an output name and workload path
- each workload file exists
- each workload parses successfully

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- the new test fails because the manifest reader or manifest file does not exist yet

- [ ] **Step 3: Add the manifest and minimal loader**

Create the manifest and test-only loader needed to validate it.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- manifest coverage passes

### Task 2: Move compare-local to the Manifest

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Read the manifest in compare-local**

Replace the duplicated inline fixture list with a simple loop over the manifest.

### Task 3: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the manifest**

Explain where the curated compare-local fixture set is defined.

### Task 4: Final Verification

**Files:**
- Review only: `Makefile`
- Review only: `cmd/cosbench-go/compare_local_manifest_test.go`
- Review only: `testdata/workloads/compare-local-fixtures.txt`

- [ ] **Step 1: Run targeted CLI tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- targeted tests pass, including manifest coverage

- [ ] **Step 2: Run compare-local**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local
```

Expected:
- compare-local completes from the manifest-driven fixture set

- [ ] **Step 3: Run full verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- full verification remains green
