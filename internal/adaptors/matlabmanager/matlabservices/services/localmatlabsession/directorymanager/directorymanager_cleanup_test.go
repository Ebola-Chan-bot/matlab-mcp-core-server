// Copyright 2025 The MathWorks, Inc.

package directorymanager_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectoryManager_Cleanup_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetCleanupTimeout(100 * time.Millisecond)
	directoryManager.SetCleanupRetry(10 * time.Millisecond)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	// Act
	err := directoryManager.Cleanup()

	// Assert
	require.NoError(t, err)
}

func TestDirectoryManager_Cleanup_WaitsForRemoveAllToPass(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetCleanupTimeout(100 * time.Millisecond)
	directoryManager.SetCleanupRetry(10 * time.Millisecond)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(assert.AnError).
		Once()

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	// Act
	err := directoryManager.Cleanup()

	// Assert
	require.NoError(t, err)
}

func TestDirectoryManager_Cleanup_Timesout(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	directoryManager := directorymanager.NewDirectoryManager(sessionDir, mockOSLayer)
	directoryManager.SetCleanupTimeout(100 * time.Millisecond)
	directoryManager.SetCleanupRetry(10 * time.Millisecond)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(assert.AnError) // Will be called many times with retry

	// Act
	err := directoryManager.Cleanup()

	// Assert
	require.Error(t, err)
}
