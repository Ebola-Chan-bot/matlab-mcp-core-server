// Copyright 2025 The MathWorks, Inc.

package httpclientfactory_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpclientfactory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCertificate(t *testing.T) []byte {
	t.Helper()

	template := &x509.Certificate{SerialNumber: big.NewInt(1)}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	return cert
}

func TestNew_HappyPath(t *testing.T) {
	// Arrange

	// Act
	factory := httpclientfactory.New()

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestHTTPClientFactory_NewClientForSelfSignedTLSServer_HappyPath(t *testing.T) {
	// Arrange
	certPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCertificate(t),
	})

	factory := httpclientfactory.New()

	// Act
	client, err := factory.NewClientForSelfSignedTLSServer(certPEMBytes)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, client, "Client should not be nil")
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

func TestHTTPClientFactory_NewClientOverUDS_HappyPath(t *testing.T) {
	// Arrange
	factory := httpclientfactory.New()

	// Act
	client := factory.NewClientOverUDS("")

	// Assert
	require.NotNil(t, client, "Client should not be nil")
}
