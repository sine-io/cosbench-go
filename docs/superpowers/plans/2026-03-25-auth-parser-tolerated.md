# Auth Parser-Tolerated Coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add parser-facing coverage proving that auth-bearing XML shapes are tolerated by the current parser/normalization path without introducing auth modeling.

**Architecture:** Add two small fixtures under `testdata/workloads/` and extend the raw and normalized parser tests to show that `<auth>` nodes do not break parsing and do not appear in the current domain model. Update the compatibility matrix to distinguish parser tolerance from auth modeling.

**Tech Stack:** Go 1.26, package-local tests in `internal/infrastructure/xml` and `internal/workloadxml`, XML fixtures under `testdata/workloads`

---

### Task 1: Add Auth-Bearing Fixtures

**Files:**
- Create: `testdata/workloads/xml-auth-tolerated-subset.xml`
- Create: `testdata/workloads/xml-auth-none-subset.xml`

- [ ] **Step 1: Create the fixtures**

Add:
- `xml-auth-tolerated-subset.xml` with workload/stage/work-level auth nodes and supported `mock` storage/work shapes
- `xml-auth-none-subset.xml` with `type="none"` auth nodes and supported surrounding structure

- [ ] **Step 2: Verify fixture readability**

Run:
```bash
sed -n '1,220p' testdata/workloads/xml-auth-tolerated-subset.xml
sed -n '1,220p' testdata/workloads/xml-auth-none-subset.xml
```

Expected:
- fixtures are small, focused, and clearly parser-tolerant rather than auth-modeled

### Task 2: Add Raw Parser Coverage

**Files:**
- Modify: `internal/infrastructure/xml/workload_parser_test.go`

- [ ] **Step 1: Write failing raw-parser tests**

Add assertions proving:
- auth-bearing XML parses successfully
- surrounding workload/stage/work/storage/operation structure is preserved

- [ ] **Step 2: Run raw-parser tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml
```

Expected:
- failures because fixtures/assertions are not fully aligned yet

- [ ] **Step 3: Fix only fixture/test expectation mismatches**

Do not add auth modeling or runtime support in this task.

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

Add assertions proving:
- auth-bearing XML still normalizes into the expected current domain structure
- auth presence does not disturb supported storage/config/work/operation parsing

- [ ] **Step 2: Run normalized-parser tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/workloadxml
```

Expected:
- failures because fixtures/assertions are not fully aligned yet

- [ ] **Step 3: Fix only fixture/test expectation mismatches**

Do not expand the domain model for auth.

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

- [ ] **Step 1: Update the auth row in the compatibility matrix**

Clarify:
- parser-tolerated auth-bearing XML is covered
- explicit auth modeling remains deferred

- [ ] **Step 2: Update board/checklist state**

Reflect that parser-facing auth coverage has advanced without changing runtime support.

### Task 5: Final Verification

**Files:**
- Review only: `testdata/workloads/xml-auth-tolerated-subset.xml`
- Review only: `testdata/workloads/xml-auth-none-subset.xml`
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
- the slice stays parser/test/doc focused and does not introduce auth modeling
