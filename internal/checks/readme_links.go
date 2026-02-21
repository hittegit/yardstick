package checks

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// ReadmeLinksCheck validates local links in README.md.
type ReadmeLinksCheck struct{}

func (ReadmeLinksCheck) Key() string { return "readme_links" }

func (ReadmeLinksCheck) Description() string {
	return "Verifies README.md local file and anchor links resolve"
}

var markdownLinkPattern = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
var markdownHeadingPattern = regexp.MustCompile(`(?m)^#{1,6}\s+(.+)$`)

func (ReadmeLinksCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
	readmePath := filepath.Join(root, "README.md")
	// #nosec G304 -- root/path are intentionally user-selected scan targets.
	b, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, nil
	}

	content := string(b)
	anchors := readmeAnchors(content)
	matches := markdownLinkPattern.FindAllStringSubmatchIndex(content, -1)
	var findings []Finding

	for _, m := range matches {
		if len(m) < 4 {
			continue
		}
		// Skip image links: ![alt](...)
		if m[0] > 0 && content[m[0]-1] == '!' {
			continue
		}

		target := strings.TrimSpace(content[m[2]:m[3]])
		target = strings.Trim(target, "<>")
		if target == "" {
			continue
		}
		lower := strings.ToLower(target)
		if strings.HasPrefix(lower, "http://") ||
			strings.HasPrefix(lower, "https://") ||
			strings.HasPrefix(lower, "mailto:") ||
			strings.HasPrefix(lower, "tel:") {
			continue
		}

		if strings.HasPrefix(target, "#") {
			if _, ok := anchors[target[1:]]; !ok {
				findings = append(findings, Finding{
					Check:   "readme_links",
					Level:   LevelWarn,
					Path:    readmePath,
					Message: "README link target not found: " + target,
				})
			}
			continue
		}

		pathPart, frag := splitFragment(target)
		if pathPart == "" {
			continue
		}
		fullPath := filepath.Join(root, filepath.FromSlash(pathPart))
		info, statErr := os.Stat(fullPath)
		if statErr != nil {
			findings = append(findings, Finding{
				Check:   "readme_links",
				Level:   LevelWarn,
				Path:    readmePath,
				Message: "README link file not found: " + target,
			})
			continue
		}
		if frag != "" && !info.IsDir() && looksLikeMarkdown(fullPath) {
			ok, anchorErr := markdownFileHasAnchor(fullPath, frag)
			if anchorErr != nil {
				return nil, anchorErr
			}
			if !ok {
				findings = append(findings, Finding{
					Check:   "readme_links",
					Level:   LevelWarn,
					Path:    readmePath,
					Message: "README link anchor not found in " + pathPart + ": #" + frag,
				})
			}
		}
	}

	return findings, nil
}

func splitFragment(target string) (string, string) {
	idx := strings.IndexByte(target, '#')
	if idx == -1 {
		return target, ""
	}
	return target[:idx], target[idx+1:]
}

func looksLikeMarkdown(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

func markdownFileHasAnchor(path, anchor string) (bool, error) {
	// #nosec G304 -- path is derived from the selected repository root.
	b, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	anchors := readmeAnchors(string(b))
	_, ok := anchors[anchor]
	return ok, nil
}

func readmeAnchors(content string) map[string]struct{} {
	out := make(map[string]struct{})
	ms := markdownHeadingPattern.FindAllStringSubmatch(content, -1)
	for _, m := range ms {
		if len(m) < 2 {
			continue
		}
		anchor := headingToAnchor(m[1])
		if anchor == "" {
			continue
		}
		out[anchor] = struct{}{}
	}
	return out
}

func headingToAnchor(heading string) string {
	heading = strings.ToLower(strings.TrimSpace(heading))
	if heading == "" {
		return ""
	}

	var b strings.Builder
	prevDash := false
	for _, r := range heading {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '-' || r == '_':
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		default:
			// Drop punctuation and symbols.
		}
	}
	anchor := strings.Trim(b.String(), "-")
	return anchor
}
