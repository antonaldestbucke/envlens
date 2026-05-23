// Package summarizer aggregates comparison reports into a high-level
// summary suitable for exit-code decisions and human-readable output.
package summarizer

import "github.com/yourorg/envlens/internal/comparator"

// Result holds aggregated statistics across all compared targets.
type Result struct {
	TotalTargets  int
	CleanTargets  int
	DirtyTargets  int
	TotalMissing  int
	TotalExtra    int
	PerTarget     []TargetSummary
}

// TargetSummary holds per-target statistics.
type TargetSummary struct {
	Name         string
	MissingCount int
	ExtraCount   int
	Clean        bool
}

// Summarize aggregates a slice of comparator.Report values into a Result.
func Summarize(reports []comparator.Report) Result {
	res := Result{
		TotalTargets: len(reports),
	}

	for _, r := range reports {
		ts := TargetSummary{
			Name:         r.Target,
			MissingCount: len(r.Missing),
			ExtraCount:   len(r.Extra),
		}
		ts.Clean = ts.MissingCount == 0 && ts.ExtraCount == 0

		res.TotalMissing += ts.MissingCount
		res.TotalExtra += ts.ExtraCount

		if ts.Clean {
			res.CleanTargets++
		} else {
			res.DirtyTargets++
		}

		res.PerTarget = append(res.PerTarget, ts)
	}

	return res
}

// HasIssues returns true when any target has missing or extra keys.
func (r Result) HasIssues() bool {
	return r.DirtyTargets > 0
}
