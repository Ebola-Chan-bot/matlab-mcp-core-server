// Copyright 2026 The MathWorks, Inc.

package dependenciesproviderresources

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/config"
	publictypes "github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type LoggerFactory interface {
	New(logger entities.Logger) publictypes.Logger
}

type ConfigFactory interface {
	New(
		internalConfig config.InternalConfig,
		internalMessageCatalog config.InternalMessageCatalog,
	) publictypes.Config
}

type Factory struct {
	loggerFactory LoggerFactory
	configFactory ConfigFactory
}

func NewFactory(
	loggerFactory LoggerFactory,
	configFactory ConfigFactory,
) *Factory {
	return &Factory{
		loggerFactory: loggerFactory,
		configFactory: configFactory,
	}
}

func (f *Factory) New(
	internal definition.DependenciesProviderResources,
) publictypes.DependenciesProviderResources {
	return &dependenciesProviderResourcesAdaptor{
		logger: f.loggerFactory.New(internal.Logger),
		config: f.configFactory.New(internal.Config, internal.MessageCatalog),
	}
}

type dependenciesProviderResourcesAdaptor struct {
	logger publictypes.Logger
	config publictypes.Config
}

func (r *dependenciesProviderResourcesAdaptor) Logger() publictypes.Logger {
	return r.logger
}

func (r *dependenciesProviderResourcesAdaptor) Config() publictypes.Config {
	return r.config
}
