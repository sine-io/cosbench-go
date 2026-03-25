# Auth Parser-Tolerated Coverage Design

## Goal

Add parser-facing coverage proving that legacy XML files containing `<auth>` elements can still be ingested by the current parser and normalization path, without claiming that authentication is modeled or executed in the local-only v1 scope.

## Why This Slice

The current compatibility docs explicitly say:

- explicit `<auth>` element modeling is deferred

That does not necessarily mean the parser should fail when `<auth>` appears in otherwise supported XML.

For migration work, it is useful to prove the narrower claim:

- auth-bearing XML shapes are parser-tolerated
- auth nodes are ignored rather than modeled
- supported workload/storage/work/operation structure still survives normalization

## Scope

### In Scope

- two representative fixtures containing `<auth>` nodes
- raw parser tests showing those XML files still parse
- normalized parser tests showing supported surrounding structure remains intact
- compatibility-matrix wording that distinguishes parser tolerance from auth modeling

### Out of Scope

- introducing an `AuthSpec` into the domain model
- auth inheritance rules
- runtime authentication behavior
- non-S3/SIO backend support

## Recommended Approach

Treat `<auth>` as **parser-tolerated but not modeled**.

This is the narrowest useful claim that fits the current scope boundary:

- XML containing `<auth>` should not necessarily break the parser
- the auth nodes themselves should not appear in the current normalized domain model

## Fixture Set

### 1. `xml-auth-tolerated-subset.xml`

Purpose:

- prove that workload-, stage-, and work-level `<auth>` nodes can coexist with supported `mock` storage and current work/operation shapes
- prove the parser still extracts workload, stage, work, storage, and operation data normally

### 2. `xml-auth-none-subset.xml`

Purpose:

- prove that `type="none"` auth-bearing XML shapes from the legacy sample family are also parser-tolerated
- keep the fixture minimal and current-scope by pairing it with supported structure rather than unsupported storage backends

## Test Placement

### Raw parser coverage

Use `internal/infrastructure/xml/workload_parser_test.go` to assert:

- parsing succeeds
- workload/stage/work/operation structure is preserved
- supported surrounding fields such as storage/config/work type still parse correctly

### Normalized parser coverage

Use `internal/workloadxml/parser_test.go` to assert:

- normalization still produces the expected domain structure
- no auth-specific fields appear in the normalized model
- auth presence does not disturb supported storage/config/operation parsing

## Documentation Update

Update `docs/xml-compat-matrix.md` so the deferred auth row distinguishes:

- parser-covered / parser-tolerated auth-bearing XML
- explicit auth modeling still deferred

This avoids the misleading interpretation that any `<auth>` node is currently a hard parse failure.

## Success Criteria

This slice is complete when:

1. representative auth-bearing fixtures exist under `testdata/workloads/`
2. raw parser tests prove those fixtures parse successfully
3. normalized parser tests prove surrounding supported structure remains intact
4. the compatibility matrix clearly states “parser-tolerated, not modeled”
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
