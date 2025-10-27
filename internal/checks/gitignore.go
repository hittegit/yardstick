package checks

import (
    "context"
    "os"
    "path/filepath"
)

// GitIgnoreCheck ensures a repository has a .gitignore file with
// a reasonable starter set of entries. This avoids committing build
// artifacts, editor settings, and OS files.
//
// Behavior
//   - If .gitignore exists, no findings are emitted.
//   - If missing, a warning is reported with guidance.
type GitIgnoreCheck struct{}

// Key returns the unique identifier for this check.
func (GitIgnoreCheck) Key() string { return "gitignore" }

// Description provides a short explanation of what this check validates.
func (GitIgnoreCheck) Description() string { return ".gitignore includes common entries or guidance is provided if missing" }

// Run executes the .gitignore validation.
func (GitIgnoreCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	path := filepath.Join(root, ".gitignore")

	// If .gitignore already exists, nothing to report.
	if _, err := os.Stat(path); err == nil {
		return nil, nil
	}

    return []Finding{{
        Check:   "gitignore",
        Level:   LevelWarn,
        Path:    path,
        Message: ".gitignore missing. Add common ignores for your ecosystem (build artifacts, editor files, OS files)",
    }}, nil
}
