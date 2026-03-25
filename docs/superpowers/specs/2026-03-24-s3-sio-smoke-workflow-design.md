# S3/SIO Smoke Test Workflow Design

## Goal

Add a real-endpoint smoke-test workflow for the S3/SIO driver without polluting the normal local test path.

The workflow should:

- exercise the existing `internal/driver/s3` adapter against a live endpoint
- stay opt-in through environment variables
- default to `t.Skip` when no live endpoint is configured
- provide a simple `make` entrypoint for operators and CI jobs that explicitly want smoke coverage

## Why This Slice

The repository already has:

- mock-backed integration coverage
- parser, execution, control-plane, and request-level adapter tests

What it still lacks is a minimal live-endpoint confirmation path. The gap is not in default correctness under local tests, but in verifying that the real S3/SIO wiring behaves as expected once requests leave the process.

## Scope

### In Scope

- a smoke test file under `internal/driver/s3`
- environment-variable driven configuration for live endpoint runs
- one happy-path end-to-end object lifecycle
- conditional multipart coverage for `sio`
- a thin `make smoke-s3` entrypoint
- short documentation for setup and execution

### Out of Scope

- always-on integration tests in `go test ./...`
- fixture-driven CLI smoke runs
- large benchmark workloads against live systems
- retry/backoff tuning work
- remote-worker behavior
- broad performance assertions

## Recommended Approach

Implement the smoke workflow as package-level Go tests in `internal/driver/s3`.

This is the smallest change that:

- reuses the real adapter and config parsing
- produces actionable test failures close to the failing API call
- avoids new binaries or shell orchestration layers

## Test Activation Model

The smoke tests must be opt-in.

Activation rule:

- if required environment variables are missing, the test calls `t.Skip(...)`
- if present, the test runs against the configured endpoint

Required environment variables:

- `COSBENCH_SMOKE_ENDPOINT`
- `COSBENCH_SMOKE_ACCESS_KEY`
- `COSBENCH_SMOKE_SECRET_KEY`

Optional environment variables:

- `COSBENCH_SMOKE_BACKEND`
  - default: `s3`
  - accepted: `s3`, `sio`
- `COSBENCH_SMOKE_REGION`
  - default: `us-east-1`
- `COSBENCH_SMOKE_PATH_STYLE`
  - default: backend-specific adapter default
- `COSBENCH_SMOKE_BUCKET_PREFIX`
  - default should be a stable `cosbench-go-smoke`

## Test Shape

The smoke test should generate a unique bucket/object suffix per run so repeated executions do not collide.

### Common happy path

For both `s3` and `sio` backends:

1. build the adapter from env-backed config
2. create a unique bucket
3. upload a small object with `PutObject`
4. read it back with `GetObject`
5. inspect metadata with `HeadObject`
6. confirm discoverability with `ListObjects`
7. delete the object
8. delete the bucket

Assertions should be intentionally minimal and stable:

- payload round-trip matches
- `HeadObject` content length is correct
- `ListObjects` includes the uploaded key

### SIO-specific multipart path

If `COSBENCH_SMOKE_BACKEND=sio`, run one additional multipart upload:

- upload a payload larger than one part through `MultipartPut`
- verify it is visible through `HeadObject`
- delete it during cleanup

This keeps SIO-specific behavior in scope without forcing multipart assertions on normal S3 runs.

## Make Target

Add a thin `make smoke-s3` target that runs:

```bash
$(GO) test ./internal/driver/s3 -run Smoke -v
```

This target is deliberately narrow. It should not wrap the full suite.

## Documentation

Update repository-facing documentation so contributors know:

- smoke tests are opt-in
- which environment variables are required
- how to run `make smoke-s3`
- that failures here indicate live-endpoint or credentials/config issues, not necessarily local unit regressions

Good documentation targets are:

- `README.md`
- `AGENTS.md`
- `BOARD.md`
- `TODO.md`

## Failure Handling

Keep cleanup best-effort:

- attempt object deletion and bucket deletion even if an earlier assertion fails
- use `t.Cleanup` where possible

The test should fail on real operation errors, but it should not leave easily avoidable residue when cleanup can still run.

## Verification

Implementation is complete for this slice when all of the following are true:

1. `go test ./internal/driver/s3 -run Smoke -v` skips cleanly with no env configured
2. `make smoke-s3` invokes the same smoke path
3. normal `go test ./...` remains green without requiring live credentials
4. documentation reflects the opt-in smoke workflow

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
