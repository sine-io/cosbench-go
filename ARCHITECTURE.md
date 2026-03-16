# Architecture (v0)

This project follows:
- **Clean Architecture**: dependency rule (inner layers know nothing about outer layers)
- **DDD Lite**: explicit domain model + boundaries; avoid over-engineering aggregates early
- **CQRS**: separate command (write) and query (read) use-cases when it clarifies; not mandatory everywhere
- **DIP**: infrastructure depends on interfaces defined in application/domain

## Proposed module layout

```
cmd/
  cosbench-go/            # root cobra app
internal/
  domain/
    workload/             # workload model (workflow, workstage, work, operation)
    storage/              # storage concepts: S3/SIO config, capabilities
    metrics/              # measurements, counters, histograms, summaries
  application/
    command/
    query/
    ports/                # interfaces (StorageClient, Clock, IDGen, MetricsSink, ...)
  interfaces/
    http/                 # gin handlers, request/response DTOs
    cli/                  # cobra commands wiring
  infrastructure/
    config/               # viper loading, config structs
    logger/               # zerolog setup
    storage/
      s3/
      sio/
    transport/
      http/
        controller/
        driver/
```

## First milestones (engineering)

1) **Parse XML** into `domain/workload` with unit tests.
2) Define `application/ports.StorageClient` and implement S3/SIO adapters.
3) Implement `driver` execution engine (workers, runtime/totalOps, ratios) with deterministic tests.
4) Implement `controller` orchestration API and driver heartbeat/result collection.

## Compatibility focus

- `<storage type="s3|sio" config="k=v;...">` semicolon-separated kv string
- Workstage/work/work/operation nodes as per legacy samples
- SIO operations/extensions: mprepare, mwrite, mfilewrite, localwrite, restore
