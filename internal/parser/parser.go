package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvFile represents a parsed .env file with its keys and source path.
type EnvFile struct {
	Path string
	Keys map[string]string
}

// Parse reads a .env file and returns an EnvFile with all parsed key-value pairs.
// It skips blank lines and comments (lines starting with '#').
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: failed to open file %q: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path: path,
		Keys: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip inline comments
		if idx := strings.Index(line, " #"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("parser: invalid syntax at %s:%d: %q", path, lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		if key == "" {
			return nil, fmt.Errorf("parser: empty key at %s:%d", path, lineNum)
		}

		env.Keys[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: error reading file %q: %w", path, err)
	}

	return env, nil
}
