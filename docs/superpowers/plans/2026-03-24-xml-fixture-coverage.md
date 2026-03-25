# XML Fixture Coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expand parser-facing XML fixture coverage for high-value local-only v1 shapes that are already implemented but under-tested.

**Architecture:** Add three focused fixtures under `testdata/workloads/` and extend two existing parser test layers to validate them. Keep this slice strictly about representative coverage and compatibility-locking; do not introduce new parser or executor behavior.

**Tech Stack:** Go 1.26, package-local Go tests, XML fixtures under `testdata/workloads`, documentation under `docs/`

---

### Task 1: Add High-Value XML Fixtures

**Files:**
- Create: `testdata/workloads/xml-inheritance-subset.xml`
- Create: `testdata/workloads/xml-attribute-subset.xml`
- Create: `testdata/workloads/xml-special-ops-subset.xml`

- [ ] **Step 1: Draft the three fixtures from the approved spec**

Create:
- `xml-inheritance-subset.xml` for config inheritance, storage override, and omitted-ratio defaulting
- `xml-attribute-subset.xml` for `trigger`, `closuredelay`, `interval`, `division`, `rampup`, `rampdown`, and `driver`
- `xml-special-ops-subset.xml` for `delay`, `cleanup`, `localwrite`, and `mfilewrite`

- [ ] **Step 2: Verify fixture shape quickly**

Run:
```bash
sed -n '1,220p' testdata/workloads/xml-inheritance-subset.xml
sed -n '1,220p' testdata/workloads/xml-attribute-subset.xml
sed -n '1,220p' testdata/workloads/xml-special-ops-subset.xml
```

Expected:
- fixtures are readable, focused, and each target a single concern set

### Task 2: Add Raw XML Parser Coverage

**Files:**
- Modify: `internal/infrastructure/xml/workload_parser_test.go`

- [ ] **Step 1: Write failing raw-parser tests**

Add tests asserting:
- inheritance fixture preserves storage placement and omitted ratio defaults to `100`
- attribute fixture preserves trigger/closure/interval/division/rampup/rampdown/driver fields
- special-op fixture preserves raw work types and storage placement for file-backed ops

- [ ] **Step 2: Run raw-parser tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml
```

Expected:
- failures because the new fixtures and assertions are not fully in place yet

- [ ] **Step 3: Adjust fixtures or assertions minimally until green**

Only fix fixture/test mismatches. Do not add parser features in this task.

- [ ] **Step 4: Re-run raw-parser tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml
```

Expected:
- package passes with the new fixture assertions

### Task 3: Add Normalized Domain Parser Coverage

**Files:**
- Modify: `internal/workloadxml/parser_test.go`

- [ ] **Step 1: Write failing normalized-parser tests**

Add tests asserting:
- inheritance fixture produces merged config and correct effective storage
- special-op fixture normalizes `delay` and `cleanup` correctly
- zero-ratio operations are absent from normalized output

- [ ] **Step 2: Run normalized-parser tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/workloadxml
```

Expected:
- failures because the new fixtures/assertions are not yet aligned

- [ ] **Step 3: Fix only fixture/test expectation mismatches**

If a failure exposes a real bug in current normalization, fix that bug minimally. Otherwise keep the slice test-only.

- [ ] **Step 4: Re-run normalized-parser tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/workloadxml
```

Expected:
- package passes with the new coverage

### Task 4: Update Compatibility Documentation

**Files:**
- Modify: `docs/xml-compat-matrix.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Reference the new fixtures in the compatibility matrix**

Update the “Real-Workload Tie-In” section so it names the new fixtures and what each one locks in.

- [ ] **Step 2: Update board/checklist state**

Reflect that broader XML fixture coverage has advanced through representative high-value fixture additions.

### Task 5: Final Verification

**Files:**
- Review only: `testdata/workloads/xml-inheritance-subset.xml`
- Review only: `testdata/workloads/xml-attribute-subset.xml`
- Review only: `testdata/workloads/xml-special-ops-subset.xml`
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
- this slice stays fixture/test/doc focused
