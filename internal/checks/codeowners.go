package checks

import (
	"context"
	"os"
	"path/filepath"
)

// CodeownersCheck ensures a repository defines ownership rules in CODEOWNERS.
type CodeownersCheck struct{}

func (CodeownersCheck) Key() string { return "codeowners" }

func (CodeownersCheck) Description() string {
	return "Ensures CODEOWNERS exists in a standard GitHub location"
}

func (CodeownersCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	candidates := []string{
		filepath.Join(root, "CODEOWNERS"),
		filepath.Join(root, ".github", "CODEOWNERS"),
		filepath.Join(root, "docs", "CODEOWNERS"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return nil, nil
		}
	}

	return []Finding{{
		Check:   "codeowners",
		Level:   LevelWarn,
		Path:    filepath.Join(root, "CODEOWNERS"),
		Message: "CODEOWNERS missing. Add ownership rules in CODEOWNERS, .github/CODEOWNERS, or docs/CODEOWNERS",
	}}, nil
}
