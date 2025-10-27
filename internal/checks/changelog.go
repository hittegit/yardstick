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
//   - If missing, a warning is reported with guidance.
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

    return []Finding{{
        Check:   "changelog",
        Level:   LevelWarn,
        Path:    path,
        Message: "CHANGELOG.md missing. Add a changelog documenting notable changes (Keep a Changelog format recommended)",
    }}, nil
}
