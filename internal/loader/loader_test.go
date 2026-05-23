package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_BasicFile(t *testing.T) {
	p := writeTempEnv(t, ".env", "FOO=bar\nBAZ=qux\n")
	ef, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Target != "default" {
		t.Errorf("expected target 'default', got %q", ef.Target)
	}
	if ef.Keys["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", ef.Keys["FOO"])
	}
}

func TestLoad_TargetName(t *testing.T) {
	p := writeTempEnv(t, ".env.staging", "KEY=value\n")
	ef, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Target != "staging" {
		t.Errorf("expected target 'staging', got %q", ef.Target)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadAll_MultipleFiles(t *testing.T) {
	p1 := writeTempEnv(t, ".env", "A=1\nB=2\n")
	p2 := writeTempEnv(t, ".env.production", "A=1\nB=2\nC=3\n")

	files, err := LoadAll([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[1].Target != "production" {
		t.Errorf("expected 'production', got %q", files[1].Target)
	}
}

func TestLoadAll_EmptyPaths(t *testing.T) {
	_, err := LoadAll([]string{})
	if err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}

func TestTargetName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{".env", "default"},
		{".env.production", "production"},
		{".env.staging", "staging"},
		{"env.local", "local"},
		{"custom.env", "custom.env"},
	}
	for _, tc := range cases {
		got := targetName(tc.input)
		if got != tc.want {
			t.Errorf("targetName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
