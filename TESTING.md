# Testing Strategy

This document defines how we test Yardstick locally and in CI to keep quality high, feedback fast, and results reproducible.

## Goals

- Fast feedback: sub‑minute local cycle for lint, vet, and unit tests.
- Deterministic: no network, no time‑of‑day dependencies, stable outputs.
- CI parity: native runner and container parity to catch environment drift.
- Dogfooding: run Yardstick against itself with `-strict`.
- Release confidence: reproducible builds and basic binary smoke tests.

## Test Types

- Unit tests (Go):
  - Packages: `internal/checks`, `internal/report`.
  - Use `t.TempDir()` for any filesystem writes, never touch the repo tree.
  - Avoid networking and environment‑dependent behavior.
- Integration tests (CLI):
  - Exercise `run(context.Background())` via flags for exit behavior.
  - Verify `-strict` fails on warnings, `-only` filters, and `-format` outputs.
- Static analysis:
  - `golangci-lint` with `govet`, `staticcheck`, `errcheck`, `gosec`, etc.
- Formatting and module hygiene:
  - `gofmt -s` for formatting; `go mod tidy` clean.
- Container parity (optional for PRs):
  - Build and test inside `ghcr.io/hittegit/yardstick-ci:go1.22-bookworm-slim`.
- Dogfooding:
  - Run `yardstick -strict` on this repo to enforce our own rules.

## Local Workflow

- Prereqs: Go 1.22+, optional `golangci-lint`, optional Docker.

Suggested loop:

```bash
# Format, vet, lint, test, build
make fmt vet lint test build

# or without Makefile
gofmt -s -w .
go vet ./...
golangci-lint run
go test ./... -cover
go build -trimpath -buildvcs=false -o yardstick .
```

Additional checks:

```bash
# Race detector (slower, run before releases)
go test -race ./...

# Coverage HTML report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Containerized parity run (mirrors CI image)
docker build -t ys-ci -f ci/Dockerfile .
docker run --rm -v "${PWD}":/work -w /work ys-ci sh -lc \
  "gofmt -s -w . && go vet ./... && go test ./... -cover && go build -buildvcs=false ./..."
```

CLI smoke tests:

```bash
# Table and JSON output
go run . -format table
go run . -format json

# Strict mode exit (expect non-zero if warnings exist)
# PowerShell:  go run . -strict; $LASTEXITCODE
# Bash:        go run . -strict; echo $?

# Only a subset of checks
go run . -only readme,license -format table

# CLI in a temp dir (read-only)\n```bash\nTMP=$(mktemp -d); go run . -path "$TMP" -format table; ls -la "$TMP"\n```

## CI Strategy
- GitHub Actions workflow: `.github/workflows/ci.yml`
  - Native tests (all events):
    - Setup Go 1.22, run `golangci-lint`, `gofmt`, `go vet`, `go test` (with coverage), `go build -buildvcs=false`.
  - Container tests (main/manual):
    - Build/push CI image, run the same steps inside the image for parity.
  - Dogfooding (self‑check):
    - Run `go run . -format json -strict` on this repo and upload the report.
  - Artifacts: coverage (`coverage.out`) and self‑report retained for 1 day.

Recommended CI enhancements (optional):
- Race detector on main/manual builds:
  - `go test -race ./...` as a separate step.
- Module tidy check:
  - `go mod tidy` then `git diff --exit-code` to ensure no drift.
- Verify JSON schema stability:
  - Add a golden file test for `report.Output` shape if schema stability becomes critical.
- CodeQL (optional):
  - Enable `github/codeql-action` for Go on a nightly schedule.

## Release Testing
- Build reproducible binaries with GoReleaser via tag `vX.Y.Z`.
- ldflags inject version, commit, and date; verify with `yardstick -version`.
- Snapshot builds locally:
```bash
goreleaser release --snapshot --clean
./dist/yardstick_*/*/yardstick -version
```

## Adding Tests

- For a new check `internal/checks/foo.go`:
  - Add `foo_test.go` with table‑driven cases.
  - Use `t.TempDir()` for files and assert findings precisely (level, path, message, fixed flag).
  - Register the check in `internal/checks/registry.go` and extend any integration tests if needed.

## Exit Criteria for PRs

- Lint clean (`golangci-lint run` passes).
- `go test ./...` passes; coverage does not decrease materially.
- `go vet ./...` clean.
- Self‑check report contains no errors and only acceptable warnings.
- No unformatted files or module drift after `gofmt -s -w .` and `go mod tidy`.

## Future End‑to‑End (optional)

- Cross‑repo smoke tests against public sample repos (e.g., a static site):
  - Matrix of target repos, run `yardstick -format json` and collect reports as artifacts.
  - Gate only on fatal errors by default; warnings are informative.
