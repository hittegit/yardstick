package checks

import (
	"context"
	"os"
	"path/filepath"
)

// ManifestCheck attempts to detect a project's primary ecosystem by looking
// for common manifest files. It is intentionally neutral so yardstick can be
// used across Go, Node, Python, Ruby, Rust, and static site projects.
//
// Behavior
//   - If a known manifest is found, emit an info finding stating which
//     ecosystem was detected and which file triggered it.
//   - If no known manifests are found, emit a warn finding suggesting the
//     user add an appropriate manifest for their stack.
//   - This check does not auto-create manifests, since that choice is
//     project specific and harder to do safely.
//
// Extending detection
//   - Add new entries to the candidates slice with the filename and a short
//     label. Keep the check simple and fast.
//   - If needed later, we can add per-ecosystem subchecks, for example
//     NodeLockfileCheck, PythonVenvCheck, etc.
//
// Examples of files detected
//   - Go:         go.mod
//   - Node:       package.json
//   - Python:     pyproject.toml or requirements.txt
//   - Ruby:       Gemfile
//   - Rust:       Cargo.toml
//   - PHP:        composer.json
//   - Static:     _config.yml, .eleventy.js, mkdocs.yml
//   - Docs only:  README.md without a manifest will still pass other checks
//                 but this one will warn.
type ManifestCheck struct{}

// Key returns the unique identifier for this check.
func (ManifestCheck) Key() string { return "manifest" }

// Description provides a short explanation of what this check validates.
func (ManifestCheck) Description() string { return "Detects project ecosystem by scanning for common manifests" }

// Run performs the manifest detection logic.
func (ManifestCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	// Known manifest candidates and a friendly label for reporting.
	type candidate struct {
		name  string
		label string
	}
	candidates := []candidate{
		{"go.mod", "Go"},
		{"package.json", "Node"},
		{"pyproject.toml", "Python"},
		{"requirements.txt", "Python"},
		{"Gemfile", "Ruby"},
		{"Cargo.toml", "Rust"},
		{"composer.json", "PHP"},
		{"_config.yml", "Static site"},      // Jekyll and similar
		{".eleventy.js", "Static site"},      // Eleventy
		{"mkdocs.yml", "Static site"},        // MkDocs
	}

	for _, c := range candidates {
		p := filepath.Join(root, c.name)
		if _, err := os.Stat(p); err == nil {
			return []Finding{{
				Check:   "manifest",
				Level:   LevelInfo,
				Path:    p,
				Message: c.label + " project detected via " + c.name,
			}}, nil
		}
	}

	// Nothing matched. Suggest adding a manifest appropriate to the stack.
	return []Finding{{
		Check:   "manifest",
		Level:   LevelWarn,
		Path:    root,
		Message: "No common project manifest found. Expected one of: go.mod, package.json, pyproject.toml, requirements.txt, Gemfile, Cargo.toml, composer.json, or a static site config",
	}}, nil
}
