// Package sorter provides utilities for sorting and ordering
// comparison reports and their keys for deterministic output.
package sorter

import (
	"sort"

	"github.com/yourorg/envlens/internal/comparator"
)

// Order defines the sort order for report keys.
type Order string

const (
	// OrderAlpha sorts keys alphabetically (ascending).
	OrderAlpha Order = "alpha"
	// OrderAlphaDesc sorts keys alphabetically (descending).
	OrderAlphaDesc Order = "alpha-desc"
	// OrderNone preserves the original insertion order.
	OrderNone Order = "none"
)

// SortReport returns a new Report with MissingKeys and ExtraKeys sorted
// according to the given Order. The original report is not modified.
func SortReport(r comparator.Report, order Order) comparator.Report {
	if order == OrderNone {
		return r
	}

	sorted := comparator.Report{
		Target:      r.Target,
		MissingKeys: copySlice(r.MissingKeys),
		ExtraKeys:   copySlice(r.ExtraKeys),
	}

	switch order {
	case OrderAlpha:
		sort.Strings(sorted.MissingKeys)
		sort.Strings(sorted.ExtraKeys)
	case OrderAlphaDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(sorted.MissingKeys)))
		sort.Sort(sort.Reverse(sort.StringSlice(sorted.ExtraKeys)))
	}

	return sorted
}

// SortReports applies SortReport to a slice of reports.
func SortReports(reports []comparator.Report, order Order) []comparator.Report {
	result := make([]comparator.Report, len(reports))
	for i, r := range reports {
		result[i] = SortReport(r, order)
	}
	return result
}

func copySlice(s []string) []string {
	if s == nil {
		return nil
	}
	out := make([]string, len(s))
	copy(out, s)
	return out
}
