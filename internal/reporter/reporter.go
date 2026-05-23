package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/comparator"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the results of comparing a reference env against targets.
type Report struct {
	Reference string
	Results   []comparator.Result
}

// Write writes the report to the given writer in the specified format.
func Write(w io.Writer, report Report, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, report)
	default:
		return writeText(w, report)
	}
}

func writeText(w io.Writer, report Report) error {
	fmt.Fprintf(w, "envlens report (reference: %s)\n", report.Reference)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	allClean := true
	for _, result := range report.Results {
		if len(result.Missing) == 0 && len(result.Extra) == 0 {
			fmt.Fprintf(w, "[OK] %s\n", result.Target)
			continue
		}
		allClean = false
		fmt.Fprintf(w, "[DIFF] %s\n", result.Target)
		for _, key := range result.Missing {
			fmt.Fprintf(w, "  - missing: %s\n", key)
		}
		for _, key := range result.Extra {
			fmt.Fprintf(w, "  + extra:   %s\n", key)
		}
	}

	if allClean {
		fmt.Fprintln(w, "All targets are in sync with the reference.")
	}
	return nil
}

func writeJSON(w io.Writer, report Report) error {
	fmt.Fprintf(w, "{\n  \"reference\": %q,\n  \"results\": [\n", report.Reference)
	for i, result := range report.Results {
		missingJSON := keysToJSONArray(result.Missing)
		extraJSON := keysToJSONArray(result.Extra)
		comma := ","
		if i == len(report.Results)-1 {
			comma = ""
		}
		fmt.Fprintf(w, "    {\"target\": %q, \"missing\": %s, \"extra\": %s}%s\n",
			result.Target, missingJSON, extraJSON, comma)
	}
	fmt.Fprintln(w, "  ]\n}")
	return nil
}

func keysToJSONArray(keys []string) string {
	if len(keys) == 0 {
		return "[]"
	}
	quoted := make([]string, len(keys))
	for i, k := range keys {
		quoted[i] = fmt.Sprintf("%q", k)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

// Summary returns a brief human-readable summary of the report, indicating
// how many targets were checked and how many had differences.
func Summary(report Report) string {
	total := len(report.Results)
	diffCount := 0
	for _, result := range report.Results {
		if len(result.Missing) > 0 || len(result.Extra) > 0 {
			diffCount++
		}
	}
	if diffCount == 0 {
		return fmt.Sprintf("%d/%d targets in sync", total, total)
	}
	return fmt.Sprintf("%d/%d targets have differences", diffCount, total)
}
