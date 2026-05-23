// Package loader provides utilities for discovering and loading .env files
// from the filesystem. It wraps the parser package to produce EnvFile values
// that carry both the parsed key-value map and a human-readable target label
// (e.g. "production", "staging", "default").
//
// Typical usage:
//
//	files, err := loader.LoadAll([]string{".env", ".env.production"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	// files[0].Target == "default"
//	// files[1].Target == "production"
package loader
