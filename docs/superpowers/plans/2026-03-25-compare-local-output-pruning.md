# Compare Local Output Pruning Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ensure compare-local output directories contain only fresh results from the current fixture manifest.

**Architecture:** Add a regression test that runs the Make target against a temp directory, then update the Makefile to recreate that directory before writing result files.

**Tech Stack:** Go integration tests, Makefile, Markdown docs

---

### Task 1: Add a Failing Regression Test

**Files:**
- Create: `cmd/cosbench-go/compare_local_make_test.go`
- Modify: `Makefile`

- [ ] **Step 1: Write the failing test**

Add a test that seeds a stale file into a temp output directory, runs `make compare-local` with `COMPARE_LOCAL_OUTPUT_DIR` overridden, and expects the stale file to disappear.

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- the stale file still exists after the run

- [ ] **Step 3: Implement the minimal Makefile change**

Recreate the compare-local output directory before generating fresh JSON summaries.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- the stale-file regression test passes

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the refresh behavior**

Explain that `make compare-local` refreshes its output directory on each run.

### Task 3: Final Verification

**Files:**
- Review only: `Makefile`
- Review only: `cmd/cosbench-go/compare_local_make_test.go`

- [ ] **Step 1: Run targeted CLI tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- targeted tests pass, including the stale-output regression test

- [ ] **Step 2: Run compare-local**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local
```

Expected:
- compare-local completes
- output directory contains only fresh files from the current run

- [ ] **Step 3: Run full verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- full verification remains green
