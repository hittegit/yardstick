package checks

import (
    "context"
    "os"
    "path/filepath"
)

// LicenseCheck ensures that a LICENSE file is present in the repository.
//
// Licensing clarity is essential for both internal and external code sharing.
// This check verifies that a LICENSE file exists. Yardstick is read-only and
// does not create files; it provides guidance when the license is missing.
type LicenseCheck struct{}

// Key returns the unique identifier for this check.
func (LicenseCheck) Key() string { return "license" }

// Description provides a short explanation of what this check validates.
func (LicenseCheck) Description() string {
	return "Ensures LICENSE file is present"
}

// Run executes the license file validation.
//
// Behavior:
//   - If LICENSE exists, no findings are produced.
//   - If missing, a warning is emitted.
//   - If --fix is used, a default MIT license is created.
func (LicenseCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	path := filepath.Join(root, "LICENSE")

	// If file already exists, nothing to report.
	if _, err := os.Stat(path); err == nil {
		return nil, nil
	}

    return []Finding{{
        Check:   "license",
        Level:   LevelWarn,
        Path:    path,
        Message: "LICENSE missing. Add a license file (e.g., MIT, Apache-2.0) appropriate to your project",
    }}, nil
}
