# Repository Guidelines

## Project Structure & Module Organization
`cmd/cosbench-go` contains the CLI runner for XML workloads, and `cmd/server` starts the HTTP control plane. Most application code lives under `internal/`: `domain` holds workload and execution models, `controlplane` manages jobs and endpoints, `executor` runs work, `driver/s3` and `infrastructure/storage` handle backends, `snapshot` persists runtime state, and `web` serves templates. Reference material lives in `docs/`. Sample XML inputs and compatibility fixtures live in `testdata/`. HTML templates and static assets live in `web/templates` and `web/static`.

## Build, Test, and Development Commands
Use the `Makefile` for the common paths:

- `make fmt` formats all Go packages with `go fmt ./...`
- `make test` runs the full test suite with `go test ./...`
- `make build` compiles all packages and binaries with `go build ./...`
- `make compare-local` runs the curated mock-backed comparison fixture set through the CLI
- `make validate` runs `go vet`, tests, and a full build for CI-style verification
- `make smoke-s3` runs the opt-in live endpoint smoke test for `internal/driver/s3`
- `make tidy` syncs `go.mod` and `go.sum`

For local runs, use the package entrypoints directly:

- `go run ./cmd/server -data-dir ./data -view-dir ./web/templates`
- `go run ./cmd/cosbench-go -workload testdata/workloads/s3-active-subset.xml -backend mock -json`
- `go run ./cmd/cosbench-go -f testdata/workloads/s3-active-subset.xml -backend mock`
- `go run ./cmd/cosbench-go testdata/workloads/s3-active-subset.xml -backend mock`
- `go build ./...` to catch compile errors across both binaries

If `/snap/bin/go` is not your Go binary, override the Makefile variable, for example `GO=$(which go) make test`.
Live smoke tests require `COSBENCH_SMOKE_ENDPOINT`, `COSBENCH_SMOKE_ACCESS_KEY`, and `COSBENCH_SMOKE_SECRET_KEY`; without them the smoke suite skips.
Repository CI runs `make validate`; keep smoke tests opt-in and out of the default CI path.
In `-json` mode, stdout is reserved for machine-readable JSON.
`make compare-local` is the fastest way to refresh local comparison evidence without live credentials.

## Coding Style & Naming Conventions
This is a Go repository; follow `gofmt` output exactly and keep package names lowercase. Exported types and functions use `CamelCase`; unexported helpers use `camelCase`. Keep packages focused on one layer or boundary, and prefer small adapters over cross-layer shortcuts. When wrapping errors, preserve context with `%w`, as in `fmt.Errorf("snapshot store: %w", err)`.

## Testing Guidelines
Write table-driven tests where inputs vary, and keep tests adjacent to the code as `*_test.go`. Test names should follow `TestXxxBehavior`, matching existing examples like `TestEngineRunTotalOps` and `TestDashboardRenders`. Use `t.TempDir()` for snapshot-backed tests and reuse XML fixtures from `testdata/workloads` or `testdata/legacy` before adding new inline samples. Run `go test ./...` before opening a PR; add targeted coverage for parser, scheduler, and reporting changes.

## Commit & Pull Request Guidelines
Current history uses short conventional prefixes such as `chore:` and `bootstrap:`. Keep that pattern with imperative summaries, for example `feat: add endpoint export CSV`. PRs should describe the behavior change, list affected packages, and include the verification commands you ran. Include screenshots when changing `web/templates` or `web/static`, and link any migration or compatibility notes when XML behavior changes.
