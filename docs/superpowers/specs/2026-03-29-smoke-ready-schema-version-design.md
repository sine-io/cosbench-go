# Smoke Ready Schema Version Design

## Goal

Add a lightweight `schema_version` to `smoke-ready` / `smoke-ready-json` and lock the current summary schema with a dedicated contract test.

## Problem

`smoke-ready` has evolved into a real machine-consumable interface with many fields:

- result
- source
- url
- artifact
- created_at

But it still has no explicit version marker, and all contract coverage is folded into one large behavior test.

That makes it harder for downstream consumers to:

- detect schema changes intentionally
- distinguish a breaking schema edit from a bug

## Desired Behavior

Add:

- top-level `schema_version`

And add one dedicated contract test that asserts the presence of the current summary schema keys without re-testing every behavioral detail.

## Scope

In scope:

- `scripts/smoke_ready.py`
- `scripts/test_smoke_ready.py`
- one new focused schema contract test
- small README wording update

Out of scope:

- changing existing summary semantics
- renaming existing fields
- adding a formal JSON Schema document

## Design

Use a simple stable version:

- `schema_version: 1`

Keep it top-level so consumers can read it before parsing nested sections.

Add `scripts/test_smoke_ready_schema.py` to assert:

- `schema_version == 1`
- core top-level objects exist
- expected summary keys exist

This test should be narrow and structural, not a duplicate of the current behavior-heavy smoke-ready tests.

## Acceptance Criteria

- JSON output includes `schema_version`
- contract test passes
- existing smoke-ready behavior remains unchanged
