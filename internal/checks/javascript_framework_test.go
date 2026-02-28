package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestJavaScriptFrameworkCheck_NoPackageJSON_NoOp(t *testing.T) {
	dir := t.TempDir()
	fs, err := (JavaScriptFrameworkCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}

func TestJavaScriptFrameworkCheck_NextJSValid_NoFindings(t *testing.T) {
	dir := t.TempDir()
	pkg := `{
  "dependencies": {"next": "^14.0.0"},
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, "app"), 0o750); err != nil {
		t.Fatalf("mkdir app: %v", err)
	}

	fs, err := (JavaScriptFrameworkCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}

func TestJavaScriptFrameworkCheck_NextJSMissingConventions_Warns(t *testing.T) {
	dir := t.TempDir()
	pkg := `{
  "dependencies": {"next": "^14.0.0"},
  "scripts": {
    "dev": "next dev"
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}

	fs, err := (JavaScriptFrameworkCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 3 {
		t.Fatalf("expected 3 findings, got %d (%+v)", len(fs), fs)
	}
	for _, f := range fs {
		if f.Check != "javascript_framework" || f.Level != LevelWarn {
			t.Fatalf("unexpected finding: %+v", f)
		}
	}
}

func TestJavaScriptFrameworkCheck_GeneralFrameworkMissingBuild_Warns(t *testing.T) {
	dir := t.TempDir()
	pkg := `{
  "devDependencies": {"vite": "^5.0.0"},
  "scripts": {
    "dev": "vite"
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}

	fs, err := (JavaScriptFrameworkCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 1 {
		t.Fatalf("expected 1 finding, got %d (%+v)", len(fs), fs)
	}
	if fs[0].Message == "" {
		t.Fatalf("expected non-empty message")
	}
}
