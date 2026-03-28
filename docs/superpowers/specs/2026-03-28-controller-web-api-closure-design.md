# Controller Web/API Closure Design

## Goal

Expand the controller-facing capability of `cosbench-go` so the unified Go service functionally covers the deferred controller web surfaces without reproducing the legacy Java UI structure verbatim.

The target is functional parity through a new unified controller API and a focused set of controller pages.

## Problems To Solve

Current controller coverage is intentionally narrow:

- dashboard
- workload upload
- endpoint management
- job detail
- history
- JSON/CSV result export

Legacy controller-web exposes a wider surface:

- matrix overview
- workload config inspection
- advanced config inspection
- stage-focused drilldown
- timeline pages and CSV
- Prometheus export
- log/config downloads

These capabilities must be added without rebuilding the legacy Freemarker application shape.

## In Scope

- controller API routes under `/api/controller/...`
- controller HTML routes under `/controller/...`
- matrix, config, advanced config, stage detail, timeline, timeline CSV, Prometheus, download-log, and download-config capability
- persistent controller-owned timeline/query data needed by those surfaces
- focused read-model/query helpers for controller pages and APIs

## Out Of Scope

- pixel-identical legacy page reproduction
- separate controller web application packaging
- remote split protocol itself
- driver pages

## Recommended Approach

Build controller capability around an API-first read model and let pages consume that same contract.

### 1. Add controller read models instead of page-specific object digging

The current web handler reaches directly into `Manager` and raw domain structs. That is enough for small pages, but not for matrix, timeline, advanced config, and Prometheus surfaces.

Add controller query/read helpers that expose:

- job matrix summaries
- normalized config views
- per-stage work summaries
- timeline series
- artifact references

Pages and APIs should both consume these read models.

### 2. Persist timeline-ready execution data under controller ownership

Timeline pages cannot be retrofitted cleanly from summary-only state. The controller needs to persist timeline-oriented samples or pre-aggregated time buckets during execution.

Recommended shape:

- controller stores per-job timeline buckets
- buckets can be queried by stage and optionally operation/work
- timeline CSV and Prometheus exporters read the same data source

This same timeline persistence becomes reusable when remote drivers later upload samples.

### 3. Reorganize controller routes into a coherent surface

Suggested HTML surface:

- `/controller`
- `/controller/matrix`
- `/controller/jobs/:id/config`
- `/controller/jobs/:id/config/advanced`
- `/controller/jobs/:id/stages/:stage`
- `/controller/jobs/:id/timeline`

Suggested API surface:

- `/api/controller/jobs`
- `/api/controller/jobs/:id`
- `/api/controller/jobs/:id/config`
- `/api/controller/jobs/:id/config/advanced`
- `/api/controller/jobs/:id/stages/:stage`
- `/api/controller/jobs/:id/timeline`
- `/api/controller/jobs/:id/timeline.csv`
- `/api/controller/jobs/:id/artifacts/log`
- `/api/controller/jobs/:id/artifacts/config`
- `/api/controller/metrics/prometheus`

### 4. Treat advanced config as an explanation view, not a legacy form clone

The new advanced config surface should explain the normalized execution shape:

- inherited storage and auth
- effective work config
- effective operation config
- explicit overrides vs inherited values

This is more useful in Go than mimicking the legacy matrix form builder.

### 5. Treat downloads as artifacts, not ad hoc handlers

Config and log downloads should be modeled as controller-managed artifacts so the same concept can later cover:

- raw submitted XML
- normalized config export
- controller event log export
- remote mission logs

## Query And Storage Additions

Expected new controller-side persisted/read state:

- timeline buckets or samples
- normalized-config export payloads
- artifact metadata for downloadable controller outputs

The storage remains file-backed in this phase.

## Testing Strategy

### Query/model tests

- matrix summarization
- advanced-config explanation output
- timeline bucket aggregation

### HTTP tests

- controller HTML route coverage
- controller JSON route coverage
- CSV export and artifact download coverage
- Prometheus content-type and payload coverage

### End-to-end checks

- submit workload
- run workload locally
- inspect controller matrix and stage pages
- fetch timeline CSV and Prometheus output

## Success Criteria

This slice is complete when:

1. all deferred controller capabilities have a new unified route/API equivalent
2. pages use controller read models rather than raw ad hoc object assembly
3. timeline data is persisted and queryable
4. Prometheus export is available from the unified service
5. config and log downloads are exposed as controller-managed artifacts

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
