# Compatibility Closure Design

## Goal

Close the remaining compatibility items that are currently parser-covered, alias-only, or semantically divergent inside the S3/SIO migration boundary.

This slice is the base layer for all later controller and remote work.

## Problems To Solve

The current repository still has meaningful runtime gaps:

- `filewrite` is present in legacy SineIO samples but not implemented in Go
- explicit XML `<auth>` nodes are tolerated but not modeled in runtime resolution
- `siov1` / `gdas` are treated mostly as compatibility aliases instead of explicit runtime profiles
- `prefetch` / `range-read` config survives parsing but not execution
- some documented S3/SIO differences remain open, including config defaults and selected API semantics

## In Scope

- `filewrite` operation semantics
- domain and XML modeling for explicit `<auth>` nodes used by the S3/SIO migration path
- auth inheritance and resolution into execution-ready config
- explicit runtime policy handling for `sio`, `siov1`, and `gdas`
- runtime execution for `is_prefetch`, `is_range_request`, `file_length`, and `chunk_length`
- closure or explicit codification of the currently documented S3/SIO delta list
- regression fixtures and comparison coverage for the new behavior

## Out Of Scope

- non-S3 backend auth schemes
- non-S3 backend drivers
- distributed driver protocol
- UI additions other than surfaces strictly needed for compatibility verification

## Recommended Approach

Implement compatibility closure as a shared config-resolution and execution enhancement, not as isolated one-off patches.

### 1. Model explicit auth as first-class domain data

Introduce `AuthSpec` in the workload domain so XML auth blocks are not discarded. The model should preserve:

- auth type/name
- raw config payload
- location in the inheritance chain

For this roadmap, the execution target is still S3/SIO-family credential resolution, so the model should be generic while the runtime resolver stays pragmatic.

### 2. Add a single config resolution pipeline

Execution should resolve effective runtime config from:

1. endpoint material
2. workload/stage/work storage config
3. inherited auth config
4. operation config

The same resolution logic must be reused by:

- preflight validation
- local execution
- remote execution later

This prevents parser/runtime drift from reappearing in another layer.

### 3. Implement `filewrite` by reusing the file-backed operation path

Legacy `filewrite` is a single-part upload backed by local file input. The Go implementation should:

- parse file selection config using the same file-backed targeting helpers used for `localwrite` / `mfilewrite`
- open the selected file
- upload with `PutObject`
- record uploaded byte count from the actual file size

This should share as much code as possible with existing file-backed operations.

### 4. Convert storage aliases into explicit runtime profiles

Do not leave `siov1` and `gdas` as blind aliases forever.

Instead:

- keep one S3-compatible adapter implementation where practical
- add backend-kind-specific config normalization and behavior branches where legacy semantics differ

Minimum required distinction:

- `sio`: current SIO v2/LTS profile
- `siov1`: legacy SIO v1 compatibility profile
- `gdas`: restore-oriented compatibility profile with its own defaults

The adapter may still share transport code, but the config and behavior policy must become explicit.

### 5. Add read-side execution semantics for prefetch and range requests

When the resolved config enables:

- `is_prefetch=true`
- `is_range_request=true`

the adapter must modify the read request rather than only preserving the keys in config strings.

Behavior target:

- prefetch mode sets the legacy-compatible request header
- range mode issues bounded byte-range requests using `file_length` and `chunk_length`
- reported byte counts reflect bytes actually read

### 6. Close known adapter deltas intentionally

The existing comparison notes already identify several differences. This project should either:

- close the delta in code
- or make it an explicit, documented compatibility policy with tests

Priority deltas:

- S3 default config strictness
- SIO path-style default
- delete tolerance
- list result surface
- slash-containing container handling

## Data Model Changes

Expected additions or changes:

- `AuthSpec` and inheritance helpers in the workload/domain model
- execution config resolver that can merge auth, endpoint, storage, and operation layers
- backend profile enumeration for `s3`, `sio`, `siov1`, `gdas`

## Testing Strategy

### Parser/domain

- fixtures for explicit auth inheritance
- fixtures for `filewrite`
- fixtures for runtime-meaningful range/prefetch shapes

### Execution

- unit tests for `filewrite`
- adapter tests for range/prefetch request shaping
- profile tests for `sio` / `siov1` / `gdas`
- regression tests for resolved config precedence

### Comparison

- extend representative compare fixtures where safe
- record whether a behavior is now `match`, `acceptable delta`, or intentionally different

## Success Criteria

This slice is complete when:

1. `filewrite` runs through the local execution engine
2. explicit `<auth>` nodes are parsed, normalized, and resolved into runtime config
3. `siov1` / `gdas` are no longer alias-only in behavior policy
4. prefetch and range-read config affects actual read execution
5. documented S3/SIO deltas are either closed or converted into explicit tested policy
6. comparison docs and fixtures reflect the new state

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
