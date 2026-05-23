package sorter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/comparator"
	"github.com/yourorg/envlens/internal/sorter"
)

func makeReport(target string, missing, extra []string) comparator.Report {
	return comparator.Report{
		Target:      target,
		MissingKeys: missing,
		ExtraKeys:   extra,
	}
}

func TestSortReport_Alpha(t *testing.T) {
	r := makeReport("prod", []string{"ZEBRA", "APPLE", "MANGO"}, []string{"ZOO", "BAR"})
	got := sorter.SortReport(r, sorter.OrderAlpha)

	wantMissing := []string{"APPLE", "MANGO", "ZEBRA"}
	wantExtra := []string{"BAR", "ZOO"}

	for i, k := range got.MissingKeys {
		if k != wantMissing[i] {
			t.Errorf("MissingKeys[%d] = %q, want %q", i, k, wantMissing[i])
		}
	}
	for i, k := range got.ExtraKeys {
		if k != wantExtra[i] {
			t.Errorf("ExtraKeys[%d] = %q, want %q", i, k, wantExtra[i])
		}
	}
}

func TestSortReport_AlphaDesc(t *testing.T) {
	r := makeReport("staging", []string{"APPLE", "MANGO", "ZEBRA"}, []string{"BAR", "ZOO"})
	got := sorter.SortReport(r, sorter.OrderAlphaDesc)

	wantMissing := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range got.MissingKeys {
		if k != wantMissing[i] {
			t.Errorf("MissingKeys[%d] = %q, want %q", i, k, wantMissing[i])
		}
	}
}

func TestSortReport_None_PreservesOrder(t *testing.T) {
	orig := []string{"ZEBRA", "APPLE", "MANGO"}
	r := makeReport("dev", orig, nil)
	got := sorter.SortReport(r, sorter.OrderNone)

	for i, k := range got.MissingKeys {
		if k != orig[i] {
			t.Errorf("MissingKeys[%d] = %q, want %q", i, k, orig[i])
		}
	}
}

func TestSortReport_DoesNotMutateOriginal(t *testing.T) {
	orig := []string{"ZEBRA", "APPLE"}
	r := makeReport("prod", orig, nil)
	_ = sorter.SortReport(r, sorter.OrderAlpha)

	if r.MissingKeys[0] != "ZEBRA" {
		t.Errorf("original report was mutated: got %q, want %q", r.MissingKeys[0], "ZEBRA")
	}
}

func TestSortReports_AppliesOrderToAll(t *testing.T) {
	reports := []comparator.Report{
		makeReport("prod", []string{"Z", "A"}, nil),
		makeReport("staging", []string{"M", "B"}, nil),
	}
	got := sorter.SortReports(reports, sorter.OrderAlpha)

	if got[0].MissingKeys[0] != "A" {
		t.Errorf("report[0] not sorted: got %q", got[0].MissingKeys[0])
	}
	if got[1].MissingKeys[0] != "B" {
		t.Errorf("report[1] not sorted: got %q", got[1].MissingKeys[0])
	}
}
