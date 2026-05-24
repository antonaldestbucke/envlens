// Package differ provides utilities for computing line-level diffs
// between two sets of .env keys, producing structured DiffLine results
// suitable for rendering or further processing.
package differ

// LineKind represents the type of a diff line.
type LineKind int

const (
	KindContext LineKind = iota // key present in both
	KindMissing                 // key missing from target
	KindExtra                   // key extra in target
)

// DiffLine represents a single line in a computed diff.
type DiffLine struct {
	Key  string
	Kind LineKind
}

// Result holds the full diff between a reference and a target env.
type Result struct {
	Target string
	Lines  []DiffLine
}

// Compute produces a Result by comparing refKeys against targetKeys.
// Keys are iterated in the order provided by refKeys for context/missing,
// then extra keys from targetKeys are appended at the end.
func Compute(target string, refKeys, targetKeys []string) Result {
	targetSet := toSet(targetKeys)
	refSet := toSet(refKeys)

	var lines []DiffLine

	for _, k := range refKeys {
		if targetSet[k] {
			lines = append(lines, DiffLine{Key: k, Kind: KindContext})
		} else {
			lines = append(lines, DiffLine{Key: k, Kind: KindMissing})
		}
	}

	for _, k := range targetKeys {
		if !refSet[k] {
			lines = append(lines, DiffLine{Key: k, Kind: KindExtra})
		}
	}

	return Result{Target: target, Lines: lines}
}

// toSet converts a slice of strings into a lookup map.
func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
