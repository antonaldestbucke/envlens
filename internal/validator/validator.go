// Package validator checks .env values against simple rules such as
// non-empty, numeric, or boolean constraints.
package validator

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule describes a validation rule applied to a specific key.
type Rule struct {
	Key      string
	Required bool
	Kind     string // "any" | "bool" | "int" | "nonempty"
}

// Violation records a single failed validation.
type Violation struct {
	Key     string
	Rule    string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s [%s]: %s", v.Key, v.Rule, v.Message)
}

// Validate applies rules to the provided key-value map and returns all
// violations found. An empty slice means the map is fully compliant.
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, r := range rules {
		val, exists := env[r.Key]

		if r.Required && !exists {
			violations = append(violations, Violation{
				Key:     r.Key,
				Rule:    "required",
				Message: "key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		switch strings.ToLower(r.Kind) {
		case "nonempty":
			if strings.TrimSpace(val) == "" {
				violations = append(violations, Violation{
					Key:     r.Key,
					Rule:    "nonempty",
					Message: "value must not be empty",
				})
			}
		case "int":
			if _, err := strconv.Atoi(strings.TrimSpace(val)); err != nil {
				violations = append(violations, Violation{
					Key:     r.Key,
					Rule:    "int",
					Message: fmt.Sprintf("value %q is not a valid integer", val),
				})
			}
		case "bool":
			norm := strings.ToLower(strings.TrimSpace(val))
			if norm != "true" && norm != "false" && norm != "1" && norm != "0" {
				violations = append(violations, Violation{
					Key:     r.Key,
					Rule:    "bool",
					Message: fmt.Sprintf("value %q is not a valid boolean", val),
				})
			}
		}
	}

	return violations
}
