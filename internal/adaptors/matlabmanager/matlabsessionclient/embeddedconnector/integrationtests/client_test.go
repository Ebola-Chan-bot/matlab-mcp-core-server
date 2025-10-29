// Copyright 2025 The MathWorks, Inc.

package embeddedconnector_integration_test

import (
	_ "embed"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/stretchr/testify/assert"
)

func startTestServer(t *testing.T, handler func(responseWriter http.ResponseWriter, request *http.Request)) embeddedconnector.ConnectionDetails {
	t.Helper()

	const expectedAPIKey = "test-api-key"

	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "POST", request.Method)
		assert.Equal(t, "/messageservice/json/secure", request.URL.Path)
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
		assert.Equal(t, expectedAPIKey, request.Header.Get("mwapikey"))

		handler(responseWriter, request)
	}))

	addressDetails := strings.Split(server.Listener.Addr().String(), ":")

	serverCert := server.Certificate()
	certPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCert.Raw,
	})

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host:           addressDetails[0],
		Port:           addressDetails[1],
		APIKey:         expectedAPIKey,
		CertificatePEM: certPEMBytes,
	}
	return connectionDetails
}
