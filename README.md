# Yardstick

A small Go CLI that scans repositories for hygiene issues, emits readable reports, and can scaffold safe defaults on request. It works across ecosystems by detecting common manifests, so it is not Go only.

## Features
- Neutral ecosystem detection via a Manifest check
- Built in hygiene checks for README, LICENSE, .gitignore, and CHANGELOG
- Auto fix mode to create safe defaults, never overwrites existing files
- Two output formats, table for humans and JSON for machines
- Strict mode to make warnings fail CI when desired
- Slim GitHub Actions CI workflow with short artifact retention

## Install
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
GO111MODULE=on go install github.com/hittegit/yardstick@latest
yardstick -h
```

## Quick Start
Run yardstick in any repository root.

```bash
# Human readable table output
yardstick -format table

# JSON output for automation
yardstick -format json

# Fail on warnings too
yardstick -strict

# Create safe defaults for common files if missing
yardstick -fix

# Write a few common support files if missing
yardstick -scaffold

# Only run a subset of checks
yardstick -only readme,license

# List available checks
yardstick -list
```

## What It Checks
- Manifest: Detects common manifests such as `go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`, and more. Reports info on a match, reports a warning if none are found
- README: Ensures `README.md` exists and includes key sections such as Overview, Installation, Usage, CI, and License
- LICENSE: Ensures `LICENSE` exists. With `-fix`, creates a permissive MIT license that you can later replace
- .gitignore: Ensures `.gitignore` exists with sensible defaults. With `-fix`, creates a starter file
- CHANGELOG: Ensures `CHANGELOG.md` exists. With `-fix`, creates a minimal starter file

All fixes are safe. Yardstick never overwrites existing files.

## Output Formats
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

## Local Development
- Prerequisites: Go 1.22, a shell, and git
- Repo layout
  - `main.go`: CLI entry point
  - `internal/checks`: Check types and built in checks
  - `internal/report`: JSON and table output
  - `internal/scaffold`: Safe file scaffolding

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

## GitHub Actions CI
This repo ships a workflow at `.github/workflows/ci.yml`.
- Builds and pushes a slim CI image to GHCR on push and workflow dispatch
- Runs formatting check, vet, tests with coverage, and a build inside that image
- Artifacts are retained for 1 day by design

To run CI only when you trigger it manually, change the workflow `on:` block to only `workflow_dispatch:` and push that change. Restore push and pull_request triggers later when ready.

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

Keep checks fast and deterministic, prefer local file inspection. If a fix is safe, gate it behind `Options{AutoFix: true}` only.

## Design Notes
- Neutral detection is intentional, yardstick should be useful in Go, Node, Python, Rust, Ruby, and static site repos
- Auto fixes are safe only, no overwrites
- No em dashes, we prefer commas and hyphens in output

## License
MIT, see `LICENSE`.
