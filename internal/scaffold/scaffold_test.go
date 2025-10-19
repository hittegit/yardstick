package scaffold

import (
    "os"
    "path/filepath"
    "testing"
)

func TestRun_CreatesMissingFiles_DoesNotOverwrite(t *testing.T) {
    dir := t.TempDir()

    // Pre-create one file to ensure it is not overwritten
    pre := filepath.Join(dir, "README.md")
    if err := os.WriteFile(pre, []byte("CUSTOM"), 0o644); err != nil {
        t.Fatalf("write pre-existing: %v", err)
    }

    if err := Run(dir); err != nil {
        t.Fatalf("Run: %v", err)
    }

    // Ensure pre-existing file preserved
    b, err := os.ReadFile(pre)
    if err != nil {
        t.Fatalf("read pre-existing: %v", err)
    }
    if string(b) != "CUSTOM" {
        t.Fatalf("expected README.md to be unchanged")
    }

    // Ensure other expected files are created
    expected := []string{".prettierrc", ".markdownlint.json", ".editorconfig"}
    for _, name := range expected {
        p := filepath.Join(dir, name)
        if _, err := os.Stat(p); err != nil {
            t.Fatalf("expected %s to exist: %v", name, err)
        }
        b, err := os.ReadFile(p)
        if err != nil || len(b) == 0 {
            t.Fatalf("expected %s to have content", name)
        }
    }
}

