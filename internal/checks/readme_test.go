package checks

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestReadmeCheck_Missing_ReadOnly(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "README.md")
    fs, err := (ReadmeCheck{}).Run(context.Background(), dir, Options{AutoFix: false})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || fs[0].Level != LevelWarn {
        t.Fatalf("expected one warn finding, got %+v", fs)
    }
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        t.Fatalf("README should not be created")
    }
}

func TestReadmeCheck_Missing_NoWriteEvenWithFix(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "README.md")
    fs, err := (ReadmeCheck{}).Run(context.Background(), dir, Options{AutoFix: true})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || fs[0].Fixed {
        t.Fatalf("expected one warn finding with Fixed=false, got %+v", fs)
    }
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        t.Fatalf("README should not be created when AutoFix is true (read-only policy)")
    }
}

func TestReadmeCheck_MissingSections(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "README.md")
    // Provide only two sections so we get warnings for the rest
    content := "# Title\n\n## Overview\ntext\n\n## Usage\ntext\n"
    if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
        t.Fatalf("write README: %v", err)
    }
    fs, err := (ReadmeCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    // Required: Overview, Installation, Usage, CI, License -> we provided 2
    if len(fs) != 3 {
        t.Fatalf("expected 3 warnings for missing sections, got %d", len(fs))
    }
}
