package checks

import (
    "context"
    "os"
    "path/filepath"
)

// StaticSiteCheck validates minimal structure for static-site projects.
// Currently focuses on Jekyll-style repos (detected via _config.yml).
// If no static-site config is detected, this check is a no-op.
type StaticSiteCheck struct{}

func (StaticSiteCheck) Key() string { return "static_site" }

func (StaticSiteCheck) Description() string {
    return "Validates minimal structure for static-site projects (e.g., Jekyll)"
}

func (StaticSiteCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
    // Jekyll detection: _config.yml at repo root
    jekyllCfg := filepath.Join(root, "_config.yml")
    if _, err := os.Stat(jekyllCfg); err != nil {
        // Not a static-site project we recognize; no-op
        return nil, nil
    }

    var out []Finding

    // Require a landing page
    idx := filepath.Join(root, "index.md")
    if _, err := os.Stat(idx); err != nil {
        out = append(out, Finding{
            Check:   "static_site",
            Level:   LevelWarn,
            Path:    root,
            Message: "index.md missing. Add a landing page for the site",
        })
    }

    // Recommend a pages/ directory with at least one markdown file
    pagesDir := filepath.Join(root, "pages")
    if st, err := os.Stat(pagesDir); err != nil || !st.IsDir() {
        out = append(out, Finding{
            Check:   "static_site",
            Level:   LevelWarn,
            Path:    pagesDir,
            Message: "pages/ directory missing. Create pages/ with markdown content",
        })
    } else {
        entries, _ := os.ReadDir(pagesDir)
        hasMD := false
        for _, e := range entries {
            if e.IsDir() {
                continue
            }
            name := e.Name()
            if filepath.Ext(name) == ".md" || filepath.Ext(name) == ".markdown" {
                hasMD = true
                break
            }
        }
        if !hasMD {
            out = append(out, Finding{
                Check:   "static_site",
                Level:   LevelWarn,
                Path:    pagesDir,
                Message: "pages/ has no markdown files. Add at least one .md page",
            })
        }
    }

    // Recommend an assets/ directory for static files
    assetsDir := filepath.Join(root, "assets")
    if st, err := os.Stat(assetsDir); err != nil || !st.IsDir() {
        out = append(out, Finding{
            Check:   "static_site",
            Level:   LevelWarn,
            Path:    assetsDir,
            Message: "assets/ directory missing. Add assets/ for images, CSS, and JS",
        })
    }

    return out, nil
}

