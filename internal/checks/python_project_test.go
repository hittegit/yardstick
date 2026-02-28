package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPythonProjectCheck_NoPythonSignals_NoOp(t *testing.T) {
	dir := t.TempDir()
	fs, err := (PythonProjectCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}

func TestPythonProjectCheck_PyprojectWithTests_NoFindings(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("[project]\nname='demo'\n"), 0o644); err != nil {
		t.Fatalf("write pyproject.toml: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, "tests"), 0o750); err != nil {
		t.Fatalf("mkdir tests: %v", err)
	}

	fs, err := (PythonProjectCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}

func TestPythonProjectCheck_RequirementsOnlyWarns(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte("pytest==8.3.0\n"), 0o644); err != nil {
		t.Fatalf("write requirements.txt: %v", err)
	}

	fs, err := (PythonProjectCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 2 {
		t.Fatalf("expected 2 findings, got %d (%+v)", len(fs), fs)
	}
	for _, f := range fs {
		if f.Check != "python_project" || f.Level != LevelWarn {
			t.Fatalf("unexpected finding: %+v", f)
		}
	}
}

func TestPythonProjectCheck_PyprojectMissingTestSignal_Warns(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("[project]\nname='demo'\n"), 0o644); err != nil {
		t.Fatalf("write pyproject.toml: %v", err)
	}

	fs, err := (PythonProjectCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 1 {
		t.Fatalf("expected 1 finding, got %d (%+v)", len(fs), fs)
	}
}
