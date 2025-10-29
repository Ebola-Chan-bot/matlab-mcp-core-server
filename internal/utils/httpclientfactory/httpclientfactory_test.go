// Copyright 2025 The MathWorks, Inc.

package httpclientfactory_test

import (
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpclientfactory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange

	// Act
	factory := httpclientfactory.New()

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestHTTPClientFactory_NewClientForSelfSignedTLSServer_HappyPath(t *testing.T) {
	// Arrange
	expectedStatusCode := http.StatusOK

	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(expectedStatusCode)
	}))

	serverCert := server.Certificate()
	certPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCert.Raw,
	})

	factory := httpclientfactory.New()

	// Act
	client, err := factory.NewClientForSelfSignedTLSServer(certPEMBytes)

	// Assert
	require.NoError(t, err)

	// Act + Assert to check the client is functional
	request, err := http.NewRequest("GET", "https://"+server.Listener.Addr().String(), nil)
	require.NoError(t, err)
	response, err := client.Do(request)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, response.Body.Close())
	})
	assert.Equal(t, expectedStatusCode, response.StatusCode)
}

func TestHTTPClientFactory_NewClientForSelfSignedTLSServer_InvalidCert(t *testing.T) {
	// Arrange
	factory := httpclientfactory.New()

	// Act
	client, err := factory.NewClientForSelfSignedTLSServer([]byte("invalid cert"))

	// Assert
	require.Error(t, err)
	assert.Nil(t, client)
}
