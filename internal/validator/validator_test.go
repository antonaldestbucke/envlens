package validator_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/validator"
)

func env(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestValidate_NoViolations(t *testing.T) {
	rules := []validator.Rule{
		{Key: "PORT", Required: true, Kind: "int"},
		{Key: "DEBUG", Required: true, Kind: "bool"},
	}
	violations := validator.Validate(env("PORT", "8080", "DEBUG", "true"), rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_RequiredMissingKey(t *testing.T) {
	rules := []validator.Rule{
		{Key: "SECRET", Required: true, Kind: "nonempty"},
	}
	violations := validator.Validate(env(), rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "required" {
		t.Errorf("expected rule=required, got %q", violations[0].Rule)
	}
}

func TestValidate_NonEmptyFails(t *testing.T) {
	rules := []validator.Rule{
		{Key: "APP_NAME", Kind: "nonempty"},
	}
	violations := validator.Validate(env("APP_NAME", "   "), rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "nonempty" {
		t.Errorf("expected rule=nonempty, got %q", violations[0].Rule)
	}
}

func TestValidate_IntFails(t *testing.T) {
	rules := []validator.Rule{
		{Key: "PORT", Kind: "int"},
	}
	violations := validator.Validate(env("PORT", "abc"), rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "int" {
		t.Errorf("expected rule=int, got %q", violations[0].Rule)
	}
}

func TestValidate_BoolFails(t *testing.T) {
	rules := []validator.Rule{
		{Key: "ENABLE_FEATURE", Kind: "bool"},
	}
	violations := validator.Validate(env("ENABLE_FEATURE", "yes"), rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_BoolAcceptsNumeric(t *testing.T) {
	rules := []validator.Rule{
		{Key: "FLAG", Kind: "bool"},
	}
	for _, val := range []string{"0", "1", "true", "false", "TRUE", "FALSE"} {
		v := validator.Validate(env("FLAG", val), rules)
		if len(v) != 0 {
			t.Errorf("value %q should be valid bool, got violation: %v", val, v)
		}
	}
}

func TestValidate_OptionalMissingKeySkipped(t *testing.T) {
	rules := []validator.Rule{
		{Key: "OPTIONAL", Required: false, Kind: "int"},
	}
	violations := validator.Validate(env(), rules)
	if len(violations) != 0 {
		t.Fatalf("optional missing key should not produce violations, got %v", violations)
	}
}
