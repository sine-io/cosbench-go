# Deferred Parser Coverage Design

## Goal

Add parser-facing coverage for deferred XML constructs that are adjacent to the current S3/SIO migration path, without claiming that the full deferred behaviors are implemented.

This slice should prove:

- these XML shapes can be parsed without breaking the current local-only v1 path
- compatibility aliases and config strings are preserved in the normalized model
- the compatibility matrix stays honest about what is parser-covered versus runtime-supported

## Why This Slice

The repository already has strong coverage for the active XML subset, but the compatibility docs still list several adjacent shapes as deferred or compatibility-only.

The most valuable next step is not to implement those features. It is to lock down that the parser can ingest them cleanly so future work starts from a stable baseline.

## Scope

### In Scope

- parser-facing coverage for `siov1` and `gdas` storage-type aliases
- parser-facing coverage for `sio` prefetch/range-read config keys
- documentation updates to distinguish parser coverage from runtime support

### Out of Scope

- explicit `<auth>` modeling
- actual prefetch/range-read execution behavior
- GDAS backend implementation
- full non-S3 driver compatibility
- legacy plugin ecosystem support

## Recommended Approach

Focus on XML fixtures and parser assertions only.

The fixtures should demonstrate that the current parser and normalization layers can preserve:

- legacy compatibility storage types such as `siov1` and `gdas`
- config keys such as:
  - `is_prefetch`
  - `is_range_request`
  - `file_length`
  - `chunk_length`

The tests should not treat those keys as behaviorally active. They should only prove that the XML survives parsing into the current domain model.

## Fixture Set

### 1. `xml-compat-storage-subset.xml`

Purpose:

- prove parser acceptance of `siov1` and `gdas` storage-type values
- prove normalization still allows `mprepare`, `mwrite`, and `restore` shapes to parse when paired with those compatibility aliases

This fixture is about compatibility alias ingestion, not about runtime execution against real GDAS.

### 2. `xml-range-prefetch-subset.xml`

Purpose:

- prove parser preservation of prefetch and range-read related config keys in the XML path
- keep those keys visible in the normalized workload model for future implementation work

Target config keys:

- `is_prefetch`
- `is_range_request`
- `file_length`
- `chunk_length`

## Test Placement

### Raw parser coverage

Use `internal/infrastructure/xml/workload_parser_test.go` to assert:

- storage types are preserved as `siov1` / `gdas`
- raw work and operation config strings contain the expected prefetch/range keys

### Normalized parser coverage

Use `internal/workloadxml/parser_test.go` to assert:

- compatibility storage aliases remain visible after normalization
- prefetch/range config keys remain present in the normalized config strings

No control-plane or execution tests are needed in this slice.

## Documentation Update

Update `docs/xml-compat-matrix.md` so it explicitly distinguishes:

- parser-covered deferred constructs
- runtime-supported constructs

The new fixtures should be referenced as evidence that:

- the XML can be ingested
- the behavior remains deferred

## Success Criteria

This slice is complete when:

1. representative fixtures exist for compatibility storage aliases and range/prefetch config shapes
2. raw parser tests cover those fixtures
3. normalized parser tests cover those fixtures
4. the compatibility matrix references the new fixtures and preserves the deferred/runtime distinction
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
