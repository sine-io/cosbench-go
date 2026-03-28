# Workflow Upload Artifact V7 Design

## Goal

Upgrade all repository workflows that use `actions/upload-artifact` to the current official release `v7.0.0` and lock that version with an automated workflow contract test.

## Problem

The repository now has clean remote smoke coverage and no remaining `download-artifact` Node20 warning, but workflow action pins are still inconsistent:

- `actions/download-artifact` has already been upgraded to the current official release
- `actions/upload-artifact` is still pinned to `v6.0.0` in multiple workflows

This is not a current functional failure, but it creates version drift and makes the workflow surface less consistent than the rest of the recent hygiene work.

## Scope

### In Scope

- upgrade every `actions/upload-artifact` reference in `.github/workflows/` to `v7.0.0`
- add one small test that scans workflow files and locks that version

### Out Of Scope

- helper changes
- workflow trigger changes
- new workflows
- changes to `checkout` or `setup-go`

## Recommended Approach

Use one generic test instead of extending each workflow-specific test individually.

### Why a generic test

- multiple workflows use `upload-artifact`
- one scan-based test keeps the policy in a single place
- future workflow additions will fail fast if they reintroduce older pins

## Success Criteria

This slice is complete when:

1. every `uses: actions/upload-artifact@...` in `.github/workflows/` is `v7.0.0`
2. the generic workflow contract test enforces that policy
3. repository tests and build remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
