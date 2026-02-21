# Changelog

All notable changes to this project will be documented here.

## Unreleased
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
