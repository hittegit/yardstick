package checks

import (
	"context"
	"os"
	"path/filepath"
)

// ContributingCheck ensures contribution guidelines are present.
type ContributingCheck struct{}

func (ContributingCheck) Key() string { return "contributing" }

func (ContributingCheck) Description() string {
	return "Ensures CONTRIBUTING.md exists in a standard GitHub location"
}

func (ContributingCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	candidates := []string{
		filepath.Join(root, "CONTRIBUTING.md"),
		filepath.Join(root, ".github", "CONTRIBUTING.md"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return nil, nil
		}
	}

	return []Finding{{
		Check:   "contributing",
		Level:   LevelWarn,
		Path:    filepath.Join(root, "CONTRIBUTING.md"),
		Message: "CONTRIBUTING.md missing. Add contributor guidelines in CONTRIBUTING.md or .github/CONTRIBUTING.md",
	}}, nil
}
