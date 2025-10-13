package report

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/erikhitt/yardstick/internal/checks"
)

// Package report converts raw findings into human or machine friendly outputs.
//
// The Output type is a lightweight DTO for JSON encoding. For terminal output,
// PrintTable renders a compact table suitable for CI logs and local use.

// Output represents a summary view of findings for JSON output.
// It flattens the data and provides quick counts per severity.
type Output struct {
	Findings []checks.Finding `json:"findings"`
	Counts   struct {
		Info  int `json:"info"`
		Warn  int `json:"warn"`
		Error int `json:"error"`
	} `json:"counts"`
}

// FromFindings converts a slice of findings into an Output, computing counts.
func FromFindings(fs []checks.Finding) Output {
	out := Output{Findings: fs}
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
	fmt.Fprintln(tw, "CHECK\tLEVEL\tPATH\tMESSAGE\tFIXED")
	for _, f := range fs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%t\n", f.Check, f.Level, f.Path, f.Message, f.Fixed)
	}
	_ = tw.Flush()
}
