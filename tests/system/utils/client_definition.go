// Copyright 2025 The MathWorks, Inc.

package utils

import "github.com/modelcontextprotocol/go-sdk/mcp"

func GetMCPCLientImplementation() *mcp.Implementation {
	// Those values don't matter for the system tests, but are required to construct an MCP client.
	return &mcp.Implementation{
		Name:    "matlab-mcp-client",
		Version: "v0.0.1",
	}
}
