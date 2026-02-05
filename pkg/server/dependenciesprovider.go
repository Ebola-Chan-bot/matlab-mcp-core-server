// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

type DependenciesProviderResources struct {
	logger logger.Logger
}

func newDependenciesProviderResources(resources definition.DependenciesProviderResources) DependenciesProviderResources {
	return DependenciesProviderResources{
		logger: newLoggerAdaptor(resources.Logger),
	}
}

func (r DependenciesProviderResources) Logger() logger.Logger {
	return r.logger
}

type DependenciesProvider[Dependencies any] func(dependenciesProviderResources DependenciesProviderResources) (Dependencies, error)

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
