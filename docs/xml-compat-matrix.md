# XML Compatibility Matrix

This matrix reflects the current local-only v1 behavior of `internal/workloadxml` and the normalization/execution path behind it.

## Supported Now

| XML element / attribute | Status | Notes |
| --- | --- | --- |
| `<workload name description trigger config>` | supported | mapped into `internal/domain.Workload` |
| workload-level `<storage type config>` | supported | inherited by stages and works through normalization |
| `<workflow config>` | supported | merged into stage/work/op config |
| `<workstage name closuredelay trigger config>` | supported | normalized into the domain stage model |
| stage-level `<storage>` | supported | inherited by contained works |
| `<work type workers runtime totalOps totalBytes ...>` | supported subset | core scheduling and limits are honored |
| `<operation type ratio division config id>` | supported | ratio and config inheritance enforced |
| COSBench-style config strings such as `containers=...;objects=...;sizes=...` | supported | parsed by the execution config layer |
| practical storage types `mock`, `s3`, `sio` | supported | `mock` is for local tests, `s3` / `sio` are the migration targets |
| local ops `init`, `prepare`, `write`, `read`, `list`, `cleanup`, `dispose`, `delete`, `head`, `restore`, `mprepare`, `mwrite`, `mfilewrite`, `localwrite`, `delay` | supported | executed through the local engine |

## Supported With Current Constraints

| XML feature | Status | Constraint |
| --- | --- | --- |
| S3 workloads from current sample set | supported | validated by representative sample XML |
| SineIO multipart workloads | supported | depends on S3-compatible endpoint and multipart path |
| stage execution ordering | supported | stages run serially in one process |
| worker concurrency | supported | local goroutines only |
| execution preflight | supported subset | catches obvious config, adapter, and file-input failures before run start |

## Deferred

| XML feature | Status | Reason |
| --- | --- | --- |
| explicit `<auth>` element modeling | deferred runtime, parser-tolerated | auth-bearing XML can be ingested without breaking current parsing, but auth nodes are not modeled in the current domain path |
| `siov1` / `gdas` compatibility storage aliases | deferred runtime, parser-covered | XML shapes are preserved for future compatibility work, but runtime support remains outside the current path |
| `sio` prefetch / range-read config shapes (`is_prefetch`, `is_range_request`, `file_length`, `chunk_length`) | deferred runtime, parser-covered | config survives parsing, but execution semantics are not implemented in the local v1 path |
| full COSBench plugin ecosystem | deferred | outside the active migration subset |
| remote worker-specific XML/runtime behaviors | deferred | local execution first |
| every historical driver/storage variant | deferred | S3/SIO only for this migration path |

## Unsupported / Not Yet Modeled Explicitly

| XML feature | Status | Current behavior |
| --- | --- | --- |
| distributed controller/worker semantics | unsupported | workload still runs in-process |
| non-S3 driver plugin semantics | unsupported | unknown storage types fail validation/factory creation |
| full percentile/statistics declarations from legacy reporting XML | unsupported | runtime summaries are computed by the Go implementation |

## Real-Workload Tie-In

- `testdata/legacy/s3-config-sample.xml` exercises the common S3 path used for init/prepare/main/cleanup/dispose.
- `testdata/legacy/sio-config-sample.xml` exercises SineIO-oriented multipart preparation and write flow.
- `testdata/workloads/s3-active-subset.xml` and `testdata/workloads/sio-multipart-subset.xml` are the active closure fixtures for the local v1 path.
- `testdata/workloads/mock-stage-aware.xml` locks stage-to-stage continuity for local `mock` runs used in smoke and representative testing.
- `testdata/workloads/xml-compat-storage-subset.xml` locks parser preservation of `siov1` / `gdas` compatibility storage aliases while keeping runtime support deferred.
- `testdata/workloads/xml-range-prefetch-subset.xml` locks parser preservation of prefetch/range-read config keys while execution support remains deferred.
- `testdata/workloads/xml-auth-tolerated-subset.xml` and `testdata/workloads/xml-auth-none-subset.xml` lock parser-tolerated auth-bearing XML shapes while auth modeling remains deferred.
- `testdata/workloads/xml-delay-stage-subset.xml` locks repeated `delay` stage shape and `closuredelay` parsing within the current XML subset.
- `testdata/workloads/xml-splitrw-subset.xml` locks split read/write target-range structure within a single main work.
- `testdata/workloads/mock-reusedata-subset.xml` locks multi-main-stage reuse-data structure under the local `mock` path.
- `testdata/workloads/xml-inheritance-subset.xml` locks config inheritance, storage override, omitted-ratio defaulting, and zero-ratio filtering behavior through the XML path.
- `testdata/workloads/xml-attribute-subset.xml` locks parsed support for `trigger`, `closuredelay`, `interval`, `division`, `rampup`, `rampdown`, and `driver`.
- `testdata/workloads/xml-special-ops-subset.xml` locks representative XML shapes for `delay`, `cleanup`, `localwrite`, and `mfilewrite`.
