# Compare Local Structured Output Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Produce stable structured JSON outputs for each compare-local fixture run.

**Architecture:** Extend the local CLI with a `-summary-file` option, then reuse that option from `make compare-local` so both local runs and the manual GitHub workflow can collect a directory of per-fixture JSON summaries.

**Tech Stack:** Go CLI code, Go tests, Makefile, GitHub Actions YAML, Markdown docs

---

### Task 1: Add a Failing CLI Test for Summary Files

**Files:**
- Modify: `cmd/cosbench-go/main_test.go`
- Modify: `cmd/cosbench-go/main.go`

- [ ] **Step 1: Write a failing test**

Add coverage for:
- parsing `-summary-file`
- writing a summary JSON file during `runCLI`

- [ ] **Step 2: Run the targeted test and verify it fails**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- the new test fails because `-summary-file` is not implemented yet

- [ ] **Step 3: Implement the minimal CLI support**

Add argument parsing and file-writing for the summary payload, keeping existing stdout behavior intact.

- [ ] **Step 4: Re-run the targeted test**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- the new summary-file coverage passes

### Task 2: Update compare-local Output Collection

**Files:**
- Modify: `Makefile`
- Modify: `.github/workflows/compare-local.yml`

- [ ] **Step 1: Write structured outputs**

Change `compare-local` so it creates a stable output directory and writes one JSON file per fixture there.

- [ ] **Step 2: Upload the structured directory**

Update the manual workflow to upload the directory instead of a single text file.

### Task 3: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the output path**

Explain where local and workflow-driven compare-local JSON summaries land and how they relate to the comparison matrix.

### Task 4: Final Verification

**Files:**
- Review only: `cmd/cosbench-go/main.go`
- Review only: `cmd/cosbench-go/main_test.go`
- Review only: `Makefile`
- Review only: `.github/workflows/compare-local.yml`

- [ ] **Step 1: Run targeted CLI tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- CLI tests pass, including summary-file coverage

- [ ] **Step 2: Run compare-local**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local
```

Expected:
- compare-local completes
- `.artifacts/compare-local/` contains one JSON file per fixture

- [ ] **Step 3: Run full verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- full test suite passes
- build passes
