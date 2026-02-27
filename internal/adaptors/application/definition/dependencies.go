// Copyright 2026 The MathWorks, Inc.

package definition

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type DependenciesProviderResources struct {
	Logger         entities.Logger
	Config         config.GenericConfig
	MessageCatalog MessageCatalog
}

type DependenciesProvider func(resources DependenciesProviderResources) (any, error)

func NewDependenciesProviderResources(
	logger entities.Logger,
	config config.GenericConfig,
	messageCatalog MessageCatalog,
) DependenciesProviderResources {
	return DependenciesProviderResources{
		Logger:         logger,
		Config:         config,
		MessageCatalog: messageCatalog,
	}
}
