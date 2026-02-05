// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/suite"
)

// ServerWithCustomDependenciesTestSuite tests SDK custom dependencies functionalities.
type ServerWithCustomDependenciesTestSuite struct {
	suite.Suite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithCustomDependenciesTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithCustomDependencies(s.T())
}

func TestServerWithCustomDependenciesTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithCustomDependenciesTestSuite))
}

func (s *ServerWithCustomDependenciesTestSuite) TestSDK_CustomDependencies_ToolUsesDependency_HappyPath() {
	// Connect to a session
	logFolder, err := os.MkdirTemp("", "server_session") // Can't use s.T().Tempdir() because too long for socket path
	s.Require().NoError(err)
	defer s.Require().NoError(os.RemoveAll(logFolder))

	instanceID := "123"

	client := mcpclient.NewClient(s.T().Context(), s.serverDetails.BinaryLocation(), nil,
		"--log-level=debug",
		"--log-folder="+logFolder,
		"--server-instance-id="+instanceID,
	)

	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	defer func() {
		s.Require().NoError(session.Close(), "closing session should not error")
	}()

	name := "World"
	expectedTextOutput := "Service Hello " + name

	// Call the tool
	result, err := session.CallTool(s.T().Context(), "greet", map[string]any{"name": name})
	s.Require().NoError(err, "should call tool successfully")

	textContent, err := session.GetTextContent(result)
	s.Require().NoError(err, "should get text content")
	s.Require().Equal(expectedTextOutput, textContent, "should return greeting message using dependency")

	// Check the logger is wired correctly
	logFile := filepath.Join(logFolder, "server-"+instanceID+".log")

	ctx, cancel := context.WithTimeout(s.T().Context(), 2*time.Second) // Timeout for the logs to write to disk
	defer cancel()

	_, err = retry.Retry(ctx, func() (struct{}, bool, error) {
		logContent, err := os.ReadFile(logFile) //nolint:gosec // G304: logFile is a controlled test path
		if err != nil {
			return struct{}{}, false, err
		}

		foundDependenciesProviderLog := strings.Contains(string(logContent), "Creating Dependencies")
		foundToolsProviderLog := strings.Contains(string(logContent), "Creating Tools")

		return struct{}{}, foundDependenciesProviderLog && foundToolsProviderLog, nil
	}, retry.NewLinearRetryStrategy(200*time.Millisecond))

	s.Require().NoError(err)
}
