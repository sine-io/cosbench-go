# Remote Smoke Recovery Backend Input Design

## Goal

Extend the dedicated `Remote Smoke Recovery` workflow so it can run either supported backend:

- `s3`
- `sio`

while still keeping recovery fixed as the scenario.

## Problem

The repository now has:

- local recovery smoke for `s3`
- local recovery smoke for `sio`
- a dedicated manual `Remote Smoke Recovery` workflow

But that workflow is currently fixed to:

- `backend=s3`
- `scenario=recovery`

So `sio + recovery` is available only through the more generic `Remote Smoke Local` workflow.

That keeps the recovery entrypoint asymmetric even though recovery parity across `s3` and `sio` is now landed.

## Scope

### In Scope

- add `workflow_dispatch.inputs.backend` to `Remote Smoke Recovery`
- default to `s3`
- allow `sio`
- keep `scenario=recovery` fixed
- keep the existing raw and summary artifact behavior intact

### Out Of Scope

- changing helper behavior
- adding a second recovery workflow
- adding matrix execution to recovery
- changing the default `CI` workflow

## Recommended Approach

Parameterize the existing `Remote Smoke Recovery` workflow rather than introducing a second workflow file.

### Why parameterize this workflow

- the workflow already exists and expresses the right operational intent
- the only missing variable is backend choice
- a second workflow would duplicate almost all of the same YAML
- keeping one named recovery workflow is easier to document and operate

## Workflow Shape

The workflow should remain:

- `workflow_dispatch`
- single job

It should gain:

- `inputs.backend`
  - default: `s3`

The run step should change from hardcoded `SMOKE_REMOTE_LOCAL_BACKEND=s3` to:

- `SMOKE_REMOTE_LOCAL_BACKEND='${{ inputs.backend }}'`

The workflow should still hardcode:

- `SMOKE_REMOTE_LOCAL_SCENARIO=recovery`

## Documentation

README should document both:

- default invocation for `s3`
- explicit `backend=sio` invocation for `sio + recovery`

## Success Criteria

This slice is complete when:

1. `Remote Smoke Recovery` accepts `backend`
2. default behavior still runs `s3 + recovery`
3. `backend=sio` runs `sio + recovery`
4. artifact and summary behavior remain unchanged

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
