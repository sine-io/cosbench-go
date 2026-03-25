# Compare Local Filter Validation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make invalid compare-local filters fail clearly instead of silently producing empty results.

**Architecture:** Add a failing integration test for an unknown filter, then validate the requested fixture name against the compare-local manifest before running the target loop.

**Tech Stack:** Go integration tests, Makefile, Markdown docs

---

### Task 1: Add a Failing Invalid-Filter Test

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`
- Modify: `Makefile`

- [ ] **Step 1: Write the failing test**

Assert that `COMPARE_LOCAL_FILTER=does-not-exist` exits non-zero and mentions the invalid fixture name.

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

- [ ] **Step 3: Implement filter validation**

Reject unknown fixture names before the compare-local loop starts.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document validation**

Mention that filtered runs only accept names from the compare-local manifest.

### Task 3: Final Verification

- [ ] **Step 1: Run targeted tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

- [ ] **Step 2: Run compare-local**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local COMPARE_LOCAL_FILTER=mock-stage-aware
```

- [ ] **Step 3: Run full verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```
