# Storage Driver Comparison Notes

This document captures code-level comparison notes between the current `cosbench-go` S3/SIO driver path and the legacy Java implementation in `../cosbench-sineio`.

It is meant to support live comparison work by identifying the most likely semantic deltas before endpoint runs begin.

## Source References

Current Go implementation:

- `internal/driver/s3/config.go`
- `internal/driver/s3/adapter.go`

Legacy Java implementation:

- `../cosbench-sineio/dev/cosbench-s3/src/com/intel/cosbench/api/S3Stor/S3Storage.java`
- `../cosbench-sineio/dev/cosbench-s3/src/com/intel/cosbench/client/S3Stor/S3Constants.java`
- `../cosbench-sineio/dev/cosbench-sineio/src/com/intel/cosbench/api/sio/SIOStorage.java`
- `../cosbench-sineio/dev/cosbench-sineio/src/com/intel/cosbench/client/sio/SIOConstants.java`

## High-Value Comparisons

### 1. S3 required config is stricter in `cosbench-go`

Legacy S3 has defaults for:

- `endpoint=http://s3.amazonaws.com`
- `path_style_access=false`

Current Go requires:

- explicit `endpoint`
- explicit `accesskey`
- explicit `secretkey`

This is acceptable for the current repository, but it is still a real behavioral delta when comparing “empty or partial config still initializes” behavior.

### 2. SIO path-style default differs

Legacy Java SIO constants define:

- `path_style_access=false` by default

Current Go does this instead:

- when backend is `sio` and `path_style_access` is absent, `PathStyle` is forced to `true`

This is a deliberate convenience choice in `cosbench-go`, but it means a live comparison should explicitly record whether a given workload relied on the legacy default or always set the field.

### 3. `storage_class` parity is now present for both single-part and multipart uploads

Legacy SIO applies `storage_class` from storage config to:

- normal object creation
- multipart upload creation

Current Go now does the same:

- `PutObjectInput`
- `CreateMultipartUploadInput`

This was previously a gap for multipart upload and has now been closed.

### 4. `part_size` storage fallback is now aligned for execution defaults

Legacy SIO loads `part_size` from storage config during adapter initialization and uses that field during multipart upload.

Current Go now merges config sources for execution-time defaults:

- adapter config parsing
- storage config fallback at execution time
- operation config override at execution time

That means:

- storage-level `part_size` now drives multipart execution when operation config omits it
- operation-level `part_size` still overrides the storage-level value

Live comparison is still useful here, but this is no longer an unhandled execution-path mismatch.

### 5. `restore_days` storage fallback is now aligned for execution defaults

Legacy SIO loads `restore_days` from storage config and uses the initialized field during restore operations.

Current Go now merges config sources for execution-time defaults:

- storage-level `restore_days` is used when operation config omits it
- operation-level `restore_days` still overrides the storage-level value

This reduces the previous source-layer mismatch to a live-verification question rather than an obvious implementation gap.

### 6. Delete tolerance differs

Legacy S3 delete paths explicitly suppress `404 Not Found` for:

- bucket deletion
- object deletion

Current Go directly returns the SDK error for:

- `DeleteBucket`
- `DeleteObject`

`CreateBucket` already tolerates “already exists/already owned” style responses, but delete behavior is stricter. This may show up as a live difference during cleanup-heavy workloads.

### 7. List result shape differs

Legacy Java list implementations build a newline-delimited byte stream of keys.

Current Go returns:

- structured `[]ObjectEntry`
- sorted by key

This is not necessarily a mismatch at the benchmark level, but it is a surface-level difference worth keeping in mind when comparing driver internals or downstream reporting assumptions.

### 8. SIO create-container behavior may differ on slash-containing container names

Legacy Java SIO explicitly splits container names on `/` before bucket creation and keeps only the first segment.

Current Go passes the bucket string through unchanged.

If any real workload depends on slash-containing container names, this is another high-priority live comparison candidate.

## Recommended Live Check Order

When live endpoint credentials are available, check these in order:

1. `testdata/legacy/sio-config-sample.xml`
   Reason: exercises `mprepare` and `mwrite`, the richest current SIO path
2. `testdata/legacy/s3-config-sample.xml`
   Reason: current mock evidence already shows a large error count, so this is the strongest S3 delta candidate
3. explicit storage-level `part_size` and `restore_days`
   Reason: fallback logic is now aligned in code, so live checks should confirm there are no remaining endpoint-specific surprises
4. cleanup/list-heavy scenarios
   Reason: delete tolerance and list surface differences may appear there

## Status

These are code-review findings and comparison hypotheses, not live-endpoint conclusions.

They should be used to:

- focus live validation
- explain unexpected differences faster
- avoid treating known scope choices as surprising regressions
