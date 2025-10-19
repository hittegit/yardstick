package checks

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestManifestCheck_DetectsGoMod(t *testing.T) {
    dir := t.TempDir()
    if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example"), 0o644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }

    fs, err := (ManifestCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 {
        t.Fatalf("expected 1 finding, got %d", len(fs))
    }
    f := fs[0]
    if f.Check != "manifest" || f.Level != LevelInfo {
        t.Fatalf("unexpected finding: %+v", f)
    }
    if f.Path != filepath.Join(dir, "go.mod") {
        t.Fatalf("unexpected path: %s", f.Path)
    }
    if f.Message == "" || f.Message != "Go project detected via go.mod" {
        t.Fatalf("unexpected message: %q", f.Message)
    }
}

func TestManifestCheck_NoManifestsWarns(t *testing.T) {
    dir := t.TempDir()
    fs, err := (ManifestCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 {
        t.Fatalf("expected 1 finding, got %d", len(fs))
    }
    f := fs[0]
    if f.Level != LevelWarn {
        t.Fatalf("expected warn, got %s", f.Level)
    }
    if f.Path != dir {
        t.Fatalf("unexpected path: %s", f.Path)
    }
    if f.Message == "" {
        t.Fatalf("expected non-empty message")
    }
}

