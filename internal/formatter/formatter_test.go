package formatter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/formatter"
)

func TestFormatDiff_MissingKeysMarkedMinus(t *testing.T) {
	ref := []string{"A", "B", "C"}
	missing := []string{"B"}
	extra := []string{}

	lines := formatter.FormatDiff(ref, missing, extra)

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[1].Symbol != "-" || lines[1].Key != "B" {
		t.Errorf("expected line[1] to be '- B', got '%s %s'", lines[1].Symbol, lines[1].Key)
	}
}

func TestFormatDiff_ExtraKeysMarkedPlus(t *testing.T) {
	ref := []string{"A"}
	missing := []string{}
	extra := []string{"Z"}

	lines := formatter.FormatDiff(ref, missing, extra)

	last := lines[len(lines)-1]
	if last.Symbol != "+" || last.Key != "Z" {
		t.Errorf("expected last line to be '+ Z', got '%s %s'", last.Symbol, last.Key)
	}
}

func TestFormatDiff_CleanTargetAllSpace(t *testing.T) {
	ref := []string{"A", "B"}
	lines := formatter.FormatDiff(ref, nil, nil)

	for _, l := range lines {
		if l.Symbol != " " {
			t.Errorf("expected space symbol, got %q", l.Symbol)
		}
	}
}

func TestRenderDiff_ContainsSymbols(t *testing.T) {
	lines := []formatter.DiffLine{
		{Symbol: " ", Key: "HOST"},
		{Symbol: "-", Key: "PORT"},
		{Symbol: "+", Key: "EXTRA"},
	}

	output := formatter.RenderDiff(lines)

	if !strings.Contains(output, "- PORT") {
		t.Error("expected output to contain '- PORT'")
	}
	if !strings.Contains(output, "+ EXTRA") {
		t.Error("expected output to contain '+ EXTRA'")
	}
	if !strings.Contains(output, "  HOST") {
		t.Error("expected output to contain '  HOST'")
	}
}

func TestTruncate_WithinLimit(t *testing.T) {
	keys := []string{"A", "B", "C"}
	result := formatter.Truncate(keys, 5)
	if len(result) != 3 {
		t.Errorf("expected 3, got %d", len(result))
	}
}

func TestTruncate_ExceedsLimit(t *testing.T) {
	keys := []string{"A", "B", "C", "D", "E"}
	result := formatter.Truncate(keys, 3)
	if len(result) != 4 {
		t.Fatalf("expected 4 items (3 + summary), got %d", len(result))
	}
	if !strings.HasPrefix(result[3], "... and") {
		t.Errorf("expected summary line, got %q", result[3])
	}
}
