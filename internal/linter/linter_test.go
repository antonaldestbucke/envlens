package linter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/linter"
)

func violationMessages(vs []linter.Violation) []string {
	out := make([]string, len(vs))
	for i, v := range vs {
		out[i] = v.String()
	}
	return out
}

func TestLint_CleanKeys(t *testing.T) {
	keys := []string{"APP_ENV", "DB_HOST", "PORT"}
	got := linter.Lint(keys)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %v", got)
	}
}

func TestLint_EmptyKey(t *testing.T) {
	keys := []string{"APP_ENV", "", "PORT"}
	got := linter.Lint(keys)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(got), got)
	}
	if !strings.Contains(got[0].Message, "empty key") {
		t.Errorf("unexpected message: %s", got[0].Message)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	keys := []string{"APP_ENV", "DB_HOST", "APP_ENV"}
	got := linter.Lint(keys)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(got), got)
	}
	if !strings.Contains(got[0].Message, "duplicate") {
		t.Errorf("unexpected message: %s", got[0].Message)
	}
}

func TestLint_InvalidCharacter(t *testing.T) {
	keys := []string{"APP-ENV", "DB.HOST"}
	got := linter.Lint(keys)
	if len(got) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(got), got)
	}
	for _, v := range got {
		if !strings.Contains(v.Message, "invalid character") {
			t.Errorf("unexpected message: %s", v.Message)
		}
	}
}

func TestLint_LeadingUnderscore(t *testing.T) {
	keys := []string{"_INTERNAL", "PUBLIC_KEY"}
	got := linter.Lint(keys)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(got), got)
	}
	if got[0].Key != "_INTERNAL" {
		t.Errorf("expected violation for _INTERNAL, got %s", got[0].Key)
	}
}

func TestLint_MultipleIssuesSameKey(t *testing.T) {
	// key starts with underscore AND has invalid char
	keys := []string{"_BAD-KEY"}
	got := linter.Lint(keys)
	if len(got) < 2 {
		t.Fatalf("expected at least 2 violations for _BAD-KEY, got %d", len(got))
	}
}

func TestViolation_String(t *testing.T) {
	v := linter.Violation{Key: "SOME_KEY", Message: "duplicate key appears 2 times"}
	s := v.String()
	if !strings.Contains(s, "SOME_KEY") || !strings.Contains(s, "duplicate") {
		t.Errorf("unexpected String output: %s", s)
	}
}
