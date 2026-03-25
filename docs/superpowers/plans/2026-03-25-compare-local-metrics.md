# Compare Local Metrics Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enrich compare-local index and workflow summary output with per-fixture metrics.

**Architecture:** Extend the compare-local integration test to require metric fields in `index.json`, then switch index generation to a small Python step that reads the manifest and summary files, and finally update the GitHub workflow summary table to show those fields.

**Tech Stack:** Go integration tests, Makefile, inline Python, GitHub Actions YAML, Markdown docs

---

### Task 1: Add a Failing Metric Assertion

**Files:**
- Modify: `cmd/cosbench-go/compare_local_make_test.go`
- Modify: `Makefile`

- [ ] **Step 1: Write the failing test**

Assert that each fixture entry in `index.json` includes `stages`, `works`, `samples`, and `errors`.

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

- [ ] **Step 3: Implement enriched index generation**

Generate `index.json` from the manifest and per-fixture summary files after the compare-local loop completes.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

### Task 2: Update the Workflow Summary

**Files:**
- Modify: `.github/workflows/compare-local.yml`

- [ ] **Step 1: Show metrics in the table**

Render `stages`, `works`, `samples`, and `errors` from `index.json` into the GitHub job summary.

### Task 3: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the enriched index**

Explain that `index.json` and the workflow summary now carry per-fixture metrics.

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
```

- [ ] **Step 3: Run full verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```
