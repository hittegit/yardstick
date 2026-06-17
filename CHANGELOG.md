# Changelog

All notable changes to this project will be documented here.

## v0.5.0 - 2026-06-17

- Added ecosystem-specific checks:
  - `javascript_framework` for broad JavaScript framework conventions with explicit Next.js compatibility validation
  - `python_project` for Python project conventions, test-layout signals, and tooling guidance
- Added tests for new JavaScript framework and Python checks
- Updated check registry, guidance, and README documentation for expanded ecosystem support
- Expanded reporting output:
  - Added top-level summary text and per-check pass/fail status in JSON output
  - Added verbose table output showing all executed checks, status, and finding counts
  - Added remediation guidance and rationale for failed checks
- Added check guidance catalog used by reporting to explain why failures matter and how to resolve them
- Added tests covering per-check status summaries and verbose table output
- Migrated `.golangci.yml` to golangci-lint v2 config format
- Bumped Go from 1.22 to 1.26 across go.mod, CI workflows, and the CI Docker image
- Bumped `actions/checkout` to v6 and `golangci-lint-action` to v9 (golangci-lint binary v2.12.2) in CI and release workflows
- Added `renovate.json` for automated dependency update PRs
- Added `#nosec G703` annotations in `readme_links.go` for path-traversal false positives on user-provided repo roots

## v0.2.0 - 2026-02-21

- Added new CI hygiene checks:
  - `codeowners`
  - `security_policy`
  - `contributing`
  - `ci_workflow`
  - `readme_links`
- Updated check registry and README documentation to include the new checks
- Hardened CI contract for downstream consumers:
  - invalid `-format` values now fail fast
  - unknown `-only` check keys now fail fast
  - `-only` parsing now trims whitespace
  - JSON findings now always include `path` and `fixed`
- Added `AGENT.md` with CI contract guardrails and maintenance guidance
- Added README section with pinned-version GitHub Actions usage for external CI

## v0.1.0 - 2026-02-21

- Initial test suite and CI
- Linting and self-check in CI
- Release workflow with multi-OS builds
- CI lint config now keeps the main linter set enabled while disabling only revive's exported rule to reduce noisy exported identifier warnings
- Expanded .gitignore defaults for Go build artifacts, caches, coverage output, and local yardstick binaries to prevent stray generated files from entering commits
