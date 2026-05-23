package main

import (
	"testing"

	"github.com/yourorg/envlens/internal/comparator"
)

func TestSplitAndTrim_Basic(t *testing.T) {
	got := splitAndTrim(".env.staging, .env.prod")
	if len(got) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(got))
	}
	if got[0] != ".env.staging" {
		t.Errorf("expected .env.staging, got %q", got[0])
	}
	if got[1] != ".env.prod" {
		t.Errorf("expected .env.prod, got %q", got[1])
	}
}

func TestSplitAndTrim_EmptySegments(t *testing.T) {
	got := splitAndTrim("a,,b, ,c")
	if len(got) != 3 {
		t.Fatalf("expected 3 parts, got %d: %v", len(got), got)
	}
}

func TestSplitAndTrim_Empty(t *testing.T) {
	got := splitAndTrim("")
	if len(got) != 0 {
		t.Fatalf("expected 0 parts, got %d", len(got))
	}
}

func TestSplitAndTrim_Whitespace(t *testing.T) {
	got := splitAndTrim("  .env.local  ,  .env.test  ")
	if len(got) != 2 {
		t.Fatalf("expected 2 parts, got %d: %v", len(got), got)
	}
	if got[0] != ".env.local" {
		t.Errorf("expected .env.local, got %q", got[0])
	}
	if got[1] != ".env.test" {
		t.Errorf("expected .env.test, got %q", got[1])
	}
}

func TestHasIssues_NoIssues(t *testing.T) {
	reports := []comparator.Report{
		{Target: "a", Missing: nil, Extra: nil},
		{Target: "b", Missing: []string{}, Extra: []string{}},
	}
	if hasIssues(reports) {
		t.Error("expected no issues")
	}
}

func TestHasIssues_WithMissing(t *testing.T) {
	reports := []comparator.Report{
		{Target: "a", Missing: []string{"KEY_A"}, Extra: nil},
	}
	if !hasIssues(reports) {
		t.Error("expected issues due to missing keys")
	}
}

func TestHasIssues_WithExtra(t *testing.T) {
	reports := []comparator.Report{
		{Target: "b", Missing: nil, Extra: []string{"EXTRA_KEY"}},
	}
	if !hasIssues(reports) {
		t.Error("expected issues due to extra keys")
	}
}
