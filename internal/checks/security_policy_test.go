package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSecurityPolicyCheck_FindsRootFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "SECURITY.md"), []byte("report here"), 0o644); err != nil {
		t.Fatalf("write SECURITY.md: %v", err)
	}

	fs, err := (SecurityPolicyCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}

func TestSecurityPolicyCheck_MissingWarns(t *testing.T) {
	dir := t.TempDir()
	fs, err := (SecurityPolicyCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 1 || fs[0].Level != LevelWarn || fs[0].Check != "security_policy" {
		t.Fatalf("unexpected findings: %+v", fs)
	}
}
