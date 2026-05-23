// Package summarizer provides aggregation utilities for envlens comparison
// reports. It converts a slice of [comparator.Report] values into a single
// [Result] that captures per-target and global statistics.
//
// Typical usage:
//
//	reports := comparator.CompareAll(reference, targets)
//	result  := summarizer.Summarize(reports)
//	if result.HasIssues() {
//		os.Exit(1)
//	}
package summarizer
