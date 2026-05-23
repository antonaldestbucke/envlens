package summarizer

import (
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable summary table to w.
func WriteText(w io.Writer, res Result) {
	fmt.Fprintf(w, "Summary: %d target(s) — %d clean, %d with issues\n",
		res.TotalTargets, res.CleanTargets, res.DirtyTargets)

	if len(res.PerTarget) == 0 {
		return
	}

	fmt.Fprintln(w, strings.Repeat("-", 48))
	for _, pt := range res.PerTarget {
		status := "✓"
		if !pt.Clean {
			status = "✗"
		}
		fmt.Fprintf(w, "  %s  %-24s missing=%-3d extra=%d\n",
			status, pt.Name, pt.MissingCount, pt.ExtraCount)
	}
	fmt.Fprintln(w, strings.Repeat("-", 48))

	if res.HasIssues() {
		fmt.Fprintf(w, "Total missing keys: %d  |  Total extra keys: %d\n",
			res.TotalMissing, res.TotalExtra)
	} else {
		fmt.Fprintln(w, "All targets are in sync with the reference.")
	}
}
