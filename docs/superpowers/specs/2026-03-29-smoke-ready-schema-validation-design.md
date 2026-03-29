# Smoke Ready Schema Validation Design

## Goal

Add a repository-local helper and Make targets that validate the current `smoke-ready-json` output against `docs/smoke-ready.schema.json`.

## Problem

The repository now has:

- `smoke-ready-json`
- `docs/smoke-ready.schema.json`
- a contract test that validates mocked helper output against that schema

But there is still no operator-facing command that answers a direct question:

> Does the current helper output conform to the published schema?

Right now the only way to answer that is by running pytest.

## Desired Behavior

Add a dedicated validation helper that:

- executes `scripts/smoke_ready.py --json`
- loads `docs/smoke-ready.schema.json`
- validates the payload with `jsonschema`
- returns exit code `0` on success and non-zero on schema mismatch or helper failure

Expose it through two Make targets:

- `make --no-print-directory smoke-ready-validate`
- `make --no-print-directory smoke-ready-validate-json`

## Scope

In scope:

- one new validation helper script
- one focused test file for that helper
- two new Make targets
- one short README note

Out of scope:

- changing `smoke_ready.py` payload semantics
- changing the schema itself
- changing remote smoke workflows
- adding CI gating around schema validation

## Design

Use a separate helper script instead of extending `scripts/smoke_ready.py`.

That keeps responsibilities clean:

- `scripts/smoke_ready.py` produces readiness payloads
- `scripts/validate_smoke_ready_schema.py` validates those payloads

The validator should support:

- human-readable default output
- `--json` for machine consumers

Recommended JSON output shape:

- `schema_path`
- `schema_version`
- `valid`
- `repo`
- `generated_at`
- `error`

The helper should shell out to `scripts/smoke_ready.py --json` instead of importing internals directly, so it validates the actual CLI contract users consume.

## Acceptance Criteria

- the repository exposes a direct schema validation helper
- Make targets exist for text and JSON modes
- tests prove the helper reports a valid payload against the current schema
- existing `smoke-ready` behavior remains unchanged
