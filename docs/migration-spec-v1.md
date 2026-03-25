# COSBench-SineIO → Go Migration Spec v1

## Goal

Re-implement `cosbench-sineio` in Go with **behavioral compatibility first** on a local-only v1 boundary:

- workload XML compatibility for the active subset used in this repository
- S3 + SIO storage compatibility
- single-process control plane + driver execution flow
- compatible core operation semantics
- compatible benchmark-level reporting through UI, JSON, CSV, and CLI summaries

This project is **not** a Java-to-Go line-by-line translation.
It is a Go-native implementation that treats the legacy project as the behavioral reference.

## Scope for v1

### In scope
- Workload XML parsing for the active subset exercised by repository fixtures
- Config inheritance across workload / workflow / stage / work / operation
- Storage types: `s3`, `sio`
- Core and SineIO operations:
  - `init`
  - `prepare`
  - `mprepare`
  - `write`
  - `mwrite`
  - `mfilewrite`
  - `localwrite`
  - `read`
  - `delete`
  - `cleanup`
  - `dispose`
  - `head`
  - `list`
  - `restore`
  - `delay`
- Driver execution engine
- Local web control plane, snapshot persistence, and job lifecycle
- JSON / CSV / CLI reporting

### Out of scope for v1
- Remote controller/driver split, worker registration, or mission dispatch
- Legacy Freemarker UI parity
- Non-S3/SIO storages
- Full historical controller web behavior parity
- Pixel-identical chart pages
- Full historical XML breadth outside the active migration subset

## Legacy model observed from `cosbench-sineio`

## XML structure

Legacy XML can include elements such as `<auth>`, but the local v1 closure focuses on the active subset exercised by this repository's fixtures and endpoint-backed configuration path.

```xml
<workload name="..." description="..." trigger="..." config="...">
  <auth ... />?
  <storage type="..." config="..." />?
  <workflow config="...">
    <workstage name="..." closuredelay="..." trigger="..." config="...">
      <auth ... />?
      <storage type="..." config="..." />?
      <work name="..." type="..." workers="..." interval="..." division="..."
            runtime="..." rampup="..." rampdown="..." afr="..."
            totalOps="..." totalBytes="..." driver="..." config="...">
        <auth ... />?
        <storage type="..." config="..." />?
        <operation type="..." ratio="..." division="..." config="..." id="..." />*
      </work>
    </workstage>+
  </workflow>
</workload>
```

## Inheritance rules

Legacy Java behavior uses `ConfigUtils.inherit(child, parent)` across levels.
For the local v1 Go implementation the behavior target is:

- workload.config is inherited into workflow.config
- workflow.config is inherited into stage.config
- stage.config is inherited into work.config
- work.config is inherited into operation.config
- stage-local `<storage>` overrides workload default storage
- work-local `<storage>` overrides stage/default storage

Explicit `<auth>` inheritance is observed in legacy XML but remains outside the active local v1 closure.

## Work type normalization

Legacy Java rewrites several special work types into a normal work containing one operation.
This is preserved in Go.

### `init`
- name defaults to `init`
- division becomes `container`
- runtime = 0
- totalBytes = 0
- totalOps = workers
- synthetic operation:
  - type = `init`
  - ratio = 100
  - config prepends: `objects=r(0,0);sizes=c(0)B`

### `dispose`
- same shape as `init`, but operation type `dispose`

### `prepare`
- name defaults to `prepare`
- division becomes `object`
- runtime = 0
- totalBytes = 0
- totalOps = workers
- synthetic operation:
  - type = `prepare`
  - ratio = 100
  - prepend `createContainer=false` if absent

### `mprepare`
- same as `prepare`, but operation type `mprepare`
- only valid for storage type: `sio`, `siov1`, `gdas`
- in Go v1 we execute `sio`; `siov1` and `gdas` may still be parsed as compatibility aliases

### `cleanup`
- name defaults to `cleanup`
- division becomes `object`
- runtime = 0
- totalBytes = 0
- totalOps = workers
- synthetic operation:
  - type = `cleanup`
  - ratio = 100
  - prepend `deleteContainer=false` if absent

### `delay`
- name defaults to `delay`
- division becomes `none`
- workers = 1
- runtime = 0
- totalBytes = 0
- totalOps = 1
- synthetic operation:
  - type = `delay`
  - ratio = 100

## Validation rules to preserve

- workload name must not be empty
- stage name must not be empty
- work name must not be empty after normalization
- workers must be > 0
- operation ratio must be 0..100
- if runtime == 0 and totalOps == 0 and totalBytes == 0 => invalid
- sum of operation ratios after filtering zero-ratio ops must equal 100
- if `totalOps > 0`, then `workers <= totalOps`
- storage type required where execution needs it
- obvious non-runnable jobs should fail preflight before execution starts

## SineIO-specific operation restrictions observed in legacy Java

Legacy work validation rejects these operations unless storage type is one of `sio`, `siov1`, `gdas`:

- `mwrite`
- `head`
- `restore`
- `mfilewrite`
- `localwrite`

For Go v1 target behavior:
- fully support under `sio`
- accept legacy aliases where practical

## Storage config parsing

Legacy `KVConfigParser` behavior is simple:
- split on `;`
- split each entry on first `=`
- trim spaces around key/value
- empty string allowed
- no escaping rules in legacy implementation

Go v1 preserves this baseline first.
Advanced escaping can be added later if real workloads require it.

## Sample compatibility cases already identified

### S3 sample
- workload-level default storage
- stage-level `init` / `prepare` / `cleanup` / `dispose`
- main mixed-operation work with ratios
- runtime-based execution

### SIO sample
- stage-local storage override
- `mprepare`, `mwrite`, `mfilewrite`, `localwrite`, `restore`, `head`
- `path_style_access`
- `part_size`
- `restore_days`
- `storage_class`
- `no_verify_ssl`
- `aws_region`
- size units using `KiB`

## Local v1 completion milestones

1. Domain model for workload / workflow / stage / work / operation / storage
2. XML parser with sample-based tests
3. Config inheritance, validation, and special work normalization
4. Storage config parser
5. Storage adapter interfaces and S3/SIO adapter wiring
6. Local driver execution engine
7. Local control plane, snapshots, and exports
8. Preflight validation and work-level diagnostics
