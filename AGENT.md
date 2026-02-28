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

## Delivery Workflow

- Delivery workflow is defined in `push_code.md` and is the single source of truth.
- For non-trivial work, follow the issue-first process in that file:
  - issue -> issue-numbered branch -> PR -> merge to `main` -> release tag
- Do not duplicate or contradict workflow steps in this file.
- If delivery process changes, update `push_code.md` first and keep this reference current.

## Downstream Sync Policy

- Downstream repos (for example `dharma-siblings`) should pin Yardstick to a release tag, not `latest`.
- Keep downstream repos in sync using automated update PRs (Renovate regex manager for `go install ...@vX.Y.Z` lines).
- Update flow for downstream repos:
  - bot opens PR with new Yardstick version
  - CI verifies compatibility
  - reviewer approves and merges
- Do not auto-merge dependency bumps that affect CI policy without a passing pipeline.

## Slash Command Spec

- Repository-specific push/promotion command behavior is documented in `push_code.md`.
- That file governs branch creation, PR handling, merge gates, and release steps.
- Keep it aligned with current GitHub workflow, branch policy, and release process.

## When Adding Checks

1. Implement `Check` in `internal/checks/<name>.go`.
2. Register it in `internal/checks/registry.go`.
3. Add focused tests in `internal/checks/<name>_test.go`.
4. Keep findings actionable and low-noise; optimize for CI signal quality.
