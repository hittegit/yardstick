package checks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

// JavaScriptFrameworkCheck validates baseline structure for JavaScript framework projects.
// It applies broad checks for framework projects and adds Next.js-specific compatibility checks.
type JavaScriptFrameworkCheck struct{}

func (JavaScriptFrameworkCheck) Key() string { return "javascript_framework" }

func (JavaScriptFrameworkCheck) Description() string {
	return "Validates baseline conventions for JavaScript framework projects, including Next.js compatibility"
}

type packageJSON struct {
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (JavaScriptFrameworkCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	pkgPath := filepath.Join(root, "package.json")
	// Not a JS project, no-op.
	if _, err := os.Stat(pkgPath); err != nil {
		return nil, nil
	}

	// #nosec G304 -- root/path are intentionally user-selected scan targets.
	b, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil, err
	}

	var pkg packageJSON
	if err := json.Unmarshal(b, &pkg); err != nil {
		return []Finding{{
			Check:   "javascript_framework",
			Level:   LevelWarn,
			Path:    pkgPath,
			Message: "package.json is not valid JSON. Fix JSON syntax so framework checks can run reliably",
		}}, nil
	}

	frameworks := detectJavaScriptFrameworks(pkg)
	// Node project without a recognized framework, skip framework-specific checks.
	if len(frameworks) == 0 {
		return nil, nil
	}

	var findings []Finding
	if !hasScript(pkg.Scripts, "dev") {
		findings = append(findings, Finding{
			Check:   "javascript_framework",
			Level:   LevelWarn,
			Path:    pkgPath,
			Message: "Missing scripts.dev in package.json. Add a dev script for local development",
		})
	}
	if !hasScript(pkg.Scripts, "build") {
		findings = append(findings, Finding{
			Check:   "javascript_framework",
			Level:   LevelWarn,
			Path:    pkgPath,
			Message: "Missing scripts.build in package.json. Add a build script for CI and production builds",
		})
	}

	// Explicit Next.js compatibility checks.
	if slices.Contains(frameworks, "next") {
		if !hasScript(pkg.Scripts, "start") {
			findings = append(findings, Finding{
				Check:   "javascript_framework",
				Level:   LevelWarn,
				Path:    pkgPath,
				Message: "Next.js project missing scripts.start. Add a start script for runtime compatibility",
			})
		}
		if !hasDir(filepath.Join(root, "app")) && !hasDir(filepath.Join(root, "pages")) {
			findings = append(findings, Finding{
				Check:   "javascript_framework",
				Level:   LevelWarn,
				Path:    root,
				Message: "Next.js project missing both app/ and pages/. Add at least one routing directory",
			})
		}
	}

	return findings, nil
}

func detectJavaScriptFrameworks(pkg packageJSON) []string {
	frameworkDeps := []string{
		"next",
		"react-scripts",
		"vite",
		"nuxt",
		"@angular/core",
		"@sveltejs/kit",
		"gatsby",
		"@remix-run/react",
	}
	var frameworks []string
	for _, dep := range frameworkDeps {
		if hasDep(pkg.Dependencies, dep) || hasDep(pkg.DevDependencies, dep) {
			frameworks = append(frameworks, dep)
		}
	}
	return frameworks
}

func hasScript(scripts map[string]string, key string) bool {
	if scripts == nil {
		return false
	}
	_, ok := scripts[key]
	return ok
}

func hasDep(deps map[string]string, name string) bool {
	if deps == nil {
		return false
	}
	_, ok := deps[name]
	return ok
}

func hasDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}
