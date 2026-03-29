# Smoke Ready JSON Schema Design

## Goal

Add a repository-owned JSON Schema document for `smoke-ready-json` and validate the helper output against it.

## Problem

`smoke-ready-json` now has:

- a top-level `schema_version`
- a stable machine-consumable summary shape
- a structural contract test

But the contract still only exists in Python assertions. Downstream tooling has no canonical schema file to reference, and contributors cannot inspect the interface without reading tests.

## Desired Behavior

Add one committed JSON Schema file that describes the current `smoke-ready-json` shape and make the existing schema contract test validate real helper output against that file.

## Scope

In scope:

- one JSON Schema file for the current `schema_version: 1`
- one focused contract test update to validate against the schema
- one short README note pointing machine consumers at the schema file

Out of scope:

- changing helper semantics
- renaming existing fields
- adding new smoke-ready fields
- publishing the schema outside the repository

## Design

Add `docs/smoke-ready.schema.json` with:

- draft 2020-12 JSON Schema metadata
- top-level required keys already emitted by the helper
- explicit object shapes for `repo_secrets`, `workflows`, and `summary`
- required summary fields for the current real-endpoint, legacy-live, and remote result blocks

Keep the schema practical rather than exhaustive:

- require stable keys and broad value types
- allow `null` for workflow latest entries that may be absent
- avoid over-constraining dynamic arrays like `required` and `blockers`

Update `scripts/test_smoke_ready_schema.py` to:

- load `docs/smoke-ready.schema.json`
- run the helper with the existing mocked workflow inputs
- validate the JSON payload with `jsonschema`
- keep one explicit assertion for `schema_version == 1`

## Acceptance Criteria

- `docs/smoke-ready.schema.json` exists and describes the current interface
- `scripts/test_smoke_ready_schema.py` validates helper output against it
- `README.md` points machine consumers at the schema file
- existing smoke-ready semantics stay unchanged
