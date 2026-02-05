// Copyright 2026 The MathWorks, Inc.

package main

import "github.com/matlab/matlab-mcp-core-server/pkg/server"

type ToolsProviderResources interface{}

func ToolsProvider(resources ToolsProviderResources) []server.Tool {
	return []server.Tool{
		NewGreetTool(),
		NewGreetStructuredTool(),
	}
}
