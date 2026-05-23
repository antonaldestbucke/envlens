package comparator

// Result holds the comparison outcome between a reference env and a target env.
type Result struct {
	Target      string
	MissingKeys []string
	ExtraKeys   []string
}

// Compare takes a reference map of keys (e.g. from .env.example) and a target
// map of keys (e.g. from .env.production) and returns a Result describing the
// differences.
func Compare(reference map[string]string, target map[string]string, targetName string) Result {
	result := Result{
		Target:      targetName,
		MissingKeys: []string{},
		ExtraKeys:   []string{},
	}

	for key := range reference {
		if _, found := target[key]; !found {
			result.MissingKeys = append(result.MissingKeys, key)
		}
	}

	for key := range target {
		if _, found := reference[key]; !found {
			result.ExtraKeys = append(result.ExtraKeys, key)
		}
	}

	sort.Strings(result.MissingKeys)
	sort.Strings(result.ExtraKeys)

	return result
}

// HasDiff returns true when the result contains any missing or extra keys.
func (r Result) HasDiff() bool {
	return len(r.MissingKeys) > 0 || len(r.ExtraKeys) > 0
}
