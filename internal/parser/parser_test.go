package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Keys["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env.Keys["APP_ENV"])
	}
	if env.Keys["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", env.Keys["DB_HOST"])
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\n\nFOO=bar\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(env.Keys))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Keys["SECRET"] != "my secret value" {
		t.Errorf("expected unquoted value, got %q", env.Keys["SECRET"])
	}
	if env.Keys["TOKEN"] != "abc123" {
		t.Errorf("expected unquoted value, got %q", env.Keys["TOKEN"])
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # default port\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Keys["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", env.Keys["PORT"])
	}
}

func TestParse_InvalidSyntax(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid syntax, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
