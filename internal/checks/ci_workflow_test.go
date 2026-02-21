package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCIWorkflowCheck_MissingWarns(t *testing.T) {
	dir := t.TempDir()
	fs, err := (CIWorkflowCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 1 || fs[0].Level != LevelWarn || fs[0].Check != "ci_workflow" {
		t.Fatalf("unexpected findings: %+v", fs)
	}
}

func TestCIWorkflowCheck_FindsWorkflowFile(t *testing.T) {
	dir := t.TempDir()
	workflows := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(workflows, 0o755); err != nil {
		t.Fatalf("mkdir workflows: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workflows, "ci.yml"), []byte("name: ci"), 0o644); err != nil {
		t.Fatalf("write workflow: %v", err)
	}

	fs, err := (CIWorkflowCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}
