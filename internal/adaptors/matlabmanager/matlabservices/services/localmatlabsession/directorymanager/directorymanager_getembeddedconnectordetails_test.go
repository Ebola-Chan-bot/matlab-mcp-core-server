// Copyright 2025 The MathWorks, Inc.

package directorymanager_test

import (
	"os"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectoryManager_GetEmbeddedConnectorDetails_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	directoryManager.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := directoryManager.SecurePortFile()
	certificateFile := directoryManager.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := directoryManager.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectoryManager_GetEmbeddedConnectorDetails_WaitsForSecurePortFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	directoryManager.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := directoryManager.SecurePortFile()
	certificateFile := directoryManager.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := directoryManager.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectoryManager_GetEmbeddedConnectorDetails_WaitsForCertificateFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	directoryManager.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := directoryManager.SecurePortFile()
	certificateFile := directoryManager.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := directoryManager.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectoryManager_GetEmbeddedConnectorDetails_WaitsForNotEmptyPortFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	directoryManager.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := directoryManager.SecurePortFile()
	certificateFile := directoryManager.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(""), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := directoryManager.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectoryManager_GetEmbeddedConnectorDetails_WaitsForNotEmptyCertificateFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	directoryManager.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := directoryManager.SecurePortFile()
	certificateFile := directoryManager.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return([]byte(""), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := directoryManager.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectoryManager_GetEmbeddedConnectorDetails_TimesoutWaitingForFilesToExists(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	directoryManager.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := directoryManager.SecurePortFile()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(nil, os.ErrNotExist)

	// Act
	port, certificate, err := directoryManager.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectoryManager_GetEmbeddedConnectorDetails_TimesoutWaitingForFileContent(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	directoryManager.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := directoryManager.SecurePortFile()
	certificateFile := directoryManager.CertificateFile()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(""), nil) // Will be called multiple times in wait loop

	// Act
	port, certificate, err := directoryManager.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}
