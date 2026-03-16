# TODO (Phase 0)

## Environment
- [x] Go toolchain available via `/snap/bin/go` (note: `go` may not be on PATH in this runtime)

## Bootstrap
- [ ] Initialize go.mod: module `github.com/sine-io/cosbench-go`
- [ ] Add base dependencies: cobra, viper, gin, zerolog
- [ ] Add lint/test baseline (go test ./...)

## Domain & Parsing (TDD)
- [ ] Unit test: parse `s3-config-sample.xml` workload basics
- [ ] Unit test: parse `sio-config-sample.xml` workload basics + SIO params
- [ ] Implement XML -> domain mapping
- [ ] Parse storage config string `k=v;k2=v2` with escaping rules (define)

## Execution engine
- [ ] Define operation scheduler and ratio selection
- [ ] Metrics model and aggregation

## Controller/Driver protocol
- [ ] Decide transport: HTTP/JSON first (gin)
- [ ] Define endpoints and DTOs
