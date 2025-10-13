package checks

import (
	"context"
	"os"
	"path/filepath"
)

// ChangelogCheck ensures a project includes a CHANGELOG.md file.
// A changelog helps consumers understand what changed across versions.
//
// Behavior
//   - If CHANGELOG.md exists, no findings are emitted.
//   - If missing, a warning is reported.
//   - With --fix, a minimal changelog is scaffolded.
type ChangelogCheck struct{}

// Key returns the unique identifier for this check.
func (ChangelogCheck) Key() string { return "changelog" }

// Description provides a one-line explanation of this check.
func (ChangelogCheck) Description() string { return "Ensures CHANGELOG.md exists" }

// Run performs the changelog validation and optional remediation.
func (ChangelogCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	path := filepath.Join(root, "CHANGELOG.md")

	// If file exists, nothing to report.
	if _, err := os.Stat(path); err == nil {
		return nil, nil
	}

	// Missing file. Either create a default or warn.
	if opts.AutoFix {
		_ = os.WriteFile(path, []byte(defaultChangelog), 0o644)
		return []Finding{{
			Check:   "changelog",
			Level:   LevelWarn,
			Path:    path,
			Message: "CHANGELOG.md missing, created a starter file",
			Fixed:   true,
		}}, nil
	}

	return []Finding{{
		Check:   "changelog",
		Level:   LevelWarn,
		Path:    path,
		Message: "CHANGELOG.md missing",
	}}, nil
}

// defaultChangelog is a minimal template used for scaffolding.
const defaultChangelog = `# Changelog

All notable changes to this project will be documented in this file.

## Unreleased
- Initial scaffolding
`
