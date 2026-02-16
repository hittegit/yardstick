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
//   - If README.md is missing, a warning with guidance is reported.
//   - If present, it verifies required sections exist.
//   - Missing sections are reported as warnings.
func (ReadmeCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	path := filepath.Join(root, "README.md")

	// #nosec G304 -- root/path are intentionally user-selected scan targets.
	b, err := os.ReadFile(path)
	if err != nil {
		// Read-only policy: never write, provide guidance only.
		return []Finding{{
			Check:   "readme",
			Level:   LevelWarn,
			Path:    path,
			Message: "README.md missing. Create README.md with sections: Overview, Installation, Usage, CI, License",
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
