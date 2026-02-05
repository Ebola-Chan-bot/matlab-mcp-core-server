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

// ServerWithCustomToolsTestSuite tests SDK custom tools functionalities.
type ServerWithCustomToolsTestSuite struct {
	suite.Suite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithCustomToolsTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithCustomTools(s.T())
}

func TestServerWithCustomToolsTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithCustomToolsTestSuite))
}

func (s *ServerWithCustomToolsTestSuite) TestSDK_CustomTools_UnstructuredContentOutput_HappyPath() {
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
	expectedTextOutput := "Hello " + name

	// Call the tool
	result, err := session.CallTool(s.T().Context(), "greet", map[string]any{"name": name})
	s.Require().NoError(err, "should call tool successfully")

	textContent, err := session.GetTextContent(result)
	s.Require().NoError(err, "should get text content")
	s.Require().Equal(expectedTextOutput, textContent, "should return greeting message")

	// Check the logger is wired correctly
	logFile := filepath.Join(logFolder, "server-"+instanceID+".log")

	ctx, cancel := context.WithTimeout(s.T().Context(), 2*time.Second) // Timeout for the logs to write to disk
	defer cancel()

	_, err = retry.Retry(ctx, func() (struct{}, bool, error) {
		logContent, err := os.ReadFile(logFile) //nolint:gosec // G304: logFile is a controlled test path
		if err != nil {
			return struct{}{}, false, err
		}

		foundLogEntry := strings.Contains(string(logContent), "Greeting user: "+name)

		return struct{}{}, foundLogEntry, nil
	}, retry.NewLinearRetryStrategy(200*time.Millisecond))

	s.Require().NoError(err)
}

func (s *ServerWithCustomToolsTestSuite) TestSDK_CustomTools_StructuredContentOutput_HappyPath() {
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
	expectedResponse := "Hello " + name

	// Call the tool
	result, err := session.CallTool(s.T().Context(), "greet-structured", map[string]any{"name": "World"})
	s.Require().NoError(err, "should call tool successfully")

	var output struct {
		Response string `json:"response"`
	}
	s.Require().NoError(session.UnmarshalStructuredContent(result, &output), "should unmarshal structured content")
	s.Require().Equal(expectedResponse, output.Response, "should return greeting message")

	// Check the logger is wired correctly
	logFile := filepath.Join(logFolder, "server-"+instanceID+".log")

	ctx, cancel := context.WithTimeout(s.T().Context(), 2*time.Second) // Timeout for the logs to write to disk
	defer cancel()

	_, err = retry.Retry(ctx, func() (struct{}, bool, error) {
		logContent, err := os.ReadFile(logFile) //nolint:gosec // G304: logFile is a controlled test path
		if err != nil {
			return struct{}{}, false, err
		}

		foundLogEntry := strings.Contains(string(logContent), "Greeting user structurally: "+name)

		return struct{}{}, foundLogEntry, nil
	}, retry.NewLinearRetryStrategy(200*time.Millisecond))

	s.Require().NoError(err)
}
