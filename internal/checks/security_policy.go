package checks

import (
	"context"
	"os"
	"path/filepath"
)

// SecurityPolicyCheck ensures a repository provides a vulnerability disclosure policy.
type SecurityPolicyCheck struct{}

func (SecurityPolicyCheck) Key() string { return "security_policy" }

func (SecurityPolicyCheck) Description() string {
	return "Ensures SECURITY.md exists in a standard GitHub location"
}

func (SecurityPolicyCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	candidates := []string{
		filepath.Join(root, "SECURITY.md"),
		filepath.Join(root, ".github", "SECURITY.md"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return nil, nil
		}
	}

	return []Finding{{
		Check:   "security_policy",
		Level:   LevelWarn,
		Path:    filepath.Join(root, "SECURITY.md"),
		Message: "SECURITY.md missing. Add vulnerability reporting guidance in SECURITY.md or .github/SECURITY.md",
	}}, nil
}
