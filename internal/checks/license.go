package checks

import (
	"context"
	"os"
	"path/filepath"
)

// LicenseCheck ensures that a LICENSE file is present in the repository.
//
// Licensing clarity is essential for both internal and external code sharing.
// This check verifies that a LICENSE file exists, and if --fix is enabled,
// it scaffolds a permissive MIT license by default. Teams may later replace it
// with their specific license requirements.
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

	// Missing LICENSE: handle based on fix flag.
	if opts.AutoFix {
		_ = os.WriteFile(path, []byte(defaultMIT), 0o644)
		return []Finding{{
			Check:   "license",
			Level:   LevelWarn,
			Path:    path,
			Message: "LICENSE missing, created MIT license",
			Fixed:   true,
		}}, nil
	}

	return []Finding{{
		Check:   "license",
		Level:   LevelWarn,
		Path:    path,
		Message: "LICENSE missing",
	}}, nil
}

// defaultMIT provides a simple permissive MIT license.
// This text is inserted automatically when a license is missing and --fix is used.
// Users should later edit or replace it if their project requires a different license.
const defaultMIT = `MIT License

Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
