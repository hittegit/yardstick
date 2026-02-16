# Changelog

All notable changes to this project will be documented here.

## Unreleased
- Initial test suite and CI
- Linting and self-check in CI
- Release workflow with multi-OS builds
- CI lint config now keeps the main linter set enabled while disabling only revive's exported rule to reduce noisy exported identifier warnings
- Expanded .gitignore defaults for Go build artifacts, caches, coverage output, and local yardstick binaries to prevent stray generated files from entering commits
