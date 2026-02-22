package report

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/hittegit/yardstick/internal/checks"
)

func TestFromFindings_Counts(t *testing.T) {
	fs := []checks.Finding{
		{Check: "a", Level: checks.LevelWarn},
		{Check: "b", Level: checks.LevelWarn},
		{Check: "c", Level: checks.LevelInfo},
		{Check: "d", Level: checks.LevelError},
	}
	out := FromFindings(fs)
	if out.Counts.Info != 1 || out.Counts.Warn != 2 || out.Counts.Error != 1 {
		t.Fatalf("unexpected counts: %+v", out.Counts)
	}
}

// normalize splits a tabwriter line into columns by collapsing runs of 2+ spaces.
func normalize(line string) []string {
	line = strings.TrimSpace(line)
	re := regexp.MustCompile(`\s{2,}`)
	line = re.ReplaceAllString(line, "\t")
	return strings.Split(line, "\t")
}

func TestPrintTable_SortsAndFormats(t *testing.T) {
	fs := []checks.Finding{
		{Check: "b", Level: checks.LevelWarn, Path: "z", Message: "m2"},
		{Check: "a", Level: checks.LevelInfo, Path: "b", Message: "m1"},
		{Check: "a", Level: checks.LevelWarn, Path: "a", Message: "m0"},
	}
	var buf bytes.Buffer
	PrintTable(&buf, fs)
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 { // header + 3 rows
		t.Fatalf("unexpected number of lines: %d\n%s", len(lines), out)
	}
	headerCols := normalize(lines[0])
	if len(headerCols) < 5 || headerCols[0] != "CHECK" || headerCols[1] != "LEVEL" || headerCols[2] != "PATH" || headerCols[3] != "MESSAGE" || headerCols[4] != "FIXED" {
		t.Fatalf("unexpected header columns: %#v (line: %q)", headerCols, lines[0])
	}
	// Sorted by check then path: a/a, a/b, b/z
	row1 := normalize(lines[1])
	row2 := normalize(lines[2])
	row3 := normalize(lines[3])
	if row1[0] != "a" || row1[2] != "a" {
		t.Fatalf("unexpected first row columns: %#v", row1)
	}
	if row2[0] != "a" || row2[2] != "b" {
		t.Fatalf("unexpected second row columns: %#v", row2)
	}
	if row3[0] != "b" || row3[2] != "z" {
		t.Fatalf("unexpected third row columns: %#v", row3)
	}
}

func TestOutputJSON_IncludesStableFindingFields(t *testing.T) {
	out := FromFindings([]checks.Finding{
		{Check: "manifest", Level: checks.LevelInfo, Path: "/repo/go.mod", Message: "ok"},
	})
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("marshal output: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"path":"/repo/go.mod"`) {
		t.Fatalf("missing path field: %s", s)
	}
	if !strings.Contains(s, `"fixed":false`) {
		t.Fatalf("missing fixed field: %s", s)
	}
}

func TestFromRun_SummaryAllPassed(t *testing.T) {
	out := FromRun([]CheckStatus{
		{Check: "manifest", Status: "pass"},
		{Check: "readme", Status: "pass"},
	}, nil)
	if out.Summary != "All checks passed (2/2)." {
		t.Fatalf("unexpected summary: %q", out.Summary)
	}
}

func TestFromRun_SummaryWithFailures(t *testing.T) {
	out := FromRun([]CheckStatus{
		{Check: "manifest", Status: "pass"},
		{Check: "readme", Status: "fail"},
		{Check: "license", Status: "fail"},
	}, nil)
	if out.Summary != "2 of 3 checks failed." {
		t.Fatalf("unexpected summary: %q", out.Summary)
	}
}

func TestPrintVerboseTable_IncludesStatusAndGuidance(t *testing.T) {
	out := Output{
		Summary: "1 of 2 checks failed.",
		Checks: []CheckStatus{
			{Check: "manifest", Status: "pass", Findings: 0},
			{
				Check:        "readme",
				Status:       "fail",
				Level:        checks.LevelWarn,
				Findings:     1,
				WhyImportant: "docs matter",
				HowToResolve: "add README sections",
			},
		},
		Findings: []checks.Finding{
			{Check: "readme", Level: checks.LevelWarn, Path: "/repo/README.md", Message: "Missing section: ## CI"},
		},
	}

	var buf bytes.Buffer
	PrintVerboseTable(&buf, out)
	s := buf.String()
	if !strings.Contains(s, "SUMMARY: 1 of 2 checks failed.") {
		t.Fatalf("missing summary: %s", s)
	}
	if !strings.Contains(s, "readme") || !strings.Contains(s, "fail") {
		t.Fatalf("missing failed check status: %s", s)
	}
	if !strings.Contains(s, "Why: docs matter") || !strings.Contains(s, "Fix: add README sections") {
		t.Fatalf("missing guidance in verbose table: %s", s)
	}
	if !strings.Contains(s, "FINDINGS") {
		t.Fatalf("missing findings section: %s", s)
	}
}
