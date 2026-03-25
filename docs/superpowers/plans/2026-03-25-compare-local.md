# Compare Local Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a local comparison convenience command that runs the curated mock-backed fixture set through the existing CLI.

**Architecture:** Keep the feature thin: a single `make compare-local` wrapper around existing `go run ./cmd/cosbench-go ... -backend mock -json` invocations. Update docs to point at the command and, if useful, refresh matrix notes based on the observed outputs.

**Tech Stack:** Makefile, existing CLI runner, Markdown docs

---

### Task 1: Add the Compare Command

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Add a failing expectation check**

Run:
```bash
make -n compare-local
```

Expected:
- target missing before implementation

- [ ] **Step 2: Add `compare-local`**

The target should run the current CLI against:
- `testdata/workloads/s3-active-subset.xml`
- `testdata/workloads/mock-stage-aware.xml`
- `testdata/workloads/mock-reusedata-subset.xml`
- `testdata/workloads/xml-splitrw-subset.xml`

using:
- `GO=...`
- `go run ./cmd/cosbench-go ... -backend mock -json`

- [ ] **Step 3: Verify command expansion**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make -n compare-local
```

Expected:
- the target expands to the expected fixture runs

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the local comparison command**

Mention:
- `make compare-local`
- that it is mock-backed and local-only
- that live endpoint comparison is still separate

- [ ] **Step 2: Optionally refresh comparison matrix wording**

If the command is run successfully, update the matrix or runbook to mention it as the fastest local refresh path.

### Task 3: Final Verification

**Files:**
- Review only: `Makefile`
- Review only: `README.md`
- Review only: `AGENTS.md`
- Review only: `docs/legacy-comparison-matrix.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Run the new compare command**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local
```

Expected:
- the curated mock-backed fixture set runs successfully

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

- [ ] **Step 4: Review scope**

Run:
```bash
git diff -- Makefile README.md AGENTS.md docs/legacy-comparison-matrix.md BOARD.md TODO.md
```

Expected:
- the slice stays automation/doc focused
