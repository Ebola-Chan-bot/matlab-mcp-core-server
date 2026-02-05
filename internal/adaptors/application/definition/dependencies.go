// Copyright 2026 The MathWorks, Inc.

package definition

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type DependenciesProviderResources struct {
	Logger entities.Logger
}

type DependenciesProvider func(resources DependenciesProviderResources) (any, error)

func NewDependenciesProviderResources(logger entities.Logger) DependenciesProviderResources {
	return DependenciesProviderResources{
		Logger: logger,
	}
}
