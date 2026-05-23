// Package filter provides utilities for filtering comparison results
// based on user-defined key patterns (prefix, suffix, or glob-style).
package filter

import (
	"strings"

	"github.com/envlens/internal/comparator"
)

// Options holds the filtering configuration.
type Options struct {
	// Prefixes restricts results to keys that start with any of these strings.
	Prefixes []string
	// Suffixes restricts results to keys that end with any of these strings.
	Suffixes []string
}

// Apply returns a new slice of Reports with keys filtered according to opts.
// If opts contains no prefixes and no suffixes, the original reports are
// returned unchanged.
func Apply(reports []comparator.Report, opts Options) []comparator.Report {
	if len(opts.Prefixes) == 0 && len(opts.Suffixes) == 0 {
		return reports
	}

	out := make([]comparator.Report, 0, len(reports))
	for _, r := range reports {
		filtered := comparator.Report{
			Target:      r.Target,
			MissingKeys: filterKeys(r.MissingKeys, opts),
			ExtraKeys:   filterKeys(r.ExtraKeys, opts),
		}
		out = append(out, filtered)
	}
	return out
}

// filterKeys returns only the keys from ks that match at least one of the
// patterns described by opts.
func filterKeys(ks []string, opts Options) []string {
	var result []string
	for _, k := range ks {
		if matchesAny(k, opts) {
			result = append(result, k)
		}
	}
	return result
}

// matchesAny reports whether key matches any prefix or suffix in opts.
func matchesAny(key string, opts Options) bool {
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, s := range opts.Suffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	return false
}
