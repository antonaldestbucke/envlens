package config

import (
	"fmt"
	"os"
)

// ValidationError describes a file-access problem found during validation.
type ValidationError struct {
	Field string
	Path  string
	Cause error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: cannot access %q: %v", e.Field, e.Path, e.Cause)
}

// Validate checks that all file paths referenced in cfg actually exist and
// are readable. It returns the first error encountered.
func Validate(cfg *Config) error {
	if err := checkReadable("reference", cfg.Reference); err != nil {
		return err
	}
	for _, t := range cfg.Targets {
		if err := checkReadable("target", t); err != nil {
			return err
		}
	}
	return nil
}

func checkReadable(field, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return &ValidationError{Field: field, Path: path, Cause: err}
	}
	f.Close()
	return nil
}
