SHELL := /bin/bash
.PHONY: all build test lint fmt vet ci release snapshot

all: lint test build

build:
	go build -trimpath -buildvcs=false -o yardstick .

test:
	go test ./... -cover

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

vet:
	go vet ./...

ci: fmt vet test

snapshot:
	goreleaser release --snapshot --clean
