package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestReadmeLinksCheck_NoReadmeNoop(t *testing.T) {
	dir := t.TempDir()
	fs, err := (ReadmeLinksCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}

func TestReadmeLinksCheck_BrokenSelfAnchorWarns(t *testing.T) {
	dir := t.TempDir()
	readme := "# Title\n\nSee [missing](#does-not-exist).\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	fs, err := (ReadmeLinksCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 1 || fs[0].Check != "readme_links" || fs[0].Level != LevelWarn {
		t.Fatalf("unexpected findings: %+v", fs)
	}
}

func TestReadmeLinksCheck_MissingLocalFileWarns(t *testing.T) {
	dir := t.TempDir()
	readme := "# Title\n\nSee [guide](docs/guide.md).\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	fs, err := (ReadmeLinksCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 1 || fs[0].Check != "readme_links" || fs[0].Level != LevelWarn {
		t.Fatalf("unexpected findings: %+v", fs)
	}
}

func TestReadmeLinksCheck_ExistingLinksPass(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "docs"), 0o755); err != nil {
		t.Fatalf("mkdir docs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "docs", "guide.md"), []byte("# Guide\n\n## Install\n"), 0o644); err != nil {
		t.Fatalf("write guide: %v", err)
	}

	readme := "" +
		"# Project\n\n" +
		"## Usage\n\n" +
		"- [Local heading](#usage)\n" +
		"- [Guide](docs/guide.md)\n" +
		"- [Guide section](docs/guide.md#install)\n" +
		"- [Site](https://example.com)\n" +
		"- ![Diagram](docs/diagram.png)\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}

	fs, err := (ReadmeLinksCheck{}).Run(context.Background(), dir, Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(fs) != 0 {
		t.Fatalf("expected no findings, got %+v", fs)
	}
}
