// Package config provides CLI flag parsing for envlens.
//
// It defines the Config struct that captures all runtime options
// and a Parse function that reads from an args slice, making it
// straightforward to unit-test without touching os.Args directly.
//
// Typical usage:
//
//	cfg, err := config.Parse(os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
package config
