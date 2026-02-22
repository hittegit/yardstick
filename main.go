package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hittegit/yardstick/internal/checks"
	"github.com/hittegit/yardstick/internal/report"
)

// Define command-line flags for configuration.
// These control how yardstick runs and what output format it uses.
var (
	flagFormat  = flag.String("format", "table", "output format: table or json")
	flagPath    = flag.String("path", ".", "path to scan")
	flagStrict  = flag.Bool("strict", false, "nonzero exit if any warn-level finding exists")
	flagOnly    = flag.String("only", "", "comma-separated list of checks to run, empty means all")
	flagList    = flag.Bool("list", false, "list available checks")
	flagVersion = flag.Bool("version", false, "print version and exit")
)

// Build-time variables injected via -ldflags at release time.
// Default values are for local development.
var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

// main is the CLI entrypoint. It parses flags and executes the program logic.
func main() {
	flag.Parse()
	if *flagVersion {
		fmt.Printf("yardstick %s (commit %s, built %s)\n", buildVersion, buildCommit, buildDate)
		return
	}
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "yardstick error: %v\n", err)
		os.Exit(2)
	}
}

// run performs the main logic of yardstick, executing checks and printing results.
func run(ctx context.Context) error {
	if *flagFormat != "table" && *flagFormat != "json" {
		return fmt.Errorf("invalid -format %q, expected table or json", *flagFormat)
	}

	allChecks := checks.All()
	available := make(map[string]struct{}, len(allChecks))
	for _, c := range allChecks {
		available[c.Key()] = struct{}{}
	}

	// If user requests the list of available checks, show and exit.
	if *flagList {
		for _, c := range allChecks {
			fmt.Printf("%s - %s\n", c.Key(), c.Description())
		}
		return nil
	}

	// Resolve absolute path for scanning.
	root, err := filepath.Abs(*flagPath)
	if err != nil {
		return err
	}

	// Parse the comma-separated list of specific checks to run (if provided).
	var sel map[string]struct{}
	if *flagOnly != "" {
		sel = make(map[string]struct{})
		var unknown []string
		for _, k := range splitCSV(*flagOnly) {
			if _, ok := available[k]; !ok {
				unknown = append(unknown, k)
				continue
			}
			sel[k] = struct{}{}
		}
		if len(unknown) > 0 {
			sort.Strings(unknown)
			return fmt.Errorf("unknown check key(s) in -only: %s", strings.Join(unknown, ", "))
		}
		if len(sel) == 0 {
			return errors.New("no valid checks selected via -only")
		}
	}

	// Run all registered checks (or a subset if specified).
	var findings []checks.Finding
	var checkStatuses []report.CheckStatus
	for _, c := range allChecks {
		// Skip any checks not listed in the --only flag.
		if sel != nil {
			if _, ok := sel[c.Key()]; !ok {
				continue
			}
		}

		// Execute each check; yardstick is read-only so AutoFix is ignored.
		fs, err := c.Run(ctx, root, checks.Options{AutoFix: false})
		if err != nil {
			return fmt.Errorf("check %s: %w", c.Key(), err)
		}
		findings = append(findings, fs...)
		checkStatuses = append(checkStatuses, statusForCheck(c, fs))
	}

	out := report.FromRun(checkStatuses, findings)

	// Render the report in the requested format.
	switch *flagFormat {
	case "json":
		// Machine-readable output for CI pipelines.
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			return err
		}
	case "table":
		// Human-readable table format with per-check status and guidance.
		report.PrintVerboseTable(os.Stdout, out)
	}

	// Evaluate whether any errors or warnings should cause a nonzero exit.
	var hasError, hasWarn bool
	for _, f := range findings {
		if f.Level == checks.LevelError {
			hasError = true
		}
		if f.Level == checks.LevelWarn {
			hasWarn = true
		}
	}

	// Exit code logic for CI integration.
	// - Errors always fail.
	// - Warnings fail only when --strict is enabled.
	if hasError || (*flagStrict && hasWarn) {
		return errors.New("policy violations found")
	}
	return nil
}

func statusForCheck(c checks.Check, fs []checks.Finding) report.CheckStatus {
	status := report.CheckStatus{
		Check:       c.Key(),
		Description: c.Description(),
		Status:      "pass",
		Findings:    len(fs),
	}
	if len(fs) == 0 {
		return status
	}

	status.Level = highestLevel(fs)
	if status.Level == checks.LevelWarn || status.Level == checks.LevelError {
		status.Status = "fail"
		if g, ok := checks.GuidanceForCheck(c.Key()); ok {
			status.WhyImportant = g.WhyImportant
			status.HowToResolve = g.HowToResolve
		}
	}
	return status
}

func highestLevel(fs []checks.Finding) checks.Level {
	level := checks.LevelInfo
	for _, f := range fs {
		switch f.Level {
		case checks.LevelError:
			return checks.LevelError
		case checks.LevelWarn:
			level = checks.LevelWarn
		}
	}
	return level
}

// splitCSV parses comma-separated values and trims surrounding whitespace.
func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if start < i {
				trimmed := strings.TrimSpace(s[start:i])
				if trimmed != "" {
					out = append(out, trimmed)
				}
			}
			start = i + 1
		}
	}
	return out
}
