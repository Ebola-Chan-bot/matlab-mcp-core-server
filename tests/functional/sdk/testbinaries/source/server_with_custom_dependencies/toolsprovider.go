// Copyright 2026 The MathWorks, Inc.

package main

import (
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
)

type ToolsProviderResources interface {
	Logger() logger.Logger
	Dependencies() Dependencies
}

func ToolsProvider(resources ToolsProviderResources) []server.Tool {
	resources.Logger().Info("Creating Tools")

	dataService := resources.Dependencies().DataService

	return []server.Tool{
		NewGreetTool(dataService),
	}
}
