// Package formatter provides utilities for formatting key lists
// and diff output for display or further processing.
package formatter

import (
	"fmt"
	"strings"
)

// DiffLine represents a single line in a formatted diff output.
type DiffLine struct {
	Symbol string // "+", "-", or " "
	Key    string
}

// FormatDiff produces a slice of DiffLines comparing missing and extra keys
// against a reference set. Missing keys are prefixed with "-", extra with "+".
func FormatDiff(reference []string, missing []string, extra []string) []DiffLine {
	missingSet := toSet(missing)
	extraSet := toSet(extra)

	var lines []DiffLine

	for _, key := range reference {
		if missingSet[key] {
			lines = append(lines, DiffLine{Symbol: "-", Key: key})
		} else {
			lines = append(lines, DiffLine{Symbol: " ", Key: key})
		}
	}

	for _, key := range extra {
		if extraSet[key] {
			lines = append(lines, DiffLine{Symbol: "+", Key: key})
		}
	}

	return lines
}

// RenderDiff converts a slice of DiffLines into a human-readable string.
func RenderDiff(lines []DiffLine) string {
	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(fmt.Sprintf("%s %s\n", line.Symbol, line.Key))
	}
	return sb.String()
}

// Truncate shortens a key list to at most n items, appending a summary line
// if items were omitted.
func Truncate(keys []string, n int) []string {
	if len(keys) <= n {
		return keys
	}
	result := make([]string, n+1)
	copy(result, keys[:n])
	result[n] = fmt.Sprintf("... and %d more", len(keys)-n)
	return result
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
