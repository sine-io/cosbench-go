# CI Automation Design

## Goal

Add a minimal, repository-local CI workflow that runs the existing build verification path automatically on code changes.

The workflow should:

- use the repository’s current `Makefile` automation instead of inventing a second verification path
- run on normal collaboration events
- avoid requiring live endpoint credentials
- stay small enough that contributors can reason about it quickly

## Why This Slice

The repository already has:

- `make build`
- `make test`
- `make vet`
- `make validate`
- opt-in `make smoke-s3`

What it still lacks is an external automation definition that actually runs these checks when code changes land.

## Scope

### In Scope

- one GitHub Actions workflow under `.github/workflows/`
- automatic execution of the existing validation path
- contributor-facing documentation that points at the workflow

### Out of Scope

- live smoke tests in default CI
- release automation
- matrix testing across many operating systems
- artifact publishing
- external secret management beyond what GitHub Actions already provides

## Recommended Approach

Use a single GitHub Actions workflow that:

- triggers on `push` and `pull_request`
- checks out the repository
- installs the Go version from `go.mod`
- runs `make validate` with `GO=go`

This is the smallest useful automation layer because it reuses the exact local verification path the repository already documents.

## Why GitHub Actions

There is currently no existing CI provider configuration in the repository.

GitHub Actions is the most natural default here because:

- the repository already uses git-based collaboration
- workflow config can live inside the repo
- no extra wrapper scripts are needed just to prove the path works

## Workflow Shape

### Trigger

- `push`
- `pull_request`

### Steps

1. checkout source
2. set up Go from `go.mod`
3. run `make validate` with `GO=go`

### Explicit non-goals

Do not run `make smoke-s3` in the default workflow. That path is intentionally opt-in and credential-dependent.

## Documentation Update

Update repository-facing docs to mention:

- that a CI workflow now exists
- that CI runs `make validate`
- that live smoke tests remain opt-in and separate

Good targets:

- `README.md`
- `AGENTS.md`
- `BOARD.md`
- `TODO.md`

## Success Criteria

This slice is complete when:

1. `.github/workflows/ci.yml` exists
2. the workflow runs `make validate`
3. smoke tests remain outside default CI
4. repository docs mention the CI path clearly
5. `go test ./...` and `go build ./...` still pass locally

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
