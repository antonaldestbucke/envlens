package config

import (
	"errors"
	"flag"
	"strings"
)

// Config holds all CLI-parsed configuration for an envlens run.
type Config struct {
	Reference string
	Targets   []string
	Format    string
	Strict    bool
	Quiet     bool
}

// ErrNoReference is returned when no reference file is provided.
var ErrNoReference = errors.New("reference file is required")

// ErrNoTargets is returned when no target files are provided.
var ErrNoTargets = errors.New("at least one target file is required")

// Parse reads flags from the provided args slice and returns a Config.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envlens", flag.ContinueOnError)

	ref := fs.String("ref", "", "path to the reference .env file")
	targets := fs.String("targets", "", "comma-separated list of target .env files")
	format := fs.String("format", "text", "output format: text or json")
	strict := fs.Bool("strict", false, "treat extra keys in targets as errors")
	quiet := fs.Bool("quiet", false, "suppress output, only use exit code")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *ref == "" {
		return nil, ErrNoReference
	}

	parsedTargets := splitAndTrim(*targets, ",")
	if len(parsedTargets) == 0 {
		return nil, ErrNoTargets
	}

	if *format != "text" && *format != "json" {
		return nil, errors.New("format must be \"text\" or \"json\"")
	}

	return &Config{
		Reference: *ref,
		Targets:   parsedTargets,
		Format:    *format,
		Strict:    *strict,
		Quiet:     *quiet,
	}, nil
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
