// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/pkg/config"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

type DependenciesProviderResources struct {
	logger logger.Logger
	config config.Config
}

func newDependenciesProviderResources(resources definition.DependenciesProviderResources) DependenciesProviderResources {
	return DependenciesProviderResources{
		logger: newLoggerAdaptor(resources.Logger),
		config: newConfigAdaptor(resources.Config, resources.MessageCatalog),
	}
}

func (r DependenciesProviderResources) Logger() logger.Logger {
	return r.logger
}

func (r DependenciesProviderResources) Config() config.Config {
	return r.config
}

type DependenciesProvider[Dependencies any] func(dependenciesProviderResources DependenciesProviderResources) (Dependencies, i18n.Error)

func (p DependenciesProvider[Dependencies]) toInternal() definition.DependenciesProvider {
	return func(resources definition.DependenciesProviderResources) (any, error) {
		if p == nil {
			return nil, nil
		}

		dependencies, err := p(newDependenciesProviderResources(resources))
		if err != nil {
			return nil, err
		}

		return dependencies, nil
	}
}
