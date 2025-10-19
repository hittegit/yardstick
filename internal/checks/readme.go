package checks

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// ReadmeCheck verifies that a project contains a README.md file
// and that it includes key sections for documentation consistency.
//
// This check encourages projects to maintain minimal structure for clarity
// when others view the repository, ensuring discoverability and maintainability.
type ReadmeCheck struct{}

// Key returns the unique identifier for this check.
func (ReadmeCheck) Key() string { return "readme" }

// Description provides a one-line explanation of what this check does.
func (ReadmeCheck) Description() string {
	return "Ensures README.md exists and includes required sections"
}

// Run executes the README validation logic.
//
// Behavior:
//   - If README.md is missing, it either reports a warning or scaffolds one if --fix was used.
//   - If present, it verifies required sections exist.
//   - Missing sections are reported as warnings.
func (ReadmeCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	path := filepath.Join(root, "README.md")

	b, err := os.ReadFile(path)
	if err != nil {
		// File missing: attempt autofix or report warning.
		if opts.AutoFix {
			_ = os.WriteFile(path, []byte(defaultReadme()), 0o644)
			return []Finding{{
				Check:   "readme",
				Level:   LevelWarn,
				Path:    path,
				Message: "README was missing, created default template",
				Fixed:   true,
			}}, nil
		}
		return []Finding{{
			Check:   "readme",
			Level:   LevelWarn,
			Path:    path,
			Message: "README.md missing",
		}}, nil
	}

	// Check for required section headers.
	content := string(b)
	required := []string{"## Overview", "## Installation", "## Usage", "## CI", "## License"}
	var findings []Finding

	for _, section := range required {
		if !strings.Contains(content, section) {
			findings = append(findings, Finding{
				Check:   "readme",
				Level:   LevelWarn,
				Path:    path,
				Message: "Missing section: " + section,
			})
		}
	}

	return findings, nil
}

// defaultReadme provides a minimal README.md template scaffolded when missing.
// This template includes all the key sections required by the ReadmeCheck.
func defaultReadme() string {
    // Avoid Markdown code fences in a Go raw string to prevent parser issues.
    // Indent shell examples instead of using ``` blocks.
    return "# Project\n\n" +
        "## Overview\n" +
        "Short summary of the project purpose.\n\n" +
        "## Installation\n" +
        "    go install ./...\n\n" +
        "## Usage\n" +
        "    project --help\n\n" +
        "## CI\n" +
        "Describe how the CI/CD pipeline works and where reports are found.\n\n" +
        "## License\n" +
        "MIT\n"
}
