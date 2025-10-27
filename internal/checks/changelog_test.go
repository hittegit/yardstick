package checks

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestChangelogCheck_Missing_ReadOnly(t *testing.T) {
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
        t.Fatalf("CHANGELOG.md should not be created")
    }
}

func TestChangelogCheck_Missing_NoWriteEvenWithFix(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "CHANGELOG.md")
    fs, err := (ChangelogCheck{}).Run(context.Background(), dir, Options{AutoFix: true})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || fs[0].Fixed {
        t.Fatalf("expected warn finding with Fixed=false, got %+v", fs)
    }
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        t.Fatalf("CHANGELOG.md should not be created when AutoFix is true (read-only policy)")
    }
}
