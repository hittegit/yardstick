package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func snapshotFlags() func() {
	format := *flagFormat
	path := *flagPath
	strict := *flagStrict
	only := *flagOnly
	list := *flagList
	version := *flagVersion
	return func() {
		*flagFormat = format
		*flagPath = path
		*flagStrict = strict
		*flagOnly = only
		*flagList = list
		*flagVersion = version
	}
}

// Test run() happy path limited to manifest, to avoid unrelated warnings.
func TestRun_OnlyManifest_NoStrict(t *testing.T) {
	t.Cleanup(snapshotFlags())
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
	t.Cleanup(snapshotFlags())
	dir := t.TempDir()
	*flagPath = dir
	*flagOnly = "readme"
	*flagStrict = true
	*flagFormat = "table"
	if err := run(context.Background()); err == nil {
		t.Fatalf("expected error due to warnings in strict mode, got nil")
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	t.Cleanup(snapshotFlags())
	*flagFormat = "yaml"
	*flagOnly = "manifest"
	if err := run(context.Background()); err == nil {
		t.Fatalf("expected invalid format error, got nil")
	}
}

func TestRun_OnlyRejectsUnknownCheck(t *testing.T) {
	t.Cleanup(snapshotFlags())
	*flagFormat = "table"
	*flagOnly = "manifest,does-not-exist"
	err := run(context.Background())
	if err == nil {
		t.Fatalf("expected error for unknown check key, got nil")
	}
	if !strings.Contains(err.Error(), "does-not-exist") {
		t.Fatalf("error missing unknown key, got: %v", err)
	}
}

func TestSplitCSV_Trimmed(t *testing.T) {
	got := splitCSV("manifest, readme ,  license,,")
	want := []string{"manifest", "readme", "license"}
	if len(got) != len(want) {
		t.Fatalf("unexpected length, got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected parse, got=%v want=%v", got, want)
		}
	}
}
