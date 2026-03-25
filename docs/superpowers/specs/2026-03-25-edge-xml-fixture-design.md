# Edge XML Fixture Design

## Goal

Expand representative XML coverage for edge-case workload shapes that still fit the current local-only v1 boundary.

This slice is about capturing realistic structure and stage behavior from legacy sample families, not adding deferred features or unsupported backends.

## Why This Slice

The repository already covers:

- active S3/SIO subset fixtures
- inheritance, attribute-heavy XML, and special-op shapes
- stage-aware mock continuity

What is still under-covered are several practical workload structures that appear in legacy sample material and matter for local behavior:

- explicit delay stages with `closuredelay`
- split read/write container targeting inside one work
- reuse-data workflows with multiple main stages between prepare and cleanup

These are not new features. They are higher-order compositions of behavior we already claim to support.

## Scope

### In Scope

- new representative fixtures under `testdata/workloads/`
- parser-level coverage for those shapes
- minimal control-plane behavior coverage where stage sequencing matters
- compatibility-matrix updates for the new fixtures

### Out of Scope

- bringing legacy `swift` or `ampli` backends into current scope
- explicit `<auth>` support
- prefetch/range-read XML semantics
- non-local controller/driver semantics

## Recommended Approach

Adapt legacy sample *shapes* into current-scope fixtures using the supported `mock` backend.

This keeps the focus on:

- stage ordering
- workload structure
- config targeting patterns
- local execution continuity

without dragging unsupported storage ecosystems into the current phase.

## Fixture Set

### 1. `xml-delay-stage-subset.xml`

Derived from the legacy `delay-stage` sample family.

Purpose:

- lock repeated `delay` stage structure
- lock stage-level `closuredelay`
- confirm delay stages coexist correctly with surrounding `init` / `prepare` / `main` / `cleanup`

### 2. `xml-splitrw-subset.xml`

Derived from the legacy `splitrw` sample family.

Purpose:

- lock a main work where `read` and `write` target different container/object ranges
- confirm ratio-bearing mixed operations can still express split target sets cleanly

### 3. `mock-reusedata-subset.xml`

Derived from the legacy `reusedata` sample family.

Purpose:

- lock a multi-main-stage workflow that reuses prepared data
- confirm local mock continuity supports `prepare -> main_1 -> main_2 -> cleanup -> dispose`

This fixture should be the primary behavior-level proof that the new mock stage-aware continuity is useful.

## Test Placement

### Raw parser coverage

Use `internal/infrastructure/xml/workload_parser_test.go` to assert:

- stage counts
- preserved stage/work attributes
- split read/write configs remain distinct
- delay-stage structure is parsed as expected

### Normalized parser coverage

Use `internal/workloadxml/parser_test.go` to assert:

- `delay` stage/work normalization still produces the expected synthetic operation
- split read/write configs survive normalization
- multi-main-stage reuse fixture keeps the expected stage layout

### Control-plane coverage

Use `internal/controlplane` tests for:

- `mock-reusedata-subset.xml`
- optionally `xml-delay-stage-subset.xml`

The goal is not full benchmark verification. The goal is to prove that the shape runs coherently under the local control plane.

## Success Criteria

This slice is complete when:

1. the three new fixtures exist under `testdata/workloads/`
2. parser tests cover their intended XML shape
3. control-plane tests cover at least the reuse-data continuity path
4. `docs/xml-compat-matrix.md` references the new fixtures
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
