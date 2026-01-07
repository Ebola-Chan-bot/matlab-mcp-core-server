// Copyright 2025-2026 The MathWorks, Inc.

package sdk

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Config interface {
	Version() string
}

type Factory struct {
	config Config
}

func NewFactory(
	config Config,
) *Factory {
	return &Factory{
		config: config,
	}
}

func (f *Factory) NewServer(name string, instructions string) *mcp.Server {
	impl := &mcp.Implementation{
		Name:    name,
		Version: f.config.Version(),
	}
	options := &mcp.ServerOptions{
		Instructions: instructions,
	}
	return mcp.NewServer(impl, options)
}
