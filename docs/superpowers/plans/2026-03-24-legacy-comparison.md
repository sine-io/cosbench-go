# Legacy Comparison Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a repeatable, repository-local legacy comparison matrix and runbook for the representative workload subset without building a full automation harness.

**Architecture:** Keep this slice documentation-first. Create a single comparison matrix document that names the representative fixtures, defines comparison dimensions, records current known status, and embeds a minimal runbook. Update repository-facing docs and status tracking to point contributors at that matrix instead of leaving the comparison process implicit.

**Tech Stack:** Markdown documentation under `docs/`, existing CLI and smoke-test commands, repository testdata references

---

### Task 1: Create the Legacy Comparison Matrix

**Files:**
- Create: `docs/legacy-comparison-matrix.md`

- [ ] **Step 1: Draft the comparison matrix structure**

Include:
- representative fixture set
- comparison dimensions
- comparison result states (`match`, `acceptable delta`, `mismatch`, `not yet run`)
- notes/follow-up column

- [ ] **Step 2: Seed the matrix with current known findings**

Populate rows for:
- `testdata/legacy/s3-config-sample.xml`
- `testdata/legacy/sio-config-sample.xml`
- `testdata/workloads/s3-active-subset.xml`
- `testdata/workloads/sio-multipart-subset.xml`
- `testdata/workloads/xml-inheritance-subset.xml`
- `testdata/workloads/xml-attribute-subset.xml`
- `testdata/workloads/xml-special-ops-subset.xml`

Use explicit statuses such as:
- parser-only comparison
- runnable with live endpoint setup
- not yet run against legacy
- acceptable delta due to scope boundary

- [ ] **Step 3: Add the runbook section**

Document:
- how to run `cosbench-go` locally on a fixture
- how to run `make smoke-s3`
- where to look in `../cosbench-sineio` for legacy references and sample configs

### Task 2: Link the Matrix from Repository Docs

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Link README to the comparison matrix**

Add a short section or bullet pointing contributors to `docs/legacy-comparison-matrix.md` for current known deltas and repeatable comparison steps.

- [ ] **Step 2: Reference the matrix from migration-gap analysis**

Update the gap analysis so “real-endpoint shakeout” and remaining risks point readers to the comparison matrix rather than leaving comparison work implicit.

### Task 3: Update Board and Checklist State

**Files:**
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Reflect that comparison infrastructure has landed**

Update:
- `BOARD.md` to show the matrix/runbook as landed
- `TODO.md` to move “Compare benchmark behavior against legacy workloads” from not-started to in-progress

- [ ] **Step 2: Keep live comparison work visible**

Do not mark the overall legacy comparison effort fully complete unless the matrix proves that. Keep remaining live-endpoint comparison work visible.

### Task 4: Final Verification

**Files:**
- Review only: `docs/legacy-comparison-matrix.md`
- Review only: `README.md`
- Review only: `docs/migration-gap-analysis.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Verify the new document exists and is readable**

Run:
```bash
sed -n '1,260p' docs/legacy-comparison-matrix.md
```

Expected:
- the matrix and runbook are present and coherent

- [ ] **Step 2: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass unchanged

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
git diff -- docs/legacy-comparison-matrix.md README.md docs/migration-gap-analysis.md BOARD.md TODO.md
```

Expected:
- this slice remains documentation/status focused
