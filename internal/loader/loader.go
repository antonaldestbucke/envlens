package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/envlens/internal/parser"
)

// EnvFile represents a parsed .env file with its associated target name.
type EnvFile struct {
	Target string
	Path   string
	Keys   map[string]string
}

// LoadAll loads and parses all .env files from the given paths.
// The target name is derived from the file name (e.g., ".env.production" -> "production").
func LoadAll(paths []string) ([]EnvFile, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no env files provided")
	}

	var files []EnvFile
	for _, p := range paths {
		ef, err := Load(p)
		if err != nil {
			return nil, err
		}
		files = append(files, ef)
	}
	return files, nil
}

// Load parses a single .env file and returns an EnvFile.
func Load(path string) (EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return EnvFile{}, fmt.Errorf("loader: cannot open %q: %w", path, err)
	}
	defer f.Close()

	keys, err := parser.Parse(f)
	if err != nil {
		return EnvFile{}, fmt.Errorf("loader: cannot parse %q: %w", path, err)
	}

	return EnvFile{
		Target: targetName(path),
		Path:   path,
		Keys:   keys,
	}, nil
}

// targetName extracts a human-readable target label from a file path.
// ".env" -> "default", ".env.production" -> "production".
func targetName(path string) string {
	base := filepath.Base(path)
	switch base {
	case ".env", "env":
		return "default"
	}
	// Strip leading ".env." or "env."
	for _, prefix := range []string{".env.", "env."} {
		if len(base) > len(prefix) && base[:len(prefix)] == prefix {
			return base[len(prefix):]
		}
	}
	return base
}
