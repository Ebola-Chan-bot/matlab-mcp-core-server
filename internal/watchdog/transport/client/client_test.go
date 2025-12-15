// Copyright 2025 The MathWorks, Inc.

package client_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/client"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
	httpmocks "github.com/matlab/matlab-mcp-core-server/mocks/utils/httpclientfactory"
	clientmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_Connect_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHttpClient := &httpmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, nil).
		Once()

	mockHTTPClientFactory.EXPECT().
		NewClientOverUDS(expectedSocketPath).
		Return(mockHttpClient).
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)

	// Act
	err := clientInstance.Connect(expectedSocketPath)

	// Assert
	require.NoError(t, err)
}

func TestClient_Connect_Timeout(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, os.ErrNotExist).
		Maybe()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)
	clientInstance.SetSocketWaitTimeout(50 * time.Millisecond)
	clientInstance.SetSocketRetryInterval(10 * time.Millisecond)

	// Act
	err := clientInstance.Connect(expectedSocketPath)

	// Assert
	require.ErrorIs(t, err, client.ErrTimeoutWaitingForSocketFile)
}

func TestClient_Connect_StatError(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, expectedError).
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)

	// Act
	err := clientInstance.Connect(expectedSocketPath)

	// Assert
	require.ErrorIs(t, err, client.ErrSocketFileInaccessible)
}

func TestClient_SendProcessPID_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHttpClient := &httpmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")
	expectedPID := 12345

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, nil).
		Once()

	mockHTTPClientFactory.EXPECT().
		NewClientOverUDS(expectedSocketPath).
		Return(mockHttpClient).
		Once()

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(req *http.Request) bool {
			if req.Method != "POST" {
				return false
			}
			if req.URL.Path != messages.ProcessToKillPath {
				return false
			}
			if req.Header.Get("Content-Type") != "application/json" {
				return false
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				return false
			}
			var reqBody messages.ProcessToKillRequest
			if err := json.Unmarshal(body, &reqBody); err != nil {
				return false
			}
			return reqBody.PID == expectedPID
		})).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
		}, nil).
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)
	err := clientInstance.Connect(expectedSocketPath)
	require.NoError(t, err)

	// Act
	_, err = clientInstance.SendProcessPID(expectedPID)

	// Assert
	require.NoError(t, err)
}

func TestClient_SendProcessPID_ClientNotConnected(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedPID := 12345

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)

	// Act
	_, err := clientInstance.SendProcessPID(expectedPID)

	// Assert
	require.ErrorIs(t, err, client.ErrClientNotConnected)
}

func TestClient_SendProcessPID_HTTPError(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHttpClient := &httpmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")
	expectedPID := 12345
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, nil).
		Once()

	mockHTTPClientFactory.EXPECT().
		NewClientOverUDS(expectedSocketPath).
		Return(mockHttpClient).
		Once()

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(nil, expectedError).
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)
	err := clientInstance.Connect(expectedSocketPath)
	require.NoError(t, err)

	// Act
	_, err = clientInstance.SendProcessPID(expectedPID)

	// Assert
	require.ErrorIs(t, err, client.ErrHTTP)
}

func TestClient_SendProcessPID_UnexpectedStatus(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHttpClient := &httpmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")
	expectedPID := 12345

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, nil).
		Once()

	mockHTTPClientFactory.EXPECT().
		NewClientOverUDS(expectedSocketPath).
		Return(mockHttpClient).
		Once()

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     "500 Internal Server Error",
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil).
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)
	err := clientInstance.Connect(expectedSocketPath)
	require.NoError(t, err)

	// Act
	_, err = clientInstance.SendProcessPID(expectedPID)

	// Assert
	require.ErrorIs(t, err, client.ErrHTTP)
}

func TestClient_SendStop_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHttpClient := &httpmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, nil).
		Once()

	mockHTTPClientFactory.EXPECT().
		NewClientOverUDS(expectedSocketPath).
		Return(mockHttpClient).
		Once()

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(req *http.Request) bool {
			if req.Method != "POST" {
				return false
			}
			if req.URL.Path != messages.ShutdownPath {
				return false
			}
			if req.Header.Get("Content-Type") != "application/json" {
				return false
			}
			return true
		})).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
		}, nil).
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)
	err := clientInstance.Connect(expectedSocketPath)
	require.NoError(t, err)

	// Act
	_, err = clientInstance.SendStop()

	// Assert
	require.NoError(t, err)
}

func TestClient_SendStop_ClientNotConnected(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)

	// Act
	_, err := clientInstance.SendStop()

	// Assert
	require.ErrorIs(t, err, client.ErrClientNotConnected)
}

func TestClient_SendStop_HTTPError(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHttpClient := &httpmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, nil).
		Once()

	mockHTTPClientFactory.EXPECT().
		NewClientOverUDS(expectedSocketPath).
		Return(mockHttpClient).
		Once()

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(nil, expectedError).
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)
	err := clientInstance.Connect(expectedSocketPath)
	require.NoError(t, err)

	// Act
	_, err = clientInstance.SendStop()

	// Assert
	require.ErrorIs(t, err, client.ErrHTTP)
}

func TestClient_Close_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHttpClient := &httpmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSocketPath := filepath.Join("tmp", "test.sock")

	mockOSLayer.EXPECT().
		Stat(expectedSocketPath).
		Return(nil, nil).
		Once()

	mockHTTPClientFactory.EXPECT().
		NewClientOverUDS(expectedSocketPath).
		Return(mockHttpClient).
		Once()

	mockHttpClient.EXPECT().
		CloseIdleConnections().
		Return().
		Once()

	clientInstance := client.NewClient(
		mockOSLayer,
		mockHTTPClientFactory,
		mockLogger,
	)
	err := clientInstance.Connect(expectedSocketPath)
	require.NoError(t, err)

	// Act
	err = clientInstance.Close()

	// Assert
	require.NoError(t, err)
}
