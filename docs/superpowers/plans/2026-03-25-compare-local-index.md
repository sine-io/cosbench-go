# Compare Local Index Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a stable top-level index file for compare-local artifacts.

**Architecture:** Extend the existing compare-local Make target to emit `index.json` beside the per-fixture summaries, and add an integration test that validates the generated index.

**Tech Stack:** Go integration tests, Makefile, JSON, Markdown docs

---

### Task 1: Add a Failing Index Test

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`
- Modify: `Makefile`

- [ ] **Step 1: Write the failing test**

Add coverage that runs `make compare-local` into a temp `compare-local` directory and expects `index.json` to exist with entries for the curated fixtures.

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- the new index assertion fails because no `index.json` is generated yet

- [ ] **Step 3: Implement the minimal index generation**

Generate `index.json` from the existing fixture manifest during the Makefile loop.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- targeted tests pass, including the new index coverage

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the index**

Explain that `index.json` is the top-level artifact entrypoint for compare-local outputs.

### Task 3: Final Verification

- [ ] **Step 1: Run targeted tests**

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
