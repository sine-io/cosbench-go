# Storage Config Fallback Design

## Goal

Close the current legacy-parity gap where `part_size` and `restore_days` are parsed from storage config but not treated as execution-time defaults by the Go implementation.

The fix should preserve the existing ability for operation config to override those values explicitly.

## Problem

Current `cosbench-go` behavior:

- `internal/driver/s3/config.go` parses `part_size` and `restore_days` from storage config
- `internal/domain/execution/opconfig.go` separately parses `part_size` and `restore_days` from operation config
- `internal/domain/execution/engine.go` uses only the operation-config parse result during `mwrite`, `mfilewrite`, and `restore`

This creates a source-layer mismatch:

- storage-level `part_size` is not the effective default for multipart execution
- storage-level `restore_days` is not the effective default for restore execution

Legacy Java behavior uses storage-level configuration as the initialized source of truth for these values.

## Scope

### In Scope

- execution-time fallback from storage config to operation config for:
  - `part_size`
  - `restore_days`
- aligned preflight validation using the same merged config view
- focused regression tests for storage-default and op-override behavior

### Out of Scope

- changing unrelated config inheritance rules
- broad refactors of workload normalization
- adding new storage config keys
- live endpoint verification

## Recommended Approach

Implement a merged execution-config path at the execution layer:

- base config comes from `work.Storage.Config`
- operation config overlays the base
- explicit operation values win

This keeps responsibility in the execution layer, where these fields are actually consumed, and avoids smearing execution semantics into XML normalization.

## Behavioral Target

### Multipart operations

For `mwrite` and `mfilewrite`:

- if operation config sets `part_size`, use it
- else if storage config sets `part_size`, use it
- else default to `5 MiB`

### Restore operations

For `restore`:

- if operation config sets `restore_days`, use it
- else if storage config sets `restore_days`, use it
- else default to `1`

## Implementation Shape

The smallest coherent implementation is:

1. add a config merge helper in the execution-config path
2. use that helper from runtime execution
3. use the same helper from preflight validation

The merge rule should be simple:

- parse storage config KV first
- parse operation config KV second
- operation entries override same-named storage entries

## Tests To Add

### Execution tests

Add tests proving:

- storage-level `part_size` drives `MultipartPut` when op config omits it
- operation-level `part_size` overrides storage-level `part_size`
- storage-level `restore_days` drives `RestoreObject` when op config omits it
- operation-level `restore_days` overrides storage-level `restore_days`

### Preflight tests

If preflight depends on the same merged config path, verify it stays aligned with runtime behavior instead of using a different interpretation.

## Success Criteria

This slice is complete when:

1. multipart execution uses storage-level `part_size` by default
2. restore execution uses storage-level `restore_days` by default
3. explicit operation config still overrides those defaults
4. regression tests cover both fallback and override behavior
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
