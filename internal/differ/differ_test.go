package differ_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

func TestCompute_AllMatch(t *testing.T) {
	ref := []string{"A", "B", "C"}
	target := []string{"A", "B", "C"}

	res := differ.Compute("prod", ref, target)

	if res.Target != "prod" {
		t.Fatalf("expected target 'prod', got %q", res.Target)
	}
	if len(res.Lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(res.Lines))
	}
	for _, l := range res.Lines {
		if l.Kind != differ.KindContext {
			t.Errorf("expected KindContext for key %q, got %v", l.Key, l.Kind)
		}
	}
}

func TestCompute_MissingKeys(t *testing.T) {
	ref := []string{"A", "B", "C"}
	target := []string{"A"}

	res := differ.Compute("staging", ref, target)

	counts := kindCounts(res.Lines)
	if counts[differ.KindContext] != 1 {
		t.Errorf("expected 1 context, got %d", counts[differ.KindContext])
	}
	if counts[differ.KindMissing] != 2 {
		t.Errorf("expected 2 missing, got %d", counts[differ.KindMissing])
	}
}

func TestCompute_ExtraKeys(t *testing.T) {
	ref := []string{"A"}
	target := []string{"A", "X", "Y"}

	res := differ.Compute("dev", ref, target)

	counts := kindCounts(res.Lines)
	if counts[differ.KindExtra] != 2 {
		t.Errorf("expected 2 extra, got %d", counts[differ.KindExtra])
	}
}

func TestCompute_EmptyRef(t *testing.T) {
	res := differ.Compute("dev", []string{}, []string{"A", "B"})

	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
	for _, l := range res.Lines {
		if l.Kind != differ.KindExtra {
			t.Errorf("expected KindExtra for %q", l.Key)
		}
	}
}

func TestCompute_EmptyTarget(t *testing.T) {
	res := differ.Compute("dev", []string{"A", "B"}, []string{})

	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
	for _, l := range res.Lines {
		if l.Kind != differ.KindMissing {
			t.Errorf("expected KindMissing for %q", l.Key)
		}
	}
}

func kindCounts(lines []differ.DiffLine) map[differ.LineKind]int {
	m := make(map[differ.LineKind]int)
	for _, l := range lines {
		m[l.Kind]++
	}
	return m
}
