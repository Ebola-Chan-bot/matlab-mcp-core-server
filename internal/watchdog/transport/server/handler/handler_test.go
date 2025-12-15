// Copyright 2025 The MathWorks, Inc.

package handler_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
	handlermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/server/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	// Act
	h := handler.New(mockLoggerFactory, mockProcessHandler)

	// Assert
	assert.NotNil(t, h)
}

func TestHandler_HandleProcessToKill_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	expectedPID := 12345
	request := messages.ProcessToKillRequest{
		PID: expectedPID,
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)

	// Act
	response, err := h.HandleProcessToKill(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ProcessToKillResponse{}, response)
}

func TestHandler_HandleShutdown_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	request := messages.ShutdownRequest{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)

	// Act
	response, err := h.HandleShutdown(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ShutdownResponse{}, response)
}

func TestHandler_HandleShutdown_CallsRegisteredFunctions(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	request := messages.ShutdownRequest{}
	functionCalled := false

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)
	h.RegisterShutdownFunction(func() {
		functionCalled = true
	})

	// Act
	response, err := h.HandleShutdown(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ShutdownResponse{}, response)
	assert.True(t, functionCalled, "Registered shutdown function should have been called")
}

func TestHandler_HandleShutdown_CallsMultipleRegisteredFunctions(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	request := messages.ShutdownRequest{}
	callOrder := []int{}
	expectedCallOrder := []int{1, 2}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)
	h.RegisterShutdownFunction(func() {
		callOrder = append(callOrder, expectedCallOrder[0])
	})
	h.RegisterShutdownFunction(func() {
		callOrder = append(callOrder, expectedCallOrder[1])
	})

	// Act
	response, err := h.HandleShutdown(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ShutdownResponse{}, response)
	assert.Equal(t, expectedCallOrder, callOrder, "Registered shutdown functions should be called in order")
}

func TestHandler_TerminateAllProcesses_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	expectedPID := 12345

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID).
		Return(nil).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)

	_, err := h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID})
	require.NoError(t, err)

	// Act
	h.TerminateAllProcesses()

	// Assert
}

func TestHandler_TerminateAllProcesses_MultiplePIDs(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	expectedPID1 := 12345
	expectedPID2 := 67890

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID1).
		Return(nil).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID2).
		Return(nil).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)

	_, err := h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID1})
	require.NoError(t, err)
	_, err = h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID2})
	require.NoError(t, err)

	// Act
	h.TerminateAllProcesses()

	// Assert
}

func TestHandler_TerminateAllProcesses_KillProcessError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	expectedPID1 := 12345
	expectedPID2 := 67890
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID1).
		Return(expectedError).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID2).
		Return(nil).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)

	_, err := h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID1})
	require.NoError(t, err)
	_, err = h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID2})
	require.NoError(t, err)

	// Act
	h.TerminateAllProcesses()

	// Assert
}

func TestHandler_TerminateAllProcesses_NoPIDs(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	h := handler.New(mockLoggerFactory, mockProcessHandler)

	// Act
	h.TerminateAllProcesses()

	// Assert
}
