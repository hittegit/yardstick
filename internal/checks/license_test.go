package checks

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestLicenseCheck_Missing_NoFix(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "LICENSE")
    fs, err := (LicenseCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || fs[0].Level != LevelWarn {
        t.Fatalf("expected warn for missing LICENSE, got %+v", fs)
    }
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        t.Fatalf("LICENSE should not be created without fix")
    }
}

func TestLicenseCheck_Missing_WithFix(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "LICENSE")
    fs, err := (LicenseCheck{}).Run(context.Background(), dir, Options{AutoFix: true})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || !fs[0].Fixed {
        t.Fatalf("expected fixed finding, got %+v", fs)
    }
    if _, err := os.Stat(path); err != nil {
        t.Fatalf("expected LICENSE to be created: %v", err)
    }
}

