// Copyright 2025 The MathWorks, Inc.

package server_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpserverfactory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServerFactory_NewServerOverUDS_HappyPath(t *testing.T) {
	// Arrange
	factory := newServerFactory()

	socketPath := filepath.Join(t.TempDir(), "test.sock")

	expectedStatusCode := http.StatusOK
	expectedFirstBody := "first hello world"
	expectedSecondBody := "second hello world"

	handlers := map[string]http.HandlerFunc{
		"GET /first": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(expectedFirstBody))
		},
		"POST /second": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(expectedSecondBody))
		},
	}

	server, err := factory.NewServerOverUDS(handlers)
	require.NoError(t, err)

	serverStopped := make(chan error, 1)
	go func() {
		serverStopped <- server.Serve(socketPath)
	}()
	defer func() {
		require.NoError(t, server.Shutdown(t.Context()))
	}()

	waitForSocketFile(t, socketPath)

	client := newUDSClient(socketPath)

	// Act & Assert
	req, err := http.NewRequest("GET", "http://unix/first", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, expectedFirstBody, string(body))

	req, err = http.NewRequest("POST", "http://unix/second", nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, expectedSecondBody, string(body))
}

func newServerFactory() *httpserverfactory.HTTPServerFactory {
	osLayer := osfacade.New()
	return httpserverfactory.New(osLayer)
}

func newUDSClient(socketPath string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
}

func waitForSocketFile(t *testing.T, socketPath string) {
	t.Helper()

	timeout := time.After(1 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		_, err := os.Stat(socketPath)
		if err == nil {
			return
		}

		require.ErrorIs(t, err, os.ErrNotExist)

		select {
		case <-timeout:
			t.Fatalf("Failed to wait for socket file: %v", err)
		case <-tick:
		}
	}
}
