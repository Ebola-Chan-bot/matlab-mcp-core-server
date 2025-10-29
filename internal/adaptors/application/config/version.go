// Copyright 2025 The MathWorks, Inc.

package config

// These variables are set via ldflags during build
var (
	// version is the semantic Version of the application
	version = "(devel)" // DO NOT USE `version = unsetVersion`, it won't work as expected
)

const (
	// unsetVersion is used as a default version if no explicit version is set during build
	unsetVersion = "(devel)"
)
