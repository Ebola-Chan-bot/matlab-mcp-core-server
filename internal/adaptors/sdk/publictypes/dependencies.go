// Copyright 2026 The MathWorks, Inc.

package publictypes

type DependenciesProviderResources interface { //nolint:iface // Semantically different interfaces
	Logger() Logger
	Config() Config
}

type DependenciesProvider[Dependencies any] func(DependenciesProviderResources) (Dependencies, Error)
