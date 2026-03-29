# Smoke Ready Validate Workflow Design

## Goal

Add a manual GitHub Actions workflow that runs the repository's `smoke-ready` schema validation entrypoint on GitHub-hosted runners.

## Problem

The repository now has:

- `make --no-print-directory smoke-ready-validate`
- `make --no-print-directory smoke-ready-validate-json`
- `docs/smoke-ready.schema.json`

But there is no stable remote entrypoint proving that the validator itself works in the GitHub-hosted environment.

## Desired Behavior

Add one manual workflow:

- name: `Smoke Ready Validate`
- trigger: `workflow_dispatch`
- run `make --no-print-directory smoke-ready-json`
- run `make --no-print-directory smoke-ready-validate-json`
- upload a small artifact directory
- write the validation JSON into the job summary

## Scope

In scope:

- one workflow file
- one focused workflow contract test
- one short README note

Out of scope:

- default CI integration
- changes to `smoke_ready.py`
- changes to the schema
- additional matrix variants

## Design

The workflow should stay minimal:

- checkout repository
- set `GH_TOKEN: ${{ github.token }}`
- capture `.artifacts/smoke-ready-validate/smoke-ready.json`
- capture `.artifacts/smoke-ready-validate/validation.json`
- upload `.artifacts/smoke-ready-validate/`
- append `validation.json` to `$GITHUB_STEP_SUMMARY`

This keeps the workflow focused on contract verification rather than broader smoke execution.

## Acceptance Criteria

- manual workflow exists and is runnable
- artifact includes both raw payload and validation result
- job summary includes the validation JSON
- existing runtime behavior remains unchanged
