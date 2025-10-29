// Copyright 2025 The MathWorks, Inc.

package directory_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-12345"

	mockOSLayer.EXPECT().
		MkdirTemp("", mock.AnythingOfType("string")).
		Return(expectedLogDir, nil).
		Once()

	// Act
	directoryInstance, err := directory.New(mockOSLayer)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, directoryInstance, "Directory instance should not be nil")
}

func TestNew_MkdirTempError(t *testing.T) {
	// Arrange
	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		MkdirTemp("", mock.AnythingOfType("string")).
		Return("", expectedError).
		Once()

	// Act
	directoryInstance, err := directory.New(mockOSLayer)

	// Assert
	require.ErrorIs(t, err, expectedError, "New should return the error from MkdirTemp")
	assert.Nil(t, directoryInstance, "Directory instance should be nil when error occurs")
}

func TestBaseDir_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-67890"

	mockOSLayer.EXPECT().
		MkdirTemp("", mock.AnythingOfType("string")).
		Return(expectedLogDir, nil).
		Once()

	directoryInstance, err := directory.New(mockOSLayer)
	require.NoError(t, err)

	// Act
	baseDir := directoryInstance.BaseDir()

	// Assert
	assert.Equal(t, expectedLogDir, baseDir, "BaseDir should return the expected log directory")
}

func TestMkdirTemp_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	baseLogDir := "/tmp/matlab-mcp-core-server-11111"
	pattern := "test-pattern-"
	expectedTempDir := "/tmp/matlab-mcp-core-server-11111/test-pattern-22222"

	mockOSLayer.EXPECT().
		MkdirTemp("", mock.AnythingOfType("string")).
		Return(baseLogDir, nil).
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp(baseLogDir, pattern).
		Return(expectedTempDir, nil).
		Once()

	directoryInstance, err := directory.New(mockOSLayer)
	require.NoError(t, err)

	// Act
	tempDir, err := directoryInstance.MkdirTemp(pattern)

	// Assert
	require.NoError(t, err, "MkdirTemp should not return an error")
	assert.Equal(t, expectedTempDir, tempDir, "MkdirTemp should return the expected temp directory")
}

func TestMkdirTemp_MkdirTempError(t *testing.T) {
	// Arrange
	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	baseLogDir := "/tmp/matlab-mcp-core-server-33333"
	pattern := "test-pattern-"
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		MkdirTemp("", mock.AnythingOfType("string")).
		Return(baseLogDir, nil).
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp(baseLogDir, pattern).
		Return("", expectedError).
		Once()

	directoryInstance, err := directory.New(mockOSLayer)
	require.NoError(t, err)

	// Act
	tempDir, err := directoryInstance.MkdirTemp(pattern)

	// Assert
	require.ErrorIs(t, err, expectedError, "MkdirTemp should return the error from OSLayer.MkdirTemp")
	assert.Empty(t, tempDir, "MkdirTemp should return empty string when error occurs")
}
