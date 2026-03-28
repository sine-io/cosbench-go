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

### 1. S3 endpoint default is now aligned, but credentials remain explicit

Legacy S3 has defaults for:

- `endpoint=http://s3.amazonaws.com`
- `path_style_access=false`

Current Go now aligns the endpoint default and still requires:

- explicit `accesskey`
- explicit `secretkey`

This narrows the old gap, but “empty or partial config still initializes” can still differ when credentials are omitted.

### 2. SIO-family path-style policy is now explicit by backend profile

Legacy Java SIO constants define:

- `path_style_access=false` by default

Current Go now does this:

- `sio`: when `path_style_access` is absent, `PathStyle` is forced to `true`
- `siov1` / `gdas`: when `path_style_access` is absent, the legacy `false` default is preserved

This turns the old alias-only behavior into an explicit profile policy. Live comparison still matters for any workload that depended on implicit defaults.

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

### 6. Delete tolerance is now intentionally aligned for common missing-object cases

Legacy S3 delete paths explicitly suppress `404 Not Found` for:

- bucket deletion
- object deletion

Current Go now suppresses common “missing bucket/object” delete errors using the same broad compatibility intent.
Live endpoint checks are still useful to confirm the tolerated error set is wide enough and not overbroad.

### 7. List result shape differs

Legacy Java list implementations build a newline-delimited byte stream of keys.

Current Go returns:

- structured `[]ObjectEntry`
- sorted by key

This is not necessarily a mismatch at the benchmark level, but it is a surface-level difference worth keeping in mind when comparing driver internals or downstream reporting assumptions.

### 8. SIO create-container slash handling is now normalized in the adapter

Legacy Java SIO explicitly splits container names on `/` before bucket creation and keeps only the first segment.

Current Go now normalizes slash-containing bucket names for SIO-family profiles before request construction.

This closes an obvious compatibility gap, though live verification is still useful for real workloads that depend on this quirk.

### 9. Prefetch and range-read request shaping is now present in the local execution path

Legacy SIO samples and code paths support:

- prefetch header-based reads
- explicit byte-range reads

Current Go now:

- parses `is_prefetch`
- parses `is_range_request`, `file_length`, and `chunk_length`
- uses those flags to shape `GetObject` requests in the S3/SIO adapter

The current Go range behavior uses the first chunk implied by `chunk_length` and `file_length`. Live comparison should verify whether that request pattern is sufficient for the intended legacy workloads.

## Recommended Live Check Order

When live endpoint credentials are available, check these in order:

1. `testdata/legacy/sio-config-sample.xml`
   Reason: exercises `mprepare` and `mwrite`, the richest current SIO path
2. `testdata/legacy/s3-config-sample.xml`
   Reason: current mock evidence already shows a large error count, so this is the strongest S3 delta candidate
3. explicit storage-level `part_size` and `restore_days`
   Reason: fallback logic is now aligned in code, so live checks should confirm there are no remaining endpoint-specific surprises
4. cleanup/list-heavy and slash-container scenarios
   Reason: delete tolerance, list surface differences, and bucket normalization should be validated on a real endpoint
5. range/prefetch scenarios
   Reason: request shaping is now implemented but still needs live parity evidence

## Status

These are code-review findings and comparison hypotheses, not live-endpoint conclusions.

They should be used to:

- focus live validation
- explain unexpected differences faster
- avoid treating known scope choices as surprising regressions
