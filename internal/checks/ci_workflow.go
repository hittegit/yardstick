package checks

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// CIWorkflowCheck ensures the repository has at least one workflow in .github/workflows.
type CIWorkflowCheck struct{}

func (CIWorkflowCheck) Key() string { return "ci_workflow" }

func (CIWorkflowCheck) Description() string {
	return "Ensures at least one workflow file exists in .github/workflows"
}

func (CIWorkflowCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	workflowsDir := filepath.Join(root, ".github", "workflows")
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		return []Finding{{
			Check:   "ci_workflow",
			Level:   LevelWarn,
			Path:    workflowsDir,
			Message: ".github/workflows missing. Add a CI workflow for tests and quality checks",
		}}, nil
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.ToLower(e.Name())
		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			return nil, nil
		}
	}

	return []Finding{{
		Check:   "ci_workflow",
		Level:   LevelWarn,
		Path:    workflowsDir,
		Message: "No workflow files found in .github/workflows. Add at least one .yml or .yaml file",
	}}, nil
}
