# Compare Local Workflow Design

## Goal

Add a manual GitHub Actions workflow that runs the existing `make compare-local` target on demand.

This should make the repository's safe local comparison set runnable in remote automation without requiring live endpoint credentials.

## Scope

### In Scope

- one `workflow_dispatch` GitHub Actions workflow
- reuse of the existing `make compare-local` target
- documentation updates pointing contributors at the new manual workflow

### Out of Scope

- automatic live endpoint comparison
- running `make smoke-s3` in default CI
- scheduled workflows
- release publishing

## Recommended Approach

Add a second workflow under `.github/workflows/` that:

- is triggered only by `workflow_dispatch`
- checks out the repository
- installs Go from `go.mod`
- runs `GO=go make compare-local`

This keeps the existing default CI small while still making the local comparison set automation-friendly.

## Success Criteria

1. a manual compare-local workflow exists
2. it runs `make compare-local`
3. docs distinguish it from the default CI and from live smoke workflows
4. `go test ./...` and `go build ./...` remain green
