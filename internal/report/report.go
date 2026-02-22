package report

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/hittegit/yardstick/internal/checks"
)

// Package report converts raw findings into human or machine friendly outputs.
//
// The Output type is a lightweight DTO for JSON encoding. For terminal output,
// PrintTable renders a compact table suitable for CI logs and local use.

// Output represents a summary view of findings for JSON output.
// It flattens the data and provides quick counts per severity.
type Output struct {
	Summary  string           `json:"summary"`
	Checks   []CheckStatus    `json:"checks"`
	Findings []checks.Finding `json:"findings"`
	Counts   struct {
		Info  int `json:"info"`
		Warn  int `json:"warn"`
		Error int `json:"error"`
	} `json:"counts"`
}

// CheckStatus describes pass/fail status for an executed check.
type CheckStatus struct {
	Check        string       `json:"check"`
	Description  string       `json:"description"`
	Status       string       `json:"status"`
	Level        checks.Level `json:"level,omitempty"`
	Findings     int          `json:"findings"`
	WhyImportant string       `json:"why_important,omitempty"`
	HowToResolve string       `json:"how_to_resolve,omitempty"`
}

// FromFindings converts a slice of findings into an Output, computing counts.
func FromFindings(fs []checks.Finding) Output {
	out := Output{Findings: fs}
	out.Summary = "Findings-only output"
	for _, f := range fs {
		switch f.Level {
		case checks.LevelInfo:
			out.Counts.Info++
		case checks.LevelWarn:
			out.Counts.Warn++
		case checks.LevelError:
			out.Counts.Error++
		}
	}
	return out
}

// FromRun builds an output payload with per-check statuses and findings.
func FromRun(statuses []CheckStatus, fs []checks.Finding) Output {
	out := FromFindings(fs)
	out.Checks = statuses

	total := len(statuses)
	failed := 0
	for _, s := range statuses {
		if s.Status == "fail" {
			failed++
		}
	}
	if total == 0 {
		out.Summary = "No checks were executed."
		return out
	}
	if failed == 0 {
		out.Summary = fmt.Sprintf("All checks passed (%d/%d).", total, total)
		return out
	}
	out.Summary = fmt.Sprintf("%d of %d checks failed.", failed, total)
	return out
}

// PrintTable writes a simple tabular report to the provided writer.
// The format is stable and greppable, which helps in CI logs.
func PrintTable(w io.Writer, fs []checks.Finding) {
	// Stable order: by check key, then by path. Keeps diffs predictable.
	sort.Slice(fs, func(i, j int) bool {
		if fs[i].Check == fs[j].Check {
			return fs[i].Path < fs[j].Path
		}
		return fs[i].Check < fs[j].Check
	})

	// tabwriter keeps columns aligned without manual padding.
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "CHECK\tLEVEL\tPATH\tMESSAGE\tFIXED")
	for _, f := range fs {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%t\n", f.Check, f.Level, f.Path, f.Message, f.Fixed)
	}
	_ = tw.Flush() //nolint:errcheck // best-effort flush for tabwriter
}

// PrintVerboseTable writes a summary plus per-check status details.
func PrintVerboseTable(w io.Writer, out Output) {
	_, _ = fmt.Fprintf(w, "SUMMARY: %s\n\n", out.Summary)

	statuses := append([]CheckStatus(nil), out.Checks...)
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Check < statuses[j].Check
	})

	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "CHECK\tSTATUS\tLEVEL\tFINDINGS\tDETAILS")
	for _, s := range statuses {
		level := string(s.Level)
		if level == "" {
			level = "-"
		}
		details := "OK"
		if s.Status == "fail" {
			details = "Why: " + s.WhyImportant + " | Fix: " + s.HowToResolve
		}
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%s\n", s.Check, s.Status, level, s.Findings, details)
	}
	_ = tw.Flush() //nolint:errcheck // best-effort flush for tabwriter

	if len(out.Findings) == 0 {
		_, _ = fmt.Fprintln(w, "\nNo findings. Repository hygiene checks look good.")
		return
	}

	_, _ = fmt.Fprintln(w, "\nFINDINGS")
	PrintTable(w, out.Findings)
}
