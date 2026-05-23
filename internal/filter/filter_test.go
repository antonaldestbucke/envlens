package filter_test

import (
	"testing"

	"github.com/envlens/internal/comparator"
	"github.com/envlens/internal/filter"
)

func makeReports() []comparator.Report {
	return []comparator.Report{
		{
			Target:      "staging",
			MissingKeys: []string{"DB_HOST", "DB_PORT", "AWS_KEY", "APP_NAME"},
			ExtraKeys:   []string{"OLD_DB_URL", "REDIS_URL"},
		},
		{
			Target:      "production",
			MissingKeys: []string{"AWS_SECRET", "APP_ENV"},
			ExtraKeys:   []string{"DEBUG"},
		},
	}
}

func TestApply_NoOptions_ReturnsUnchanged(t *testing.T) {
	reports := makeReports()
	got := filter.Apply(reports, filter.Options{})
	if len(got) != len(reports) {
		t.Fatalf("expected %d reports, got %d", len(reports), len(got))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	reports := makeReports()
	opts := filter.Options{Prefixes: []string{"DB_"}}
	got := filter.Apply(reports, opts)

	if len(got[0].MissingKeys) != 2 {
		t.Errorf("expected 2 missing DB_ keys in staging, got %v", got[0].MissingKeys)
	}
	for _, k := range got[0].MissingKeys {
		if k != "DB_HOST" && k != "DB_PORT" {
			t.Errorf("unexpected key %q passed prefix filter", k)
		}
	}
	if len(got[0].ExtraKeys) != 0 {
		t.Errorf("expected no extra DB_ keys, got %v", got[0].ExtraKeys)
	}
}

func TestApply_SuffixFilter(t *testing.T) {
	reports := makeReports()
	opts := filter.Options{Suffixes: []string{"_SECRET"}}
	got := filter.Apply(reports, opts)

	if len(got[1].MissingKeys) != 1 || got[1].MissingKeys[0] != "AWS_SECRET" {
		t.Errorf("expected [AWS_SECRET] in production missing keys, got %v", got[1].MissingKeys)
	}
}

func TestApply_PrefixAndSuffix(t *testing.T) {
	reports := makeReports()
	opts := filter.Options{
		Prefixes: []string{"AWS_"},
		Suffixes: []string{"_URL"},
	}
	got := filter.Apply(reports, opts)

	// staging: missing AWS_KEY; extra OLD_DB_URL, REDIS_URL
	if len(got[0].MissingKeys) != 1 || got[0].MissingKeys[0] != "AWS_KEY" {
		t.Errorf("unexpected missing keys: %v", got[0].MissingKeys)
	}
	if len(got[0].ExtraKeys) != 2 {
		t.Errorf("expected 2 extra keys (_URL), got %v", got[0].ExtraKeys)
	}
}

func TestApply_NoMatch_EmptyKeys(t *testing.T) {
	reports := makeReports()
	opts := filter.Options{Prefixes: []string{"NONEXISTENT_"}}
	got := filter.Apply(reports, opts)

	for _, r := range got {
		if len(r.MissingKeys) != 0 || len(r.ExtraKeys) != 0 {
			t.Errorf("expected empty keys for target %s, got missing=%v extra=%v",
				r.Target, r.MissingKeys, r.ExtraKeys)
		}
	}
}
