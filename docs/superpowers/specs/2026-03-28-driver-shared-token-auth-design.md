# Driver Shared Token Auth Design

## Goal

Add the thinnest useful authentication layer to the remote controller/driver protocol so driver-side write operations are no longer unauthenticated.

This slice should harden the transport boundary without redesigning the current protocol, changing the browser UX, or introducing a complex credential model.

## Problem

The current remote split protocol is functionally present:

- drivers can register
- drivers can heartbeat
- drivers can claim work
- drivers can upload events and samples
- drivers can complete missions

But every driver write endpoint currently trusts any caller that can reach the HTTP interface. That leaves the controller/driver protocol without even a minimal authenticity check.

## Scope

### In Scope

- one shared token for all drivers
- bearer-token validation on driver write endpoints
- app/config wiring so the controller can require a token
- agent HTTP client wiring so drivers send the token automatically
- combined-mode loopback support with the same token
- tests for missing token, wrong token, and correct token

### Out Of Scope

- per-driver credentials
- token rotation or expiry
- browser login or role-based UI auth
- XML-level auth configuration
- signed requests or mTLS

## Recommended Approach

Use a single shared bearer token for all driver-to-controller writes.

### Why this approach

- smallest change set
- protects the highest-risk surface first
- preserves the current protocol shape
- keeps future migration space open for per-driver or stronger auth later

### Alternatives considered

1. Per-driver token issuance
   Better isolation, but requires controller-side credential lifecycle and new bootstrapping behavior before the protocol is even minimally protected.

2. Browser/session-based auth
   Irrelevant to the current risk. The immediate gap is machine-to-machine protocol writes, not UI access control.

3. No auth until a later “security phase”
   Too weak. The current controller/driver split is already real enough that the write protocol should not remain open.

## Authentication Model

### Token source

The controller receives one configured shared token.

Suggested source of truth:

- environment variable exposed through `app.Config`

For example:

- `COSBENCH_DRIVER_SHARED_TOKEN`

The driver agent receives the same token through configuration and attaches it to every write request.

### Header shape

Use standard bearer auth:

```http
Authorization: Bearer <token>
```

### Endpoints protected in this first slice

Protect only driver write endpoints:

- `POST /api/driver/register`
- `POST /api/driver/heartbeat`
- `POST /api/driver/missions/claim`
- `POST /api/driver/missions/:id/events`
- `POST /api/driver/missions/:id/samples`
- `POST /api/driver/missions/:id/complete`

Do not protect driver read endpoints yet:

- `/api/driver/self`
- `/api/driver/missions`
- `/api/driver/missions/:id`
- `/api/driver/workers`
- `/api/driver/logs`

This keeps the first slice narrowly focused on machine-write protection.

## Runtime Behavior

### Missing token on controller

If the controller is started without a configured shared token:

- remote driver write endpoints should reject requests
- the rejection should be explicit and deterministic

Recommended status:

- `503 Service Unavailable`

Reason: the protocol is unavailable because required auth configuration is absent.

### Missing or malformed Authorization header

- reject with `401 Unauthorized`

### Wrong token

- reject with `403 Forbidden`

### Correct token

- current behavior remains unchanged

## Combined Mode

`combined` mode must continue to work without manual HTTP header assembly in tests.

Recommended behavior:

- app bootstraps one loopback driver agent
- the same configured shared token is injected into that agent’s HTTP client
- combined-mode tests continue to exercise the real HTTP write endpoints

This preserves the current benefit of loopback verification.

## Code Boundaries

### `internal/app`

Responsibilities:

- extend `app.Config` with driver shared token
- read the token from environment or pass it through config
- inject it into combined-mode loopback agent setup

### `internal/web`

Responsibilities:

- one reusable auth helper for driver write handlers
- no protocol business logic should move into the auth helper

### `internal/driver/agent`

Responsibilities:

- HTTP client automatically sets bearer header on write requests
- no call site should have to assemble auth headers manually

### `cmd/server`

Responsibilities:

- expose or document the token configuration path

## Testing Strategy

### Web/controller tests

Add tests proving:

- missing token configuration causes protected endpoints to reject
- missing `Authorization` header yields `401`
- wrong token yields `403`
- correct token preserves existing success behavior

### Agent tests

Add tests proving:

- register/claim/report fails with no token when the controller requires one
- register/claim/report succeeds with the correct token

### App tests

Add tests proving:

- `combined` mode still processes a mission when the token is configured

## Success Criteria

This slice is complete when:

1. all driver write endpoints require bearer auth
2. the controller fails closed when no shared token is configured
3. the driver agent automatically sends the configured token
4. combined mode still works through the real HTTP loopback path
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
