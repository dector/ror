package env

import "fmt"

var version = "0.1.0-00"

// This variable is meant to be set at compile time using ldflags, e.g.:
// go build -ldflags "-X 'github.com/dector/ror/internal/env.commit=$(git describe --tags --always --dirty)'"
var commit string

// This variable is meant to be set at compile time using ldflags, e.g.:
// go build -ldflags "-X 'github.com/dector/ror/internal/env.time=$(date +%Y-%m-%dT%H:%M:%SZ)'"
var time string

// GetShortVersion returns the short version string.
func GetShortVersion() string {
	return fmt.Sprintf("v%s", version)
}

func GetDefaultVersion() string {
	return fmt.Sprintf("v%s:%s", version, commit)
}

// GetLongVersion returns the full version string (includes commit and build time).
func GetLongVersion() string {
	return fmt.Sprintf("version: %s\ncommit: %s\ntime: %s", version, commit, time)
}
