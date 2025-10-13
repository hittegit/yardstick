package scaffold

import (
	"os"
	"path/filepath"
)

// Package scaffold writes a small set of starter files into a repository
// when they are missing. Yardstick only scaffolds files that are safe to
// generate without project-specific choices, keeping behavior predictable.
//
// Files created by default
//   - README.md: a minimal skeleton with common sections
//   - .prettierrc: opinionated but common prose wrapping settings
//   - .markdownlint.json: lint settings tuned for long lines and inline HTML
//   - .editorconfig: consistent line endings, final newline, and Go tab width
//
// The caller controls when this runs via the --scaffold flag. Existing files
// are never overwritten.
func Run(root string) error {
	files := map[string]string{
		"README.md": `# Project

## Overview
Brief summary of what this project does and who it serves.

## Installation
Describe how to install or build the project.

## Usage
Explain common commands or links to deeper docs.

## CI
Describe where CI lives and how to read reports.

## License
MIT
`,
		".prettierrc":      "{\n  \"proseWrap\": \"always\"\n}\n",
		".markdownlint.json": "{\n  \"default\": true,\n  \"MD013\": { \"line_length\": 120 },\n  \"MD033\": false\n}\n",
		".editorconfig": `[root]
root = true

[*]
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.{md,markdown}]
max_line_length = 120
indent_style = space
indent_size = 2

[*.go]
indent_style = tab
max_line_length = 120
`,
	}

	for name, content := range files {
		path := filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			// File already exists, do not overwrite.
			continue
		}
		_ = os.WriteFile(path, []byte(content), 0o644)
	}
	return nil
}
