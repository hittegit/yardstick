package checks

import "context"

// Level represents the severity level of a finding.
// Each check can report results at different levels depending on importance.
// - info:  informational, no action required.
// - warn:  something is suboptimal but not a failure.
// - error: policy violation or missing critical element.
type Level string

// Common severity constants used across all checks.
// Adding a new Level requires updating reports and tests accordingly.
const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Finding describes the result of a single check.
// It represents one piece of evidence found by a rule execution.
type Finding struct {
	// Check is the unique key of the check that produced this finding.
	Check string `json:"check"`

	// Level indicates how severe the finding is (info, warn, or error).
	Level Level `json:"level"`

	// Path points to the file or directory the finding concerns.
	// Optional, may be empty for project-wide findings.
	Path string `json:"path,omitempty"`

	// Message provides a short human-readable description of the issue.
	Message string `json:"message"`

	// Fixed is true if the issue was automatically corrected when --fix was used.
	Fixed bool `json:"fixed,omitempty"`
}

// Options contains runtime flags passed into each check.
// Used to control whether a check should attempt auto-remediation.
type Options struct {
	// AutoFix indicates whether yardstick should attempt to fix simple issues automatically.
	AutoFix bool
}

// Check is the interface that all yardstick checks must implement.
//
// Each check:
//   - returns a unique Key identifying itself.
//   - provides a Description for listing and documentation.
//   - implements Run, which performs the check and returns one or more Findings.
//
// Example usage pattern:
//
//   type MyCheck struct {}
//   func (MyCheck) Key() string         { return "mycheck" }
//   func (MyCheck) Description() string { return "verifies something important" }
//   func (MyCheck) Run(ctx context.Context, root string, opts Options) ([]Finding, error) {
//       // perform validation logic here
//       return []Finding{{Check: "mycheck", Level: LevelWarn, Message: "example"}}, nil
//   }
type Check interface {
	Key() string
	Description() string
	Run(ctx context.Context, root string, opts Options) ([]Finding, error)
}
