# Deferred Parser Coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add parser-facing coverage for deferred-but-adjacent XML shapes: compatibility storage aliases and range/prefetch config keys.

**Architecture:** Add two focused fixtures under `testdata/workloads/` and extend the two parser test layers to assert raw and normalized preservation of those shapes. Keep the slice strictly parser-facing: no execution tests, no new runtime support, no changes to deferred status.

**Tech Stack:** Go 1.26, package-local tests in `internal/infrastructure/xml` and `internal/workloadxml`, XML fixtures under `testdata/workloads`

---

### Task 1: Add Deferred Parser Fixtures

**Files:**
- Create: `testdata/workloads/xml-compat-storage-subset.xml`
- Create: `testdata/workloads/xml-range-prefetch-subset.xml`

- [ ] **Step 1: Create the fixtures**

Add:
- `xml-compat-storage-subset.xml` for `siov1` / `gdas` compatibility storage aliases
- `xml-range-prefetch-subset.xml` for prefetch/range-read config keys

- [ ] **Step 2: Verify fixture readability**

Run:
```bash
sed -n '1,220p' testdata/workloads/xml-compat-storage-subset.xml
sed -n '1,220p' testdata/workloads/xml-range-prefetch-subset.xml
```

Expected:
- fixtures clearly express parser-only compatibility coverage

### Task 2: Add Raw Parser Coverage

**Files:**
- Modify: `internal/infrastructure/xml/workload_parser_test.go`

- [ ] **Step 1: Write failing raw-parser tests**

Add assertions for:
- `siov1` / `gdas` storage types are preserved in parsed workload storage specs
- prefetch/range-read config keys remain present in raw config strings

- [ ] **Step 2: Run raw-parser tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml
```

Expected:
- failures because the fixtures and new assertions are not yet in place or aligned

- [ ] **Step 3: Adjust only fixtures/test expectations**

Do not add runtime support in this task.

- [ ] **Step 4: Re-run raw-parser tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml
```

Expected:
- raw parser package passes

### Task 3: Add Normalized Parser Coverage

**Files:**
- Modify: `internal/workloadxml/parser_test.go`

- [ ] **Step 1: Write failing normalized-parser tests**

Add assertions for:
- compatibility storage aliases remain visible after normalization
- prefetch/range config keys remain visible in normalized config strings

- [ ] **Step 2: Run normalized-parser tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/workloadxml
```

Expected:
- failures because the fixtures/assertions are not fully aligned yet

- [ ] **Step 3: Fix only fixture/test expectation mismatches**

Do not expand execution behavior or deferred support here.

- [ ] **Step 4: Re-run normalized-parser tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/workloadxml
```

Expected:
- normalized parser package passes

### Task 4: Update Compatibility Documentation

**Files:**
- Modify: `docs/xml-compat-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Update the compatibility matrix**

Reference the new fixtures and clarify that they prove parser coverage while runtime support remains deferred.

- [ ] **Step 2: Update board/checklist state**

Reflect that parser-facing coverage for deferred XML constructs has advanced.

### Task 5: Final Verification

**Files:**
- Review only: `testdata/workloads/xml-compat-storage-subset.xml`
- Review only: `testdata/workloads/xml-range-prefetch-subset.xml`
- Review only: `internal/infrastructure/xml/workload_parser_test.go`
- Review only: `internal/workloadxml/parser_test.go`
- Review only: `docs/xml-compat-matrix.md`

- [ ] **Step 1: Run focused parser verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml ./internal/workloadxml
```

Expected:
- both parser packages pass

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
git diff -- testdata/workloads internal/infrastructure/xml/workload_parser_test.go internal/workloadxml/parser_test.go docs/xml-compat-matrix.md BOARD.md TODO.md
```

Expected:
- the slice stays parser/test/doc focused
