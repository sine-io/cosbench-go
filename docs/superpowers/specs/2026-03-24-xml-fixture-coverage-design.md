# XML Fixture Coverage Design

## Goal

Expand repository fixture coverage for the highest-value XML shapes already claimed by the local-only v1 migration boundary.

This slice is about **coverage**, not new parser features. It should lock in currently supported XML behavior through representative fixtures and parser-level assertions.

## Why This Slice

The repository already has:

- legacy S3 and SIO sample XML fixtures
- representative active-subset fixtures for S3 and SIO multipart paths
- parser tests that prove the happy path for a few sample shapes

What it still lacks is a focused fixture set for several XML shapes that the compatibility docs already describe as supported:

- inheritance across workload / workflow / stage / work / operation config
- stage-level and work-level storage overrides
- attribute-heavy work definitions such as `trigger`, `closuredelay`, `interval`, `division`, `rampup`, `rampdown`, and `driver`
- special-op XML shapes such as `delay`, `cleanup`, `localwrite`, and `mfilewrite`
- zero-ratio filtering and omitted-ratio defaulting

## Scope

### In Scope

- add a small number of new XML fixtures under `testdata/workloads/`
- add parser tests in `internal/infrastructure/xml`
- add normalized-domain tests in `internal/workloadxml`
- update compatibility documentation to mention the new representative fixtures

### Out of Scope

- new XML schema support not already implemented
- control-plane execution tests for every new fixture
- remote-worker XML/runtime semantics
- explicit `<auth>` modeling
- broad fixture generation tooling

## Recommended Approach

Add three narrowly scoped fixtures instead of one giant catch-all file.

This keeps each fixture readable, gives failures a clear cause, and avoids mixing unrelated parser concerns in the same sample.

## Fixture Set

### 1. `xml-inheritance-subset.xml`

Purpose:

- prove config inheritance from workload → workflow → stage → work → operation
- prove stage-level storage override and work-level storage override
- prove omitted operation ratio defaults to `100`

Expected assertions:

- merged stage config is visible after normalization
- merged work config is visible after normalization
- merged operation config is visible after normalization
- stage-local storage overrides workload default
- work-local storage overrides stage default

### 2. `xml-attribute-subset.xml`

Purpose:

- prove XML attributes that are parsed but only lightly covered today

Target attributes:

- workload `trigger`
- stage `closuredelay`
- stage `trigger`
- work `interval`
- work `division`
- work `rampup`
- work `rampdown`
- work `driver`

Expected assertions:

- raw parser preserves all these fields
- normalized domain model still exposes the expected values

### 3. `xml-special-ops-subset.xml`

Purpose:

- prove special work normalization for XML shapes that currently have little or no fixture coverage

Target work types:

- `delay`
- `cleanup`
- `localwrite`
- `mfilewrite`

Expected assertions:

- `delay` normalizes to one synthetic operation and one worker
- `cleanup` injects `deleteContainer=false` when absent
- `localwrite` and `mfilewrite` survive parsing with SIO-compatible storage and valid operation lists

## Test Placement

### Raw parser coverage

Use `internal/infrastructure/xml/workload_parser_test.go` to assert:

- fixture parses successfully
- raw attribute fields and raw storage placement are preserved
- omitted ratio defaults to `100`

### Normalized domain coverage

Use `internal/workloadxml/parser_test.go` to assert:

- normalization output matches expected inherited config
- special works become the expected synthetic operations
- zero-ratio operations are absent from the normalized result

This split mirrors the repository’s existing distinction between raw XML mapping and normalized domain parsing.

## Documentation Update

Update `docs/xml-compat-matrix.md` so the “Real-Workload Tie-In” section points at the new fixtures and explains what each one covers.

The matrix should remain honest: it should describe these fixtures as representative coverage for already-supported shapes, not as evidence of brand-new XML support.

## Success Criteria

This slice is complete when:

1. the three new fixtures exist under `testdata/workloads/`
2. raw parser tests cover their intended attributes and structures
3. normalized parser tests cover inheritance, special work normalization, and zero-ratio filtering
4. `go test ./...` stays green
5. `docs/xml-compat-matrix.md` references the new fixtures

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
