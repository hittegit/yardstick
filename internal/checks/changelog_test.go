package checks

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestChangelogCheck_Missing_NoFix(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "CHANGELOG.md")
    fs, err := (ChangelogCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || fs[0].Level != LevelWarn {
        t.Fatalf("expected warn for missing CHANGELOG.md, got %+v", fs)
    }
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        t.Fatalf("CHANGELOG.md should not be created without fix")
    }
}

func TestChangelogCheck_Missing_WithFix(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "CHANGELOG.md")
    fs, err := (ChangelogCheck{}).Run(context.Background(), dir, Options{AutoFix: true})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || !fs[0].Fixed {
        t.Fatalf("expected fixed finding, got %+v", fs)
    }
    if _, err := os.Stat(path); err != nil {
        t.Fatalf("expected CHANGELOG.md to be created: %v", err)
    }
}

