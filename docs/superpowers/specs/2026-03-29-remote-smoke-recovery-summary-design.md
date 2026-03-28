# Remote Smoke Recovery Summary Design

## Goal

Add a stable, downloadable summary artifact to the `Remote Smoke Recovery` workflow so the recovery result can be consumed outside the GitHub job summary UI.

## Problem

The repository already has:

- a dedicated `Remote Smoke Recovery` workflow
- a job summary generated from `.artifacts/remote-smoke/summary.md`
- a broad raw artifact upload of `.artifacts/remote-smoke/`

What is still missing is a compact, stable artifact that contains only the recovery summary payloads:

- `summary.json`
- `summary.md`

Without that, downstream use has to download the full raw smoke artifact and know where the summary files live.

## Scope

### In Scope

- add one small builder script that creates `.artifacts/remote-smoke-recovery-summary/`
- upload that directory as a second artifact from `Remote Smoke Recovery`
- add a lightweight contract test for the new workflow behavior

### Out Of Scope

- helper changes
- recovery scenario logic changes
- matrix workflow changes
- default CI changes

## Recommended Approach

Create a small script that copies the existing summary files into a stable output directory and upload that directory as a second artifact.

### Why a builder script

- it keeps the workflow YAML simple
- it mirrors the pattern already used by the matrix aggregate path
- it makes the output path explicit and testable

## Artifact Contract

The new artifact should be named:

- `remote-smoke-recovery-summary`

Its contents should include:

- `summary.json`
- `summary.md`

The existing raw artifact:

- `remote-smoke-recovery-output`

should remain unchanged.

## Success Criteria

This slice is complete when:

1. `Remote Smoke Recovery` still uploads the raw smoke artifact
2. it also uploads a dedicated `remote-smoke-recovery-summary` artifact
3. the summary artifact contains `summary.json` and `summary.md`
4. repository tests and build remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
