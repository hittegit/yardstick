package checks

import (
    "context"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestReadmeCheck_Missing_NoFix(t *testing.T) {
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
        t.Fatalf("README should not be created without fix")
    }
}

func TestReadmeCheck_Missing_WithFix(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "README.md")
    fs, err := (ReadmeCheck{}).Run(context.Background(), dir, Options{AutoFix: true})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 1 || !fs[0].Fixed {
        t.Fatalf("expected one fixed finding, got %+v", fs)
    }
    b, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("read scaffolded README: %v", err)
    }
    if !strings.Contains(string(b), "## Overview") {
        t.Fatalf("scaffolded README missing expected content")
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

