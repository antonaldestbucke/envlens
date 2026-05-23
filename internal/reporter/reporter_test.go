package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/comparator"
)

func makeReport() Report {
	return Report{
		Reference: ".env.example",
		Results: []comparator.Result{
			{
				Target:  ".env.production",
				Missing: []string{"DB_PASSWORD", "SECRET_KEY"},
				Extra:   []string{"OLD_FLAG"},
			},
			{
				Target:  ".env.staging",
				Missing: []string{},
				Extra:   []string{},
			},
		},
	}
}

func TestWriteText_ContainsReference(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeReport(), FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), ".env.example") {
		t.Error("expected reference file name in output")
	}
}

func TestWriteText_ShowsMissingKeys(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, makeReport(), FormatText)
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in text output")
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Error("expected SECRET_KEY in text output")
	}
}

func TestWriteText_ShowsExtraKeys(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, makeReport(), FormatText)
	if !strings.Contains(buf.String(), "OLD_FLAG") {
		t.Error("expected OLD_FLAG in text output")
	}
}

func TestWriteText_CleanTarget(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, makeReport(), FormatText)
	if !strings.Contains(buf.String(), "[OK] .env.staging") {
		t.Error("expected [OK] marker for clean target")
	}
}

func TestWriteJSON_ContainsTarget(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeReport(), FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, ".env.production") {
		t.Error("expected .env.production in JSON output")
	}
}

func TestWriteJSON_ContainsMissingKey(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, makeReport(), FormatJSON)
	if !strings.Contains(buf.String(), "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in JSON output")
	}
}

func TestWriteJSON_EmptyArrays(t *testing.T) {
	var buf bytes.Buffer
	report := Report{
		Reference: ".env.example",
		Results: []comparator.Result{
			{Target: ".env.local", Missing: []string{}, Extra: []string{}},
		},
	}
	Write(&buf, report, FormatJSON)
	if !strings.Contains(buf.String(), "[]") {
		t.Error("expected empty JSON arrays for clean result")
	}
}
