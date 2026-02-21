package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestContributingCheck_FindsRootFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "CONTRIBUTING.md"), []byte("how to contribute"), 0o644); err != nil {
		t.Fatalf("write CONTRIBUTING.md: %v", err)
	}

	fs, err := (ContributingCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}

func TestContributingCheck_MissingWarns(t *testing.T) {
	dir := t.TempDir()
	fs, err := (ContributingCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 1 || fs[0].Level != LevelWarn || fs[0].Check != "contributing" {
		t.Fatalf("unexpected findings: %+v", fs)
	}
}
