GO ?= /snap/bin/go

.PHONY: build compare-local fmt smoke-s3 test tidy validate vet

build:
	$(GO) build ./...

compare-local:
	@echo "== s3-active-subset =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/s3-active-subset.xml -backend mock -json -quiet
	@echo "== mock-stage-aware =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/mock-stage-aware.xml -backend mock -json -quiet
	@echo "== mock-reusedata-subset =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/mock-reusedata-subset.xml -backend mock -json -quiet
	@echo "== xml-splitrw-subset =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/xml-splitrw-subset.xml -backend mock -json -quiet

smoke-s3:
	$(GO) test ./internal/driver/s3 -run Smoke -v

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

vet:
	$(GO) vet ./...

validate: vet test build
