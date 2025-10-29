// Copyright 2025 The MathWorks, Inc.

package processhandler_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/processhandler"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	processhandlermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/processhandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSWrapper := &processhandlermocks.MockOSWrapper{}
	defer mockOSWrapper.AssertExpectations(t)

	// Act
	processHandlerInstance := processhandler.New(mockOSWrapper)

	// Assert
	assert.NotNil(t, processHandlerInstance, "ProcessHandler instance should not be nil")
}

func TestProcessHandler_WatchProcessAndGetTerminationChan_HappyPath(t *testing.T) {
	// Arrange
	mockOSWrapper := &processhandlermocks.MockOSWrapper{}
	defer mockOSWrapper.AssertExpectations(t)

	mockProcess := &osfacademocks.MockProcess{}
	defer mockProcess.AssertExpectations(t)

	processPid := 1234

	mockOSWrapper.EXPECT().
		WaitForProcessToComplete(processPid).
		Return().
		Once()

	processHandlerInstance := processhandler.New(mockOSWrapper)

	// Act
	terminationChan := processHandlerInstance.WatchProcessAndGetTerminationChan(processPid)

	// Assert
	select {
	case <-terminationChan:
		// Expected - process terminated
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Expected termination channel to be closed when process disappears")
	}
}

func TestProcessHandler_KillProcess_HappyPath(t *testing.T) {
	// Arrange
	mockOSWrapper := &processhandlermocks.MockOSWrapper{}
	defer mockOSWrapper.AssertExpectations(t)

	mockProcess := &osfacademocks.MockProcess{}
	defer mockProcess.AssertExpectations(t)

	processPid := 1234

	mockOSWrapper.EXPECT().
		FindProcess(processPid).
		Return(mockProcess).
		Once()

	mockProcess.EXPECT().
		Kill().
		Return(nil).
		Once()

	processHandlerInstance := processhandler.New(mockOSWrapper)

	// Act
	err := processHandlerInstance.KillProcess(processPid)

	// Assert
	require.NoError(t, err, "KillProcess should not return an error when process exists and is killed successfully")
}

func TestProcessHandler_KillProcess_ProcessExistsButKillFails(t *testing.T) {
	// Arrange
	mockOSWrapper := &processhandlermocks.MockOSWrapper{}
	defer mockOSWrapper.AssertExpectations(t)

	mockProcess := &osfacademocks.MockProcess{}
	defer mockProcess.AssertExpectations(t)

	processPid := 1234
	expectedError := assert.AnError

	mockOSWrapper.EXPECT().
		FindProcess(processPid).
		Return(mockProcess).
		Once()

	mockProcess.EXPECT().
		Kill().
		Return(expectedError).
		Once()

	processHandlerInstance := processhandler.New(mockOSWrapper)

	// Act
	err := processHandlerInstance.KillProcess(processPid)

	// Assert
	assert.ErrorIs(t, err, expectedError, "KillProcess should return the error from Kill method")
}

func TestProcessHandler_KillProcess_ProcessDoesNotExist(t *testing.T) {
	// Arrange
	mockOSWrapper := &processhandlermocks.MockOSWrapper{}
	defer mockOSWrapper.AssertExpectations(t)

	processPid := 1234

	mockOSWrapper.EXPECT().
		FindProcess(processPid).
		Return(nil).
		Once()

	processHandlerInstance := processhandler.New(mockOSWrapper)

	// Act
	err := processHandlerInstance.KillProcess(processPid)

	// Assert
	require.NoError(t, err, "KillProcess should not return an error when process doesn't exist")
}
