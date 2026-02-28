package checks

import (
	"context"
	"os"
	"path/filepath"
)

// PythonProjectCheck validates baseline conventions for Python projects.
type PythonProjectCheck struct{}

func (PythonProjectCheck) Key() string { return "python_project" }

func (PythonProjectCheck) Description() string {
	return "Validates baseline conventions for Python projects"
}

func (PythonProjectCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	pyprojectPath := filepath.Join(root, "pyproject.toml")
	requirementsPath := filepath.Join(root, "requirements.txt")

	hasPyproject := fileExists(pyprojectPath)
	hasRequirements := fileExists(requirementsPath)
	// Not a Python project we recognize, no-op.
	if !hasPyproject && !hasRequirements {
		return nil, nil
	}

	var findings []Finding

	// Modern packaging/tooling guidance when requirements.txt is the only manifest.
	if hasRequirements && !hasPyproject {
		findings = append(findings, Finding{
			Check:   "python_project",
			Level:   LevelWarn,
			Path:    requirementsPath,
			Message: "requirements.txt found without pyproject.toml. Add pyproject.toml for modern tooling and metadata interoperability",
		})
	}

	if !hasPythonTestSignal(root) {
		findings = append(findings, Finding{
			Check:   "python_project",
			Level:   LevelWarn,
			Path:    root,
			Message: "No Python test layout/config detected. Add tests/ or pytest/tox configuration to support CI validation",
		})
	}

	return findings, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasPythonTestSignal(root string) bool {
	testsDir := filepath.Join(root, "tests")
	pytestINI := filepath.Join(root, "pytest.ini")
	toxINI := filepath.Join(root, "tox.ini")
	noxfile := filepath.Join(root, "noxfile.py")
	setupCFG := filepath.Join(root, "setup.cfg")

	if st, err := os.Stat(testsDir); err == nil && st.IsDir() {
		return true
	}
	return fileExists(pytestINI) || fileExists(toxINI) || fileExists(noxfile) || fileExists(setupCFG)
}
