# Legacy Comparison Matrix

This document records the current, repeatable comparison view between `cosbench-go` and the legacy `cosbench-sineio` project for the repository's representative workload subset.

Companion code-level notes live in `docs/storage-driver-comparison-notes.md`.
Execution guidance for future live runs lives in `docs/legacy-live-run-checklist.md`.

## Result Labels

- `match`: behavior is believed to align closely enough for the local-only v1 scope, based on current evidence
- `acceptable delta`: behavior or artifact shape differs, but the difference is consistent with the documented local-only scope
- `mismatch`: evidence shows a concrete behavioral gap that should be investigated
- `not yet run`: no direct comparison evidence has been captured yet

## Comparison Dimensions

Each row uses the same dimensions:

- XML parse outcome
- normalized workload shape
- accepted storage/backend path
- execution outcome category
- result surface availability (`CLI`, `JSON`, `CSV`)
- notable semantic differences or missing evidence

## Matrix

| Fixture | Legacy Reference Status | `cosbench-go` Status | Comparison Class | Result | Notes |
| --- | --- | --- | --- | --- | --- |
| `testdata/legacy/s3-config-sample.xml` | direct legacy sample from `../cosbench-sineio/release/conf/config-samples/s3-config-sample.xml` | parser-covered; local CLI run completed on 2026-03-24 with `-backend mock`; JSON summary produced | runnable with live endpoint setup | acceptable delta | mock-override evidence: 5 stages / 5 works / 1,508,071 samples / 1,206,814 errors. This confirms current parser + execution path can ingest the legacy sample, but the mixed read/write workload shows high error volume under mock-backed local execution, so direct legacy/live comparison is still pending |
| `testdata/legacy/sio-config-sample.xml` | direct legacy sample from `../cosbench-sineio/release/conf/config-samples/sio-config-sample.xml` | parser-covered; local CLI run completed on 2026-03-24 with `-backend mock`; JSON summary produced | runnable with live endpoint setup | acceptable delta | mock-override evidence: 2 stages / 2 works / 268,090 samples / 0 errors. Strongest current candidate for a first live legacy-side comparison of `mprepare` + `mwrite` |
| `testdata/workloads/s3-active-subset.xml` | go-curated representative derived from active legacy semantics | parser-covered; local `make compare-local` run completed on 2026-03-25 with `-backend mock`; JSON summary produced | runnable now | acceptable delta | mock-backed local compare evidence: 2 stages / 2 works / 203,792 samples / 142,691 errors. This fixture has no prepare stage, so mock-backed read failures are expected local evidence rather than proof of legacy divergence |
| `testdata/workloads/sio-multipart-subset.xml` | go-curated representative derived from active legacy semantics | parser-covered; local execution path and opt-in live smoke support exist | runnable with live endpoint setup | acceptable delta | narrower than the legacy sample but aligned with the active multipart path |
| `testdata/workloads/mock-stage-aware.xml` | go-curated representative focused on local stage continuity | local `make compare-local` run completed on 2026-03-25 with `-backend mock`; JSON summary produced | runnable now | match | mock-backed local compare evidence: 6 stages / 6 works / 8 samples / 0 errors. Confirms stage-aware mock continuity is behaving as intended |
| `testdata/workloads/mock-reusedata-subset.xml` | go-curated representative derived from legacy reuse-data structure | local `make compare-local` run completed on 2026-03-25 with `-backend mock`; JSON summary produced | runnable now | acceptable delta | mock-backed local compare evidence: 6 stages / 6 works / 7 samples / 0 errors. Useful local proof that prepared data survives across multiple main stages |
| `testdata/workloads/xml-inheritance-subset.xml` | no direct legacy sample; compares against documented inheritance behavior | parser and normalization covered | parser-only comparison | acceptable delta | locks workload/workflow/stage/work/op config inheritance, storage override, omitted-ratio defaulting, zero-ratio filtering |
| `testdata/workloads/xml-attribute-subset.xml` | no direct legacy sample; compares against documented XML model | parser and normalization covered | parser-only comparison | acceptable delta | locks `trigger`, `closuredelay`, `interval`, `division`, `rampup`, `rampdown`, `driver` |
| `testdata/workloads/xml-special-ops-subset.xml` | informed by legacy sample family and SineIO changelog | parser and normalization covered | parser-only comparison | acceptable delta | locks XML shapes for `delay`, `cleanup`, `localwrite`, `mfilewrite`; not yet directly compared against a live legacy run |
| `testdata/workloads/xml-splitrw-subset.xml` | go-curated representative derived from legacy split read/write shape | local `make compare-local` run completed on 2026-03-25 with `-backend mock`; JSON summary produced | runnable now | acceptable delta | mock-backed local compare evidence: 5 stages / 5 works / 8 samples / 1 error. The read side expects data in a separate container range from write, so mock-only evidence is informative but not a legacy parity verdict |

## Current Known Deltas

- `cosbench-go` is intentionally local-only for this phase, so distributed controller/worker behavior is an expected delta, not a mismatch.
- Report surfaces are comparable at a summary/export level, but they are not expected to be byte-identical to the legacy Java outputs.
- Some repository fixtures are curated representatives of active legacy semantics rather than one-to-one copies of legacy XML files.
- The legacy S3 sample currently shows a large error count when run through `cosbench-go` with `-backend mock`; this is evidence worth preserving, not hiding, until a direct live legacy comparison explains whether the behavior is expected or divergent.

## Runbook

### 1. Confirm parser-level behavior

Run:

```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/infrastructure/xml ./internal/workloadxml
```

Use this for rows marked `parser-only comparison`.

### 2. Run `cosbench-go` locally on a representative fixture

Example:

```bash
GO=$(which go || echo /snap/bin/go) go run ./cmd/cosbench-go -workload testdata/workloads/s3-active-subset.xml -backend mock -json
```

This is useful for checking normalized shape, CLI summary, and local execution category without requiring live credentials.

### 2a. Run the curated local comparison set

Run:

```bash
GO=$(which go || echo /snap/bin/go) make compare-local
```

This is the fastest way to refresh the safe mock-backed fixture set in one pass.
It currently covers:

- `s3-active-subset.xml`
- `mock-stage-aware.xml`
- `mock-reusedata-subset.xml`
- `xml-splitrw-subset.xml`

### 3. Run live endpoint smoke coverage

Set:

- `COSBENCH_SMOKE_ENDPOINT`
- `COSBENCH_SMOKE_ACCESS_KEY`
- `COSBENCH_SMOKE_SECRET_KEY`

Optional:

- `COSBENCH_SMOKE_BACKEND`
- `COSBENCH_SMOKE_REGION`
- `COSBENCH_SMOKE_PATH_STYLE`
- `COSBENCH_SMOKE_BUCKET_PREFIX`

Then run:

```bash
GO=$(which go || echo /snap/bin/go) make smoke-s3
```

This confirms live adapter connectivity and the minimal object lifecycle. It does not replace workload-level comparison by itself.

For the full live-run sequence and recording checklist, use:

```bash
docs/legacy-live-run-checklist.md
```

### 4. Locate legacy references

Primary legacy sample directory:

```bash
../cosbench-sineio/release/conf/config-samples/
```

Useful files there include:

- `s3-config-sample.xml`
- `sio-config-sample.xml`
- `sio-config-smoke-test.xml`
- `delay-stage-config-sample.xml`

For code-level driver differences to keep in mind during live checks, also review:

```bash
docs/storage-driver-comparison-notes.md
```

### 5. Record findings

When a row is actually exercised against a live endpoint or legacy run, update:

- `Legacy Reference Status`
- `cosbench-go Status`
- `Result`
- `Notes`

Keep the row explicit even when the result is still `not yet run`.
