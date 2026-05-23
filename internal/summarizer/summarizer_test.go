package summarizer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/comparator"
	"github.com/yourorg/envlens/internal/summarizer"
)

func makeReports(specs []struct {
	name    string
	missing []string
	extra   []string
}) []comparator.Report {
	var out []comparator.Report
	for _, s := range specs {
		out = append(out, comparator.Report{
			Target:  s.name,
			Missing: s.missing,
			Extra:   s.extra,
		})
	}
	return out
}

func TestSummarize_AllClean(t *testing.T) {
	reports := makeReports([]struct {
		name    string
		missing []string
		extra   []string
	}{
		{"prod", nil, nil},
		{"staging", nil, nil},
	})
	res := summarizer.Summarize(reports)
	if res.TotalTargets != 2 {
		t.Errorf("expected 2 targets, got %d", res.TotalTargets)
	}
	if res.CleanTargets != 2 {
		t.Errorf("expected 2 clean targets, got %d", res.CleanTargets)
	}
	if res.HasIssues() {
		t.Error("expected no issues")
	}
}

func TestSummarize_SomeDirty(t *testing.T) {
	reports := makeReports([]struct {
		name    string
		missing []string
		extra   []string
	}{
		{"prod", []string{"SECRET"}, nil},
		{"staging", nil, []string{"DEBUG"}},
		{"local", nil, nil},
	})
	res := summarizer.Summarize(reports)
	if res.DirtyTargets != 2 {
		t.Errorf("expected 2 dirty targets, got %d", res.DirtyTargets)
	}
	if res.TotalMissing != 1 {
		t.Errorf("expected 1 missing key, got %d", res.TotalMissing)
	}
	if res.TotalExtra != 1 {
		t.Errorf("expected 1 extra key, got %d", res.TotalExtra)
	}
	if !res.HasIssues() {
		t.Error("expected issues to be detected")
	}
}

func TestSummarize_Empty(t *testing.T) {
	res := summarizer.Summarize(nil)
	if res.TotalTargets != 0 {
		t.Errorf("expected 0 targets, got %d", res.TotalTargets)
	}
	if res.HasIssues() {
		t.Error("empty summary should have no issues")
	}
}

func TestSummarize_PerTargetDetails(t *testing.T) {
	reports := makeReports([]struct {
		name    string
		missing []string
		extra   []string
	}{
		{"prod", []string{"A", "B"}, []string{"C"}},
	})
	res := summarizer.Summarize(reports)
	if len(res.PerTarget) != 1 {
		t.Fatalf("expected 1 per-target entry, got %d", len(res.PerTarget))
	}
	pt := res.PerTarget[0]
	if pt.Name != "prod" {
		t.Errorf("expected target name 'prod', got %q", pt.Name)
	}
	if pt.MissingCount != 2 {
		t.Errorf("expected 2 missing, got %d", pt.MissingCount)
	}
	if pt.ExtraCount != 1 {
		t.Errorf("expected 1 extra, got %d", pt.ExtraCount)
	}
	if pt.Clean {
		t.Error("target should not be clean")
	}
}
