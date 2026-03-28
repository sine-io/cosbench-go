# Driver Pages In Unified Service Design

## Goal

Add driver-facing pages to the same Go service so operators can inspect driver health and mission activity without introducing a separate driver-web application.

This project depends on the remote split protocol and reuses the unified route model defined in the roadmap spec.

## User Constraint

The requested scope for this slice is intentionally narrower than full legacy driver-web parity:

- include driver pages
- keep them inside the unified service
- distinguish controller and driver surfaces through routes and APIs

This allows a simpler design than rebuilding the Java `driver-web` package.

## In Scope

- driver HTML routes under `/driver/...`
- driver read models backed by shared remote state
- pages for driver overview, missions, workers, and logs/health summaries
- reuse of existing driver APIs from the remote split project
- navigation and layout updates so controller and driver surfaces coexist cleanly

## Out Of Scope

- separate driver web process or packaging
- legacy route or template parity
- controller page work
- remote protocol creation itself

## Recommended Page Set

Minimum useful driver surface:

- `/driver`
  - driver overview and health
- `/driver/missions`
  - active and recent mission list
- `/driver/missions/:id`
  - mission detail, work-unit summary, event stream, result summary
- `/driver/workers`
  - local worker pool visibility and concurrency status
- `/driver/logs`
  - downloadable or browser-visible driver log summaries

The controller/driver distinction should be obvious in navigation and page headings.

## API Reuse Model

Driver pages should not invent a second data contract.

They should consume the same unified APIs introduced by the remote split, for example:

- `/api/driver/self`
- `/api/driver/missions`
- `/api/driver/missions/:id`
- `/api/driver/workers`
- `/api/driver/logs`

If a page needs data the API does not yet expose, extend the API first and then consume it in the page.

## Read Model Expectations

Driver pages need a focused read layer separate from controller summaries.

Likely read concerns:

- current driver identity, mode, and heartbeat status
- current worker-pool state
- mission queue, lease, and completion state
- per-mission event/log excerpts
- recent upload/retry health

Controller-owned authoritative state can be mirrored, but driver pages should clearly indicate whether a field is:

- local runtime state
- last acknowledged controller state

## UI Direction

Do not recreate the legacy Freemarker driver UI.

The new driver UI should:

- use the current Go template system
- share a unified layout with role-specific navigation
- favor functional tables and detail views over decorative parity

The important migration target is capability, not appearance.

## Testing Strategy

### HTTP/page tests

- route coverage for all driver pages
- expected page content for overview, missions, and workers
- graceful empty-state handling when no missions are active

### API/page integration

- page rendering against realistic driver mission state
- detail pages reflecting retries, failures, and uploads

### Human verification

- run combined mode and confirm driver pages populate
- run controller-only plus driver-only and confirm driver pages reflect live state

## Success Criteria

This slice is complete when:

1. driver pages exist under `/driver/...`
2. pages reuse unified driver APIs instead of bespoke handler-only state assembly
3. operators can inspect driver health, missions, workers, and logs
4. controller and driver navigation can coexist in the same service without role confusion

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
