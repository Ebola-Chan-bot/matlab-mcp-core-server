// Copyright 2025 The MathWorks, Inc.

package socket_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
	socketmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/socket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockDirectory := &socketmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	// Act
	factory := socket.NewFactory(mockDirectory)

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_Socket_HappyPath(t *testing.T) {
	// Arrange
	mockDirectory := &socketmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	baseDir := filepath.Join("tmp", "watchdog")
	id := "abc123"

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(id).
		Once()

	factory := socket.NewFactory(mockDirectory)

	// Act
	socketInstance, err := factory.Socket()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, socketInstance)
}

func TestFactory_Socket_Singleton(t *testing.T) {
	// Arrange
	mockDirectory := &socketmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	baseDir := filepath.Join("tmp", "watchdog")
	id := "abc123"

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(id).
		Once()

	factory := socket.NewFactory(mockDirectory)

	// Act
	firstSocketInstance, firstErr := factory.Socket()
	secondSocketInstance, secondErr := factory.Socket()

	// Assert
	assert.NoError(t, firstErr)
	assert.NoError(t, secondErr)
	assert.Equal(t, firstSocketInstance, secondSocketInstance)
}

func TestFactory_Socket_ReturnCachedError(t *testing.T) {
	// Arrange
	mockDirectory := &socketmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	baseDir := filepath.Join("tmp", "s_super_long_path_that_causes_an_error_for_the_socket_to_be_created_due_to_os_limits")
	id := "abc123"

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(id).
		Once()

	factory := socket.NewFactory(mockDirectory)

	// Act
	firstSocketInstance, firstErr := factory.Socket()
	secondSocketInstance, secondErr := factory.Socket()

	// Assert
	assert.Equal(t, firstErr, secondErr)
	assert.Nil(t, firstSocketInstance)
	assert.Nil(t, secondSocketInstance)
}

func TestSocket_HappyPath(t *testing.T) {
	// Arrange
	mockDirectory := &socketmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	baseDir := filepath.Join("tmp", "watchdog")
	id := "abc123"

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(id).
		Once()

	// Act
	socketInstance, err := socket.NewSocket(mockDirectory)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, socketInstance)
}

func TestSocket_PathTooLong(t *testing.T) {
	// Arrange
	mockDirectory := &socketmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	baseDir := filepath.Join("tmp", "s_super_long_path_that_causes_an_error_for_the_socket_to_be_created_due_to_os_limits")
	id := "abc123"

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(id).
		Once()

	// Act
	socketInstance, err := socket.NewSocket(mockDirectory)

	// Assert
	require.ErrorIs(t, err, socket.ErrSocketPathTooLong)
	assert.Nil(t, socketInstance)
}

func TestSocket_Path_HappyPath(t *testing.T) {
	// Arrange
	mockDirectory := &socketmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	baseDir := filepath.Join("tmp", "watchdog")
	id := "abc123"
	expectedPath := filepath.Join(baseDir, "watchdog-"+id+".sock")

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(id).
		Once()

	socketInstance, err := socket.NewSocket(mockDirectory)
	require.NoError(t, err)

	// Act
	path := socketInstance.Path()

	// Assert
	assert.Equal(t, expectedPath, path)
}
