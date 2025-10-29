// Copyright 2025 The MathWorks, Inc.

package execsuite

import (
	"context"
	"os/exec"

	"github.com/matlab/matlab-mcp-core-server/tests/system/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Test_ServerStartsAndClientCanConnect tests that we can start the binaries, and that a client can connect.
func (s *Suite) Test_ServerStartsAndClientCanConnect() {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectionTimeout)
	defer cancel()

	transport := &mcp.CommandTransport{
		Command: exec.Command(s.matlabMCPServerBinariesPath), //nolint:gosec // Using generated path for testing
	}

	// Act
	client := mcp.NewClient(utils.GetMCPCLientImplementation(), nil)
	session, err := client.Connect(ctx, transport, nil)

	// Assert
	s.Require().NoError(err, "Client connection should succeed")
	session.Close() //nolint:errcheck,gosec // Ignore error on session close as it doesn't impact test outcome
}
