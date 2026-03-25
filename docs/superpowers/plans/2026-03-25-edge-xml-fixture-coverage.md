# Edge XML Fixture Coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add representative edge-case XML fixtures for delay-stage, split read/write targeting, and reuse-data stage structure within the current local-only v1 scope.

**Architecture:** Add three focused fixtures under `testdata/workloads/`, extend parser-facing tests to lock their normalized shape, and add a minimal control-plane proof for the reuse-data path under `mock`. Keep the slice within current supported backends and stage semantics.

**Tech Stack:** Go 1.26, package-local tests in `internal/infrastructure/xml`, `internal/workloadxml`, and `internal/controlplane`, XML fixtures under `testdata/workloads`

---

### Task 1: Add Edge-Case Fixtures

**Files:**
- Create: `testdata/workloads/xml-delay-stage-subset.xml`
- Create: `testdata/workloads/xml-splitrw-subset.xml`
- Create: `testdata/workloads/mock-reusedata-subset.xml`

- [ ] **Step 1: Create the fixtures**

Add:
- `xml-delay-stage-subset.xml` for repeated delay stages and `closuredelay`
- `xml-splitrw-subset.xml` for split read/write target ranges
- `mock-reusedata-subset.xml` for multi-main-stage reuse after prepare

- [ ] **Step 2: Verify fixture readability**

Run:
```bash
sed -n '1,220p' testdata/workloads/xml-delay-stage-subset.xml
sed -n '1,220p' testdata/workloads/xml-splitrw-subset.xml
sed -n '1,220p' testdata/workloads/mock-reusedata-subset.xml
```

Expected:
- each fixture is focused and maps to one legacy sample family

### Task 2: Add Parser-Facing Coverage

**Files:**
- Modify: `internal/infrastructure/xml/workload_parser_test.go`
- Modify: `internal/workloadxml/parser_test.go`

- [ ] **Step 1: Write failing parser tests**

Add assertions for:
- delay-stage structure and `closuredelay`
- split read/write range preservation
- multi-main-stage reuse layout

- [ ] **Step 2: Run parser package tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml ./internal/workloadxml
```

Expected:
- failures because fixtures and assertions are not fully aligned yet

- [ ] **Step 3: Adjust only fixture/test expectations**

Do not add new parser features in this task. Fix fixture/test mismatches only.

- [ ] **Step 4: Re-run parser package tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml ./internal/workloadxml
```

Expected:
- both parser packages pass with the new edge-case coverage

### Task 3: Add Minimal Control-Plane Proof for Reuse Data

**Files:**
- Modify: `internal/controlplane/manager_test.go`

- [ ] **Step 1: Write a failing control-plane test for reuse-data continuity**

Add a test proving:
- `mock-reusedata-subset.xml` succeeds under the local control plane
- multiple main stages can reuse prepared data before cleanup/dispose

- [ ] **Step 2: Run control-plane tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failure if the fixture shape exposes a continuity or stage-structure mismatch

- [ ] **Step 3: Fix only fixture/test expectation mismatches**

If the failure reveals a real bug, fix it minimally; otherwise keep this task coverage-only.

- [ ] **Step 4: Re-run control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- control-plane package passes with reuse-data coverage

### Task 4: Update Compatibility Documentation

**Files:**
- Modify: `docs/xml-compat-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Update the compatibility matrix**

Reference the three new fixtures in the “Real-Workload Tie-In” section and describe what each one proves.

- [ ] **Step 2: Update board/checklist state**

Reflect that representative edge XML fixture expansion advanced with delay-stage, splitrw, and reuse-data shapes.

### Task 5: Final Verification

**Files:**
- Review only: `testdata/workloads/xml-delay-stage-subset.xml`
- Review only: `testdata/workloads/xml-splitrw-subset.xml`
- Review only: `testdata/workloads/mock-reusedata-subset.xml`
- Review only: `internal/infrastructure/xml/workload_parser_test.go`
- Review only: `internal/workloadxml/parser_test.go`
- Review only: `internal/controlplane/manager_test.go`
- Review only: `docs/xml-compat-matrix.md`

- [ ] **Step 1: Run focused verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml ./internal/workloadxml ./internal/controlplane
```

Expected:
- all targeted packages pass

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
git diff -- testdata/workloads internal/infrastructure/xml/workload_parser_test.go internal/workloadxml/parser_test.go internal/controlplane/manager_test.go docs/xml-compat-matrix.md BOARD.md TODO.md
```

Expected:
- this slice remains fixture/test/doc focused
