// Copyright 2025-2026 The MathWorks, Inc.

package sdk

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Factory struct {
	configFactory ConfigFactory
}

func NewFactory(
	configFactory ConfigFactory,
) *Factory {
	return &Factory{
		configFactory: configFactory,
	}
}

func (f *Factory) NewServer(name string, instructions string) (*mcp.Server, messages.Error) {
	config, err := f.configFactory.Config()
	if err != nil {
		return nil, err
	}

	impl := &mcp.Implementation{
		Name:    name,
		Version: config.Version(),
	}
	options := &mcp.ServerOptions{
		Instructions: instructions,
	}

	return mcp.NewServer(impl, options), nil
}
