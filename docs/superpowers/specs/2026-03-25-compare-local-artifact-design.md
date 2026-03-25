# Compare Local Artifact Design

## Goal

Improve the manual `compare-local` GitHub Actions workflow by preserving its output as an artifact.

This should make remote execution of the local comparison set more useful without changing what the workflow actually runs.

## Scope

### In Scope

- capture `make compare-local` output into a file
- upload that file as a workflow artifact
- document that the manual workflow produces downloadable output

### Out of Scope

- changing the `compare-local` command itself
- automatic matrix updates
- live endpoint automation

## Recommended Approach

Keep the workflow command the same in substance:

- run `make compare-local`

but tee the output into a file and upload that file using `actions/upload-artifact`.

## Success Criteria

1. the manual compare-local workflow uploads an artifact
2. the artifact contains the command output
3. docs mention the artifact availability
4. `go test ./...` and `go build ./...` remain green
