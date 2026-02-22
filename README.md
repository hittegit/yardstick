# Yardstick

![CI](https://github.com/hittegit/yardstick/actions/workflows/ci.yml/badge.svg)

A small Go CLI that scans repositories for hygiene issues and emits readable reports. It works across ecosystems by detecting common manifests, so it is not Go only.

## Overview
Yardstick runs a fast set of repository hygiene checks and renders machine or human friendly reports. It is ecosystem-aware, so you can drop it into Go, Node, Python, Rust, or static site repositories and get useful results with zero configuration. Yardstick is read-only and never writes files.

## Features
- Neutral ecosystem detection via a Manifest check
- Built in hygiene checks for README, LICENSE, .gitignore, and CHANGELOG
- Read-only by default; provides clear guidance to fix issues
- Two output formats, table for humans and JSON for machines
- Strict mode to make warnings fail CI when desired
- Slim GitHub Actions CI workflow with short artifact retention

## Installation
You need Go 1.22 or newer.

- From source in this repo
```bash
# Build and run locally
go run . -h

# Or build a binary
go build -o yardstick .
./yardstick -h
```

- From module path
```bash
# Installs the yardstick binary into your GOPATH/bin or GOBIN
GO111MODULE=on go install github.com/hittegit/yardstick@v0.2.0
yardstick -h
```

### Version
You can print the version embedded at build time:
```bash
yardstick -version
```

## Usage
Run yardstick in any repository root.

```bash
# Human readable table output
yardstick -format table

# JSON output for automation
yardstick -format json

# Fail on warnings too
yardstick -strict

# Only run a subset of checks
yardstick -only readme,license

# List available checks
yardstick -list
```

## What It Checks
- Manifest: Detects common manifests such as `go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`, and more. Reports info on a match, reports a warning if none are found
- README: Ensures `README.md` exists and includes key sections such as Overview, Installation, Usage, CI, and License
- README Links: Validates local README links and markdown anchors for `README.md`
- LICENSE: Ensures `LICENSE` exists and advises adding an appropriate license if missing
- .gitignore: Ensures `.gitignore` exists and advises on sensible defaults if missing
- CHANGELOG: Ensures `CHANGELOG.md` exists and advises adding one if missing
- CODEOWNERS: Ensures repository ownership rules are defined in a standard GitHub CODEOWNERS location
- Security Policy: Ensures `SECURITY.md` exists in a standard GitHub location
- Contributing: Ensures `CONTRIBUTING.md` exists in a standard GitHub location
- CI Workflow: Ensures at least one `.yml` or `.yaml` workflow exists in `.github/workflows`

Yardstick is read-only. It never writes files.

## Output
- Table, compact, greppable, stable column order
- JSON, machine friendly, includes counts by severity

Example table output
```
CHECK       LEVEL  PATH          MESSAGE                                FIXED
changelog   warn   /repo/CHANGELOG.md  CHANGELOG.md missing               false
gitignore   warn   /repo/.gitignore     .gitignore missing                 false
license     warn   /repo/LICENSE        LICENSE missing                    false
manifest    info   /repo/go.mod         Go project detected via go.mod     false
readme      warn   /repo/README.md      Missing section: ## CI             false
```

Example JSON shape
```json
{
  "findings": [
    {"check":"manifest","level":"info","path":"/repo/go.mod","message":"Go project detected via go.mod","fixed":false}
  ],
  "counts": {"info":1, "warn":0, "error":0}
}
```

## Exit Codes
- Non zero exit when errors are present
- With `-strict`, non zero exit when warnings are present
- CLI errors also return a non zero exit

## Using In External CI
Pin Yardstick to a released version in downstream repositories, then run JSON output and enforce your policy in the workflow.

GitHub Actions example:
```yaml
name: quality

on:
  pull_request:
  push:
    branches: [main]

jobs:
  yardstick:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5

      - uses: actions/setup-go@v6
        with:
          go-version: '1.22.x'

      - name: Install pinned yardstick release
        run: go install github.com/hittegit/yardstick@v0.2.0

      - name: Run yardstick and save JSON report
        run: |
          yardstick -path . -format json > yardstick.json
          cat yardstick.json

      - name: Fail build when errors are present
        run: |
          test "$(jq -r '.counts.error' yardstick.json)" = "0"
```

To fail on warnings too, either:
- run `yardstick -strict -format json` and rely on its exit code, or
- assert `.counts.warn == 0` in your workflow script.

## Local Development
- Prerequisites: Go 1.22, a shell, and git
- Repo layout
  - `main.go`: CLI entry point
  - `internal/checks`: Check types and built in checks
  - `internal/report`: JSON and table output

Common tasks
```bash
# Format and vet
gofmt -s -w .
go vet ./...

# Run tests
go test ./...

# Build and try locally
go build -o yardstick .
./yardstick -format table
```

## CI
This repo ships a workflow at `.github/workflows/ci.yml`.
- Builds and pushes a slim CI image to GHCR on push and workflow dispatch
- Lints (golangci-lint), formats, vets, runs tests with coverage, and builds
- Artifacts are retained for 1 day by design

To run CI only when you trigger it manually, change the workflow `on:` block to only `workflow_dispatch:` and push that change. Restore push and pull_request triggers later when ready.

### Dogfooding
CI also runs Yardstick against this repository with `-strict` to ensure we stay compliant with our own checks.

## Adding A New Check
1. Create a new file under `internal/checks`, for example `codeowners.go`
2. Implement the `Check` interface
3. Add an instance to the registry in `internal/checks/registry.go`

Interface shape
```go
// type Check interface {
//   Key() string
//   Description() string
//   Run(ctx context.Context, root string, opts Options) ([]Finding, error)
// }
```

Keep checks fast and deterministic, prefer local file inspection. Yardstick is read-only; checks provide guidance but do not write.

## Design Notes
- Neutral detection is intentional, yardstick should be useful in Go, Node, Python, Rust, Ruby, and static site repos
- Read-only, no writes
- No em dashes, we prefer commas and hyphens in output

## License
MIT, see `LICENSE`.




