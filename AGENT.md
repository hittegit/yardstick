# AGENT.md

## Purpose

Yardstick is a read-only CLI intended to run in CI for any repository and report hygiene findings with deterministic output and policy-driven exit codes.

## CI Contract

- Input target: `-path` points to the repository to scan (default `.`).
- Output formats:
  - `-format table` for human logs
  - `-format json` for machine parsing
- JSON schema stability:
  - top-level: `summary`, `checks`, `findings`, `counts`
  - `checks[]` keys: `check`, `description`, `status`, `level`, `findings`, `why_important`, `how_to_resolve`
  - `findings[]` keys: `check`, `level`, `path`, `message`, `fixed`
  - `counts` keys: `info`, `warn`, `error`
- Exit behavior:
  - exit non-zero when any `error` findings exist
  - with `-strict`, exit non-zero when any `warn` findings exist
  - argument/usage errors also exit non-zero

## Flag Semantics To Preserve

- `-only` must run a subset of checks by key.
- Unknown `-only` keys must fail fast (do not silently pass CI).
- `-format` must accept only `table` or `json`; invalid values must fail.
- `-list` prints available check keys and descriptions.

## Non-Negotiable Guardrails

- Never write into the scanned repository.
- Keep checks local and deterministic (no network calls).
- Keep finding levels constrained to `info`, `warn`, `error`.
- Keep check keys stable once released; downstream CI may parse them.

## Code Map

- `main.go`: CLI parsing, selection/validation, check execution, output, exit policy.
- `internal/checks`: check implementations and registry (`registry.go`).
- `internal/report`: JSON DTO and table renderer.
- `main_test.go`: CLI behavior and policy tests.
- `internal/report/report_test.go`: output/counting contract tests.

## Local Validation

- `go test ./...`
- `go vet ./...`
- `go run . -list`
- `go run . -format json`
- `go run . -format json -strict`

## Linting Requirements

- Configuration and documentation files must remain lint-clean.
- Run markdown lint before committing doc changes:
  - `markdownlint '*.md'`
- Run YAML lint before committing workflow/config changes:
  - `yamllint .golangci.yml .goreleaser.yaml .github/workflows/*.yml`
- Keep `.markdownlint.yaml` and YAML formatting aligned with these checks.
- Do not merge when markdown or YAML lint reports errors.

## When Adding Checks

1. Implement `Check` in `internal/checks/<name>.go`.
2. Register it in `internal/checks/registry.go`.
3. Add focused tests in `internal/checks/<name>_test.go`.
4. Keep findings actionable and low-noise; optimize for CI signal quality.
