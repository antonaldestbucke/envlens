// Package linter provides rule-based linting for parsed .env key sets.
// It checks for common issues such as duplicate keys, empty key names,
// and keys containing invalid characters.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Violation represents a single linting issue found in a .env file.
type Violation struct {
	Key     string
	Message string
}

// String returns a human-readable representation of the violation.
func (v Violation) String() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Message)
}

// Lint runs all linting rules against the provided key slice and returns
// any violations found. The keys slice should contain raw key names as
// returned by the parser.
func Lint(keys []string) []Violation {
	var violations []Violation

	seen := make(map[string]int)
	for i, key := range keys {
		seen[key]++
		_ = i
	}

	reported := make(map[string]bool)
	for _, key := range keys {
		if key == "" {
			violations = append(violations, Violation{Key: key, Message: "empty key name is not allowed"})
			continue
		}

		if seen[key] > 1 && !reported[key] {
			violations = append(violations, Violation{
				Key:     key,
				Message: fmt.Sprintf("duplicate key appears %d times", seen[key]),
			})
			reported[key] = true
		}

		if v := checkInvalidChars(key); v != nil {
			violations = append(violations, *v)
		}

		if strings.HasPrefix(key, "_") {
			violations = append(violations, Violation{
				Key:     key,
				Message: "key should not start with an underscore",
			})
		}
	}

	return violations
}

// checkInvalidChars returns a Violation if the key contains characters
// outside the allowed set (A-Z, a-z, 0-9, underscore).
func checkInvalidChars(key string) *Violation {
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return &Violation{
				Key:     key,
				Message: fmt.Sprintf("contains invalid character %q", r),
			}
		}
	}
	return nil
}
