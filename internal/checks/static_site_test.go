package checks

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestStaticSiteCheck_NoConfigNoop(t *testing.T) {
    dir := t.TempDir()
    fs, err := (StaticSiteCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 0 {
        t.Fatalf("expected no findings when no static site config, got %d", len(fs))
    }
}

func TestStaticSiteCheck_Jekyll_MissingPieces(t *testing.T) {
    dir := t.TempDir()
    if err := os.WriteFile(filepath.Join(dir, "_config.yml"), []byte("title: site"), 0o644); err != nil {
        t.Fatalf("write _config.yml: %v", err)
    }
    fs, err := (StaticSiteCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) < 2 { // index.md, pages/, assets/
        t.Fatalf("expected warnings for missing index/pages/assets, got %d", len(fs))
    }
}

func TestStaticSiteCheck_Jekyll_MinimalOK(t *testing.T) {
    dir := t.TempDir()
    if err := os.WriteFile(filepath.Join(dir, "_config.yml"), []byte("title: site"), 0o644); err != nil {
        t.Fatalf("write _config.yml: %v", err)
    }
    if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Home"), 0o644); err != nil {
        t.Fatalf("write index.md: %v", err)
    }
    if err := os.Mkdir(filepath.Join(dir, "pages"), 0o755); err != nil {
        t.Fatalf("mkdir pages: %v", err)
    }
    if err := os.WriteFile(filepath.Join(dir, "pages", "about.md"), []byte("# About"), 0o644); err != nil {
        t.Fatalf("write about.md: %v", err)
    }
    if err := os.Mkdir(filepath.Join(dir, "assets"), 0o755); err != nil {
        t.Fatalf("mkdir assets: %v", err)
    }
    fs, err := (StaticSiteCheck{}).Run(context.Background(), dir, Options{})
    if err != nil {
        t.Fatalf("run: %v", err)
    }
    if len(fs) != 0 {
        t.Fatalf("expected no warnings for minimal structure, got %d", len(fs))
    }
}

