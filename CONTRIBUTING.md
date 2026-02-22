# Contributing

Thanks for your interest in contributing to Yardstick!

## Development setup

- Go 1.22+
- `make` (optional convenience)

## Common tasks

```bash
# Run linters and tests
make lint test

# Or without Makefile
golangci-lint run
go test ./...
```

## Branching and PRs

- Use feature branches, submit small PRs with a clear scope.
- Add/adjust tests for behavior changes.
- Keep code readable and well-commented.

## Commit messages

- Use conventional prefixes when possible (feat, fix, chore, ci, docs, refactor, test).
- Keep subject lines short; add details in the body if helpful.

## Code style

- Go: defer to `gofmt`/`go vet` and linters; tabs in `.go` files.
- YAML/JSON: 2-space indentation.
- Markdown: trailing spaces allowed for intentional line breaks.

## Security

- Do not include secrets in code or tests.
- See SECURITY.md for reporting vulnerabilities.

## Releases

- Tag `vX.Y.Z` to trigger the release workflow.
- Release artifacts are attached to the GitHub Release.
