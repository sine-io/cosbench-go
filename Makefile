GO ?= /snap/bin/go

.PHONY: build fmt smoke-s3 test tidy validate vet

build:
	$(GO) build ./...

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
