// Package merger combines multiple env file key sets into a unified
// reference, optionally flagging keys that appear only in some sources.
package merger

import "github.com/user/envlens/internal/loader"

// Result holds the merged key set and metadata about key origin.
type Result struct {
	// Keys is the deduplicated union of all keys across sources.
	Keys []string
	// Partial maps each key to the list of source names that define it.
	Partial map[string][]string
}

// Merge combines the key sets from all loaded env files into a single Result.
// Keys present in every source are considered "complete"; keys present in only
// some sources are recorded in Partial so callers can surface warnings.
func Merge(envs []loader.Env) Result {
	seen := make(map[string][]string) // key -> source names

	for _, e := range envs {
		for k := range e.Keys {
			seen[k] = append(seen[k], e.Name)
		}
	}

	total := len(envs)
	partial := make(map[string][]string)
	keys := make([]string, 0, len(seen))

	for k, sources := range seen {
		keys = append(keys, k)
		if len(sources) < total {
			partial[k] = sources
		}
	}

	return Result{
		Keys:    keys,
		Partial: partial,
	}
}

// IsPartial reports whether key k is absent from at least one source.
func (r Result) IsPartial(k string) bool {
	_, ok := r.Partial[k]
	return ok
}
