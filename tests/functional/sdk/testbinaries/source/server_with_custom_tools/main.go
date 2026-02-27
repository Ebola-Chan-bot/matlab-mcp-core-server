// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/pkg/server"
)

func main() {
	serverDefinition := server.Definition[any]{
		Name:         "server-with-custom-tools",
		Title:        "Server With Custom Tools",
		Instructions: "This is a test server with custom tools",

		ToolsProvider: func(toolsProviderResources server.ToolsProviderResources[any]) []server.Tool {
			return ToolsProvider(toolsProviderResources)
		},
	}
	serverInstance := server.New(serverDefinition)

	exitCode := serverInstance.StartAndWaitForCompletion(context.Background())

	os.Exit(exitCode)
}
