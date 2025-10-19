package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hittegit/yardstick/internal/checks"
	"github.com/hittegit/yardstick/internal/report"
	"github.com/hittegit/yardstick/internal/scaffold"
)

// Define command-line flags for configuration.
// These control how yardstick runs and what output format it uses.
var (
	flagFormat   = flag.String("format", "table", "output format: table or json")
	flagPath     = flag.String("path", ".", "path to scan")
	flagFix      = flag.Bool("fix", false, "attempt to auto-fix selected issues")
	flagScaffold = flag.Bool("scaffold", false, "write common repo files if missing")
	flagStrict   = flag.Bool("strict", false, "nonzero exit if any warn-level finding exists")
	flagOnly     = flag.String("only", "", "comma-separated list of checks to run, empty means all")
	flagList     = flag.Bool("list", false, "list available checks")
)

// main is the CLI entrypoint. It parses flags and executes the program logic.
func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "yardstick error: %v\n", err)
		os.Exit(2)
	}
}

// run performs the main logic of yardstick, executing checks and printing results.
func run(ctx context.Context) error {
	// If user requests the list of available checks, show and exit.
	if *flagList {
		for _, c := range checks.All() {
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
		for _, k := range splitCSV(*flagOnly) {
			sel[k] = struct{}{}
		}
	}

	// Optionally scaffold missing standard files like README.md or LICENSE.
	if *flagScaffold {
		if err := scaffold.Run(root); err != nil {
			return fmt.Errorf("scaffold: %w", err)
		}
	}

	// Run all registered checks (or a subset if specified).
	var findings []checks.Finding
	for _, c := range checks.All() {
		// Skip any checks not listed in the --only flag.
		if sel != nil {
			if _, ok := sel[c.Key()]; !ok {
				continue
			}
		}

		// Execute each check with optional autofix.
		fs, err := c.Run(ctx, root, checks.Options{AutoFix: *flagFix})
		if err != nil {
			return fmt.Errorf("check %s: %w", c.Key(), err)
		}
		findings = append(findings, fs...)
	}

	// Render the report in the requested format.
	switch *flagFormat {
	case "json":
		// Machine-readable output for CI pipelines.
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report.FromFindings(findings)); err != nil {
			return err
		}
	default:
		// Human-readable table format.
		report.PrintTable(os.Stdout, findings)
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

// splitCSV is a minimal helper for parsing comma-separated values without trimming spaces.
func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if start < i {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}
