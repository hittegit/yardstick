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
//   - If missing, a warning is reported.
//   - With --fix, a default .gitignore is created.
type GitIgnoreCheck struct{}

// Key returns the unique identifier for this check.
func (GitIgnoreCheck) Key() string { return "gitignore" }

// Description provides a short explanation of what this check validates.
func (GitIgnoreCheck) Description() string {
	return ".gitignore includes common entries or is scaffolded if missing"
}

// Run executes the .gitignore validation and optional remediation.
func (GitIgnoreCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	path := filepath.Join(root, ".gitignore")

	// If .gitignore already exists, nothing to report.
	if _, err := os.Stat(path); err == nil {
		return nil, nil
	}

	// Missing file. Either create a default or warn.
	if opts.AutoFix {
		_ = os.WriteFile(path, []byte(defaultIgnore), 0o644)
		return []Finding{{
			Check:   "gitignore",
			Level:   LevelWarn,
			Path:    path,
			Message: ".gitignore missing, created default entries",
			Fixed:   true,
		}}, nil
	}

	return []Finding{{
		Check:   "gitignore",
		Level:   LevelWarn,
		Path:    path,
		Message: ".gitignore missing",
	}}, nil
}

// defaultIgnore provides a minimal, sensible default .gitignore.
// Tailor this over time to your ecosystem.
const defaultIgnore = `# Build artifacts
bin/
dist/
*.out

# Go
*.test
coverage.txt

# Editors
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db
`
