// Copyright 2025 The MathWorks, Inc.

package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerConfig interface {
	Version() string
}

func NewMCPSDKServer(config ServerConfig) *mcp.Server {
	impl := &mcp.Implementation{
		Name:    name,
		Version: config.Version(),
	}
	options := &mcp.ServerOptions{
		Instructions: instructions,
	}
	return mcp.NewServer(impl, options)
}
