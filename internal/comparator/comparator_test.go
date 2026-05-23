package comparator_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/comparator"
)

func ref() map[string]string {
	return map[string]string{
		"APP_ENV":    "production",
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"SECRET_KEY": "changeme",
	}
}

func TestCompare_NoMissingKeys(t *testing.T) {
	target := map[string]string{
		"APP_ENV":    "staging",
		"DB_HOST":    "db.internal",
		"DB_PORT":    "5432",
		"SECRET_KEY": "s3cr3t",
	}
	result := comparator.Compare(ref(), target, "staging")
	if result.HasDiff() {
		t.Errorf("expected no diff, got missing=%v extra=%v", result.MissingKeys, result.ExtraKeys)
	}
}

func TestCompare_DetectsMissingKeys(t *testing.T) {
	target := map[string]string{
		"APP_ENV": "staging",
		"DB_HOST": "db.internal",
	}
	result := comparator.Compare(ref(), target, "staging")
	if len(result.MissingKeys) != 2 {
		t.Fatalf("expected 2 missing keys, got %d: %v", len(result.MissingKeys), result.MissingKeys)
	}
	if result.MissingKeys[0] != "DB_PORT" || result.MissingKeys[1] != "SECRET_KEY" {
		t.Errorf("unexpected missing keys order: %v", result.MissingKeys)
	}
}

func TestCompare_DetectsExtraKeys(t *testing.T) {
	target := map[string]string{
		"APP_ENV":    "staging",
		"DB_HOST":    "db.internal",
		"DB_PORT":    "5432",
		"SECRET_KEY": "s3cr3t",
		"EXTRA_VAR":  "unexpected",
	}
	result := comparator.Compare(ref(), target, "staging")
	if len(result.ExtraKeys) != 1 || result.ExtraKeys[0] != "EXTRA_VAR" {
		t.Errorf("expected [EXTRA_VAR] as extra keys, got %v", result.ExtraKeys)
	}
	if len(result.MissingKeys) != 0 {
		t.Errorf("expected no missing keys, got %v", result.MissingKeys)
	}
}

func TestCompare_EmptyTarget(t *testing.T) {
	target := map[string]string{}
	result := comparator.Compare(ref(), target, "empty")
	if len(result.MissingKeys) != 4 {
		t.Errorf("expected 4 missing keys, got %d", len(result.MissingKeys))
	}
}

func TestResult_HasDiff(t *testing.T) {
	noDiff := comparator.Result{Target: "t", MissingKeys: []string{}, ExtraKeys: []string{}}
	if noDiff.HasDiff() {
		t.Error("expected HasDiff to be false")
	}
	hasDiff := comparator.Result{Target: "t", MissingKeys: []string{"FOO"}, ExtraKeys: []string{}}
	if !hasDiff.HasDiff() {
		t.Error("expected HasDiff to be true")
	}
}
