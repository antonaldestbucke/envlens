package config

import (
	"testing"
)

func TestParse_ValidArgs(t *testing.T) {
	cfg, err := Parse([]string{
		"-ref", ".env",
		"-targets", "prod.env,staging.env",
		"-format", "json",
		"-strict",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Reference != ".env" {
		t.Errorf("expected reference .env, got %q", cfg.Reference)
	}
	if len(cfg.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(cfg.Targets))
	}
	if cfg.Format != "json" {
		t.Errorf("expected format json, got %q", cfg.Format)
	}
	if !cfg.Strict {
		t.Error("expected strict to be true")
	}
}

func TestParse_DefaultFormat(t *testing.T) {
	cfg, err := Parse([]string{"-ref", ".env", "-targets", "prod.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "text" {
		t.Errorf("expected default format text, got %q", cfg.Format)
	}
	if cfg.Strict {
		t.Error("expected strict to be false by default")
	}
	if cfg.Quiet {
		t.Error("expected quiet to be false by default")
	}
}

func TestParse_MissingReference(t *testing.T) {
	_, err := Parse([]string{"-targets", "prod.env"})
	if err != ErrNoReference {
		t.Errorf("expected ErrNoReference, got %v", err)
	}
}

func TestParse_MissingTargets(t *testing.T) {
	_, err := Parse([]string{"-ref", ".env"})
	if err != ErrNoTargets {
		t.Errorf("expected ErrNoTargets, got %v", err)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := Parse([]string{"-ref", ".env", "-targets", "prod.env", "-format", "xml"})
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestParse_TargetsTrimmed(t *testing.T) {
	cfg, err := Parse([]string{"-ref", ".env", "-targets", " prod.env , staging.env "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets[0] != "prod.env" {
		t.Errorf("expected trimmed target, got %q", cfg.Targets[0])
	}
	if cfg.Targets[1] != "staging.env" {
		t.Errorf("expected trimmed target, got %q", cfg.Targets[1])
	}
}

func TestParse_QuietFlag(t *testing.T) {
	cfg, err := Parse([]string{"-ref", ".env", "-targets", "prod.env", "-quiet"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Quiet {
		t.Error("expected quiet to be true")
	}
}
