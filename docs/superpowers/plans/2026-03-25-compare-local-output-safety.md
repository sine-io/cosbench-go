# Compare Local Output Safety Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make compare-local output cleanup safe even when the output directory is overridden.

**Architecture:** Add a failing test for unsafe output directories, then make the Makefile reject unsafe paths and clear only the contents of a dedicated `compare-local` directory.

**Tech Stack:** Go integration tests, Makefile, Markdown docs

---

### Task 1: Add a Failing Safety Test

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`
- Modify: `Makefile`

- [ ] **Step 1: Write the failing test**

Add coverage that runs `make compare-local` with an unsafe output directory and expects failure.

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- the unsafe-directory test fails because compare-local still accepts the path

- [ ] **Step 3: Implement the minimal Makefile guard**

Reject unsafe output directories and switch cleanup to directory-content pruning.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- stale-output and unsafe-directory tests both pass

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the dedicated output directory rule**

Explain that overrides should still target a directory whose basename is `compare-local`.

### Task 3: Final Verification

- [ ] **Step 1: Run targeted CLI tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

- [ ] **Step 2: Run compare-local**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local
```

- [ ] **Step 3: Run full verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```
