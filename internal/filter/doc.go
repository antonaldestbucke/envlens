// Package filter provides post-processing utilities for envlens comparison
// results.
//
// After [comparator.Compare] produces a set of [comparator.Report] values,
// filter.Apply can narrow those reports to only the keys that match a set of
// user-supplied prefixes or suffixes. This is useful when a project uses
// naming conventions such as "DB_", "AWS_", or "_SECRET" to group related
// environment variables and the operator only wants to audit one group at a
// time.
//
// Example:
//
//	opts := filter.Options{
//		Prefixes: []string{"DB_", "REDIS_"},
//	}
//	filtered := filter.Apply(reports, opts)
package filter
