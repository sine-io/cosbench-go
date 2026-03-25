# Compare Local Filter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add filtered compare-local runs for one curated fixture at a time.

**Architecture:** Extend the compare-local integration test to cover a filtered run, then add a Makefile filter env var and wire the manual workflow to pass a workflow_dispatch input into that filter.

**Tech Stack:** Go integration tests, Makefile, GitHub Actions YAML, Markdown docs

---

### Task 1: Add a Failing Filter Test

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`
- Modify: `Makefile`

- [ ] **Step 1: Write the failing test**

Assert that `COMPARE_LOCAL_FILTER=mock-stage-aware` only produces the matching summary and one index entry.

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

- [ ] **Step 3: Implement filter support**

Update compare-local to skip non-matching manifest rows when the filter is set.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

### Task 2: Update the Manual Workflow

**Files:**
- Modify: `.github/workflows/compare-local.yml`

- [ ] **Step 1: Add a workflow input**

Expose a `fixture` input and pass it into `COMPARE_LOCAL_FILTER`.

### Task 3: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document filtered runs**

Explain how to run one compare-local fixture locally and through the manual workflow.

### Task 4: Final Verification

- [ ] **Step 1: Run targeted tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

- [ ] **Step 2: Run compare-local**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local
GO=$(which go || echo /snap/bin/go) make compare-local COMPARE_LOCAL_FILTER=mock-stage-aware
```

- [ ] **Step 3: Run full verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```
