package checks

// Guidance captures why a check matters and how to resolve failures.
type Guidance struct {
	WhyImportant string
	HowToResolve string
}

var guidanceByCheck = map[string]Guidance{
	"manifest": {
		WhyImportant: "A manifest helps tools and contributors understand the project stack and dependency model.",
		HowToResolve: "Add a standard manifest for your ecosystem, for example go.mod, package.json, pyproject.toml, or Cargo.toml.",
	},
	"static_site": {
		WhyImportant: "A minimal static-site structure improves reliability for builds, hosting, and navigation.",
		HowToResolve: "Add index.md, a pages/ directory with markdown content, and an assets/ directory for static files.",
	},
	"readme": {
		WhyImportant: "A complete README reduces onboarding friction and clarifies project usage for contributors and consumers.",
		HowToResolve: "Create or update README.md to include Overview, Installation, Usage, CI, and License sections.",
	},
	"readme_links": {
		WhyImportant: "Broken README links reduce trust and block readers from important docs and setup instructions.",
		HowToResolve: "Fix invalid local links and anchors in README.md so each referenced file and heading exists.",
	},
	"license": {
		WhyImportant: "A license defines legal reuse terms and protects both maintainers and users.",
		HowToResolve: "Add a LICENSE file with the license your project intends to use, for example MIT or Apache-2.0.",
	},
	"gitignore": {
		WhyImportant: "A .gitignore prevents accidental commits of build artifacts, secrets, and machine-local files.",
		HowToResolve: "Add a .gitignore tuned to your stack to exclude artifacts, editor files, and OS-specific files.",
	},
	"changelog": {
		WhyImportant: "A changelog helps users and maintainers track behavior changes across releases.",
		HowToResolve: "Add CHANGELOG.md and document notable changes per release, ideally using Keep a Changelog format.",
	},
	"codeowners": {
		WhyImportant: "CODEOWNERS clarifies review responsibility and improves governance in collaborative repositories.",
		HowToResolve: "Add CODEOWNERS in a standard location and map key paths to responsible reviewers.",
	},
	"security_policy": {
		WhyImportant: "A security policy provides a clear process for responsible vulnerability reporting.",
		HowToResolve: "Add SECURITY.md with reporting channels, expected response timelines, and disclosure expectations.",
	},
	"contributing": {
		WhyImportant: "Contribution guidelines reduce confusion and improve consistency for incoming changes.",
		HowToResolve: "Add CONTRIBUTING.md covering setup, coding standards, test expectations, and PR process.",
	},
	"ci_workflow": {
		WhyImportant: "CI workflows enforce baseline quality checks before changes are merged.",
		HowToResolve: "Add at least one workflow file under .github/workflows to run build and test checks.",
	},
}

// GuidanceForCheck returns guidance for a check key.
func GuidanceForCheck(key string) (Guidance, bool) {
	g, ok := guidanceByCheck[key]
	return g, ok
}
