// Copyright 2025 The MathWorks, Inc.

package system_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/system/config"
	"github.com/matlab/matlab-mcp-core-server/tests/system/execsuite"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSystem(t *testing.T) {
	matlabMCPServerBinariesPath := getMATLABMCPCoreServerPath(t)

	suite.Run(t, execsuite.NewSuite(matlabMCPServerBinariesPath))
}

const matlabMCPCoreServerBinaryPathEnvironmentVariable = "MATLAB_MCP_CORE_SERVER_BINARY_PATH"

// getMATLABMCPCoreServerPath retrieves the absolute path of the matlab-mcp-sever binaries.
// It requires `make` to have been run, or at a minimum `make build` to generate the binaries.
//
// For CI, we allow overrides, by using MATLAB_MCP_CORE_SERVER_BINARY_PATH environment variable
func getMATLABMCPCoreServerPath(t *testing.T) string {
	path := filepath.Join(
		"..",
		"..",
		".bin",
		config.OSDescriptor,
		config.MATLABMCPCoreServerBinariesFilename,
	)

	if value := os.Getenv(matlabMCPCoreServerBinaryPathEnvironmentVariable); value != "" {
		path = value
	}

	path, err := filepath.Abs(path)

	require.NoError(t, err, "Failed to get absolute path")
	require.NotEmpty(t, path, "matlab-mcp-core-server binary path cannot be empty")
	require.FileExists(t, path, "matlab-mcp-core-server binary does not exist")

	return path
}
