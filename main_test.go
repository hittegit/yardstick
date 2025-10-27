package main

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

// Test run() happy path limited to manifest, to avoid unrelated warnings.
func TestRun_OnlyManifest_NoStrict(t *testing.T) {
    dir := t.TempDir()
    // Create a minimal go.mod so ManifestCheck emits info.
    if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example"), 0o644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }
    *flagPath = dir
    *flagOnly = "manifest"
    *flagStrict = false
    *flagFormat = "json"
    if err := run(context.Background()); err != nil {
        t.Fatalf("run returned error: %v", err)
    }
}

// Test strict mode fails on warnings: use readme check only in an empty dir.
func TestRun_StrictFailsOnWarn(t *testing.T) {
    dir := t.TempDir()
    *flagPath = dir
    *flagOnly = "readme"
    *flagStrict = true
    *flagFormat = "table"
    if err := run(context.Background()); err == nil {
        t.Fatalf("expected error due to warnings in strict mode, got nil")
    }
}
