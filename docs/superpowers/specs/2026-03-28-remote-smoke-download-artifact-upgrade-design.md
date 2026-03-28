# Remote Smoke Download Artifact Upgrade Design

## Goal

Upgrade the `actions/download-artifact` usage in `Remote Smoke Matrix` from the current Node20-targeting release to the current official release that removes the remaining deprecation warning.

## Problem

The matrix workflow already forces JavaScript actions onto Node24, so the workflow is functionally green. The remaining issue is a warning from GitHub Actions noting that `actions/download-artifact@v6.0.0` targets Node20 and is only being coerced onto Node24.

This is not a functional failure, but it is noise in the workflow and it weakens the “clean signal” of remote smoke verification.

## Scope

### In Scope

- upgrade `actions/download-artifact` in `.github/workflows/remote-smoke-matrix.yml`
- lock the version in the existing workflow contract test

### Out Of Scope

- helper changes
- matrix row changes
- upload-artifact version changes
- default CI changes

## Recommended Approach

Use the current official release tag from `actions/download-artifact`:

- `v8.0.1`

This is the smallest change that directly addresses the warning surfaced by the latest `Remote Smoke Matrix` run.

## Success Criteria

This slice is complete when:

1. the matrix workflow uses `actions/download-artifact@v8.0.1`
2. the workflow contract test locks that exact version
3. repository tests and build remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
