// Package redactor masks sensitive values in parsed .env maps
// before they are passed to reporters or loggers.
package redactor

import "strings"

// DefaultSensitiveSuffixes contains common key suffixes considered sensitive.
var DefaultSensitiveSuffixes = []string{
	"_SECRET",
	"_PASSWORD",
	"_PASSWD",
	"_TOKEN",
	"_KEY",
	"_PRIVATE",
	"_CREDENTIAL",
	"_CREDENTIALS",
	"_DSN",
	"_API_KEY",
}

// Mask is the string used to replace sensitive values.
const Mask = "***REDACTED***"

// Options controls redaction behaviour.
type Options struct {
	// ExtraSuffixes appends caller-supplied suffixes to the defaults.
	ExtraSuffixes []string
	// ReplaceAll replaces every value regardless of key name.
	ReplaceAll bool
}

// Redact returns a shallow copy of env with sensitive values replaced by Mask.
// The original map is never modified.
func Redact(env map[string]string, opts Options) map[string]string {
	suffixes := buildSuffixes(opts)
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.ReplaceAll || isSensitive(k, suffixes) {
			out[k] = Mask
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitive reports whether key is considered sensitive using default suffixes.
func IsSensitive(key string) bool {
	return isSensitive(key, buildSuffixes(Options{}))
}

func isSensitive(key string, suffixes []string) bool {
	upper := strings.ToUpper(key)
	for _, s := range suffixes {
		if strings.HasSuffix(upper, s) {
			return true
		}
	}
	return false
}

func buildSuffixes(opts Options) []string {
	suffixes := make([]string, len(DefaultSensitiveSuffixes))
	copy(suffixes, DefaultSensitiveSuffixes)
	for _, s := range opts.ExtraSuffixes {
		suffixes = append(suffixes, strings.ToUpper(s))
	}
	return suffixes
}
