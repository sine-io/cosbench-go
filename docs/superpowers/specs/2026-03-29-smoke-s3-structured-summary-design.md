# Smoke S3 Structured Summary Design

## Goal

Upgrade `Smoke S3` and `Smoke S3 Matrix` so they emit structured smoke result summaries instead of only raw text artifacts.

## Problem

Right now:

- `Smoke S3` uploads only `smoke-s3-output.txt`
- `Smoke S3 Matrix` row artifacts only contain the same raw text
- the matrix aggregate marks rows merely as `present` or `missing`
- `smoke-ready` has to infer `executed` vs `skipped` by parsing text output

That works, but it is brittle and duplicates smoke-result parsing logic outside the workflows that generated the evidence.

## Desired Behavior

Each `Smoke S3` run should emit:

- raw text output
- a structured `summary.json`

with a normalized result such as:

- `executed`
- `skipped`
- `failed`

Each `Smoke S3 Matrix` row should do the same, and the aggregate script should consume those structured row summaries rather than guessing from raw text.

## Scope

In scope:

- one small summary script for raw smoke output
- `smoke-s3.yml`
- `smoke-s3-matrix.yml`
- `aggregate_smoke_s3_matrix.py`
- workflow/aggregation tests
- small README note

Out of scope:

- changing the Go smoke tests themselves
- changing `smoke-ready` in this round
- changing legacy live or remote smoke workflows

## Design

Add `scripts/summarize_smoke_s3_output.py` that accepts:

- input text file path
- output summary path

and classifies the output:

- all smoke tests skipped => `skipped`
- smoke tests executed and passed => `executed`
- anything else => `failed`

Update `Smoke S3` workflow to:

- keep `smoke-s3-output.txt`
- generate `.artifacts/smoke-s3-summary/summary.json`
- upload both together
- include summary content in the job summary when available

Update `Smoke S3 Matrix` row jobs to do the same per backend.

Update `aggregate_smoke_s3_matrix.py` to prefer structured row summaries and render row statuses as:

- `executed`
- `skipped`
- `failed`
- `missing`

## Acceptance Criteria

- `Smoke S3` uploads raw text plus structured summary
- `Smoke S3 Matrix` row artifacts include structured summaries
- matrix aggregate consumes structured row summaries
- aggregate markdown/status reflects executed vs skipped vs failed
- existing raw text output remains available
