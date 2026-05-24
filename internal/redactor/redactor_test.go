package redactor_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/redactor"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":      "myapp",
		"DB_PASSWORD":   "s3cr3t",
		"API_KEY":       "abc123",
		"AUTH_TOKEN":    "tok_xyz",
		"REDIS_DSN":     "redis://localhost",
		"LOG_LEVEL":     "info",
		"STRIPE_SECRET": "sk_live_xxx",
	}
}

func TestRedact_LeavesNonSensitiveUnchanged(t *testing.T) {
	out := redactor.Redact(baseEnv(), redactor.Options{})
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", out["APP_NAME"])
	}
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", out["LOG_LEVEL"])
	}
}

func TestRedact_MasksSensitiveValues(t *testing.T) {
	out := redactor.Redact(baseEnv(), redactor.Options{})
	sensitiveKeys := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "REDIS_DSN", "STRIPE_SECRET"}
	for _, k := range sensitiveKeys {
		if out[k] != redactor.Mask {
			t.Errorf("expected %s to be redacted, got %q", k, out[k])
		}
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	original := env["DB_PASSWORD"]
	redactor.Redact(env, redactor.Options{})
	if env["DB_PASSWORD"] != original {
		t.Error("original map was mutated")
	}
}

func TestRedact_ReplaceAll(t *testing.T) {
	out := redactor.Redact(baseEnv(), redactor.Options{ReplaceAll: true})
	for k, v := range out {
		if v != redactor.Mask {
			t.Errorf("expected all values redacted, but %s=%q", k, v)
		}
	}
}

func TestRedact_ExtraSuffixes(t *testing.T) {
	env := map[string]string{
		"GITHUB_PAT": "ghp_abc",
		"PUBLIC_URL": "https://example.com",
	}
	out := redactor.Redact(env, redactor.Options{ExtraSuffixes: []string{"_PAT"}})
	if out["GITHUB_PAT"] != redactor.Mask {
		t.Errorf("expected GITHUB_PAT redacted, got %q", out["GITHUB_PAT"])
	}
	if out["PUBLIC_URL"] != "https://example.com" {
		t.Errorf("expected PUBLIC_URL unchanged, got %q", out["PUBLIC_URL"])
	}
}

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"AUTH_TOKEN", true},
		{"APP_NAME", false},
		{"LOG_LEVEL", false},
		{"PRIVATE_KEY", true},
	}
	for _, tc := range cases {
		got := redactor.IsSensitive(tc.key)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
