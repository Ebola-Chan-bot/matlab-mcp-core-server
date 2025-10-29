// Copyright 2025 The MathWorks, Inc.

package execsuite

import (
	"context"
	"os/exec"

	"github.com/matlab/matlab-mcp-core-server/tests/system/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Test_ListAvailableMATLABsTool tests that the ListAvailableMATLABs tool works as expected.
func (s *Suite) Test_ListAvailableMATLABsTool() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectionTimeout)
	defer cancel()

	transport := &mcp.CommandTransport{
		Command: exec.Command( //nolint:gosec // Using generated path for testing
			s.matlabMCPServerBinariesPath,
			"--use-single-matlab-session=false",
		),
	}

	client := mcp.NewClient(utils.GetMCPCLientImplementation(), nil)
	session, err := client.Connect(ctx, transport, nil)
	s.Require().NoError(err, "Client connection should succeed")
	defer session.Close() //nolint:errcheck // Ignore error on session close as it doesn't impact test outcome

	// Act
	params := &mcp.CallToolParams{
		Name:      "list_available_matlabs",
		Arguments: map[string]any{},
	}
	results, err := session.CallTool(ctx, params)

	// Assert
	s.Require().NoError(err, "Tool call should succeed")
	s.NotNil(results, "Results should not be nil")

	// Here, we need  further qualification on `results.StructuredContent`.
	// This requires to have controlled information about MATLAB installations in TC job hosts.
	// Maybe that requires installing MATLAB with `mpm`?
}
