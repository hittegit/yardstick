package checks

// All returns the complete list of registered checks.
//
// This registry is how yardstick discovers which checks to execute.
// Each check implements the Check interface (see types.go), and must be
// explicitly listed here to be included in the scanning process.
//
// When you create a new check:
//   1. Define it in a new file under internal/checks (e.g., mycheck.go).
//   2. Implement the Check interface (Key, Description, Run).
//   3. Add an instance to the slice below.
//
// Keeping this list explicit helps ensure that all checks are intentionally
// loaded, rather than automatically discovered, which improves predictability
// and makes CI output easier to reason about.
func All() []Check {
    return []Check{
        ManifestCheck{},  // Detect project ecosystem by scanning for common manifests
        ReadmeCheck{},    // Ensures README.md exists and has required sections
        LicenseCheck{},   // Ensures LICENSE file is present
        GitIgnoreCheck{}, // Ensures .gitignore covers common entries
        ChangelogCheck{}, // Ensures CHANGELOG.md exists
    }
}
