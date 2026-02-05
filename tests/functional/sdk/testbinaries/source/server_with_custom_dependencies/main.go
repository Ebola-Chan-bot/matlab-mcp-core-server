// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/pkg/server"
)

func main() {
	serverDefinition := server.Definition[Dependencies]{
		Name:         "server-with-custom-dependencies",
		Title:        "Server With Custom Dependencies",
		Instructions: "This is a test server with custom dependencies",

		DependenciesProvider: func(dependenciesProviderResources server.DependenciesProviderResources) (Dependencies, error) {
			return DependenciesProvider(dependenciesProviderResources)
		},
		ToolsProvider: func(toolsProviderResources server.ToolsProviderResources[Dependencies]) []server.Tool {
			return ToolsProvider(toolsProviderResources)
		},
	}
	serverInstance := server.New(serverDefinition)

	exitCode := serverInstance.StartAndWaitForCompletion(context.Background())

	os.Exit(exitCode)
}
