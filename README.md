# cosbench-go

Go re-implementation of COSBench with XML workload compatibility (S3 + SIO) and Go best practices.

## Goals
- Full migration (controller + driver) implemented in Go
- Compatible with existing COSBench workload XML format used by `cosbench-sineio`
- Focus storages: **S3** and **SIO**
- Architecture: DDD Lite + CQRS + Clean Architecture + DIP
- Libraries: gin, zerolog, cobra, viper
- Development approach: TDD-first for core parsing + scheduling + metrics

## Status
Bootstrap in progress.

## References
- Legacy project (spec reference): `../cosbench-sineio`
- XML samples:
  - `../cosbench-sineio/release/conf/config-samples/s3-config-sample.xml`
  - `../cosbench-sineio/release/conf/config-samples/sio-config-sample.xml`
