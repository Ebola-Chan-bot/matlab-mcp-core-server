// Copyright 2025 The MathWorks, Inc.

package watchdog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog"
	transportmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport"
	socketmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/socket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandler := &mocks.MockServerHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	// Act
	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandler,
		mockServerFactory,
		mockSocketFactory,
	)

	// Assert
	assert.NotNil(t, watchdogInstance, "Watchdog instance should not be nil")
}

func TestWatchdog_StartAndWaitForCompletion_GracefulShutdown(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandler := &mocks.MockServerHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockServer := &transportmocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	serverStarted := make(chan struct{})
	expectedParentPID := 1234
	shutdownFuncC := make(chan func())
	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(socketPath).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(mockServer, nil).
		Once()

	mockServer.EXPECT().
		Start(socketPath).
		Run(func(_ string) {
			close(serverStarted)
		}).
		Return(nil).
		Once()

	mockServer.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockServerHandler.EXPECT().
		RegisterShutdownFunction(mock.AnythingOfType("func()")).
		Run(func(fn func()) {
			shutdownFuncC <- fn
		}).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockServerHandler.EXPECT().
		TerminateAllProcesses().
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandler,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	<-serverStarted
	shutdownFcn := <-shutdownFuncC

	shutdownFcn()

	// Assert
	require.NoError(t, <-errC)
}

func TestWatchdog_StartAndWaitForCompletion_ParentProcessTermination(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandler := &mocks.MockServerHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockServer := &transportmocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	serverStarted := make(chan struct{})
	expectedParentPID := 1234
	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(socketPath).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(mockServer, nil).
		Once()

	mockServer.EXPECT().
		Start(socketPath).
		Run(func(_ string) {
			close(serverStarted)
		}).
		Return(nil).
		Once()

	mockServer.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockServerHandler.EXPECT().
		RegisterShutdownFunction(mock.AnythingOfType("func()")).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockServerHandler.EXPECT().
		TerminateAllProcesses().
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandler,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	<-serverStarted
	close(parentTerminationC)

	// Assert
	require.NoError(t, <-errC)
}

func TestWatchdog_StartAndWaitForCompletion_OSSignalInterrupt(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandler := &mocks.MockServerHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockServer := &transportmocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	serverStarted := make(chan struct{})
	expectedParentPID := 1234
	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(socketPath).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(mockServer, nil).
		Once()

	mockServer.EXPECT().
		Start(socketPath).
		Run(func(_ string) {
			close(serverStarted)
		}).
		Return(nil).
		Once()

	mockServer.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockServerHandler.EXPECT().
		RegisterShutdownFunction(mock.AnythingOfType("func()")).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockServerHandler.EXPECT().
		TerminateAllProcesses().
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandler,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	<-serverStarted
	interruptSignalC <- os.Interrupt

	// Assert
	require.NoError(t, <-errC)
}

func TestWatchdog_StartAndWaitForCompletion_SocketFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandler := &mocks.MockServerHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandler,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestWatchdog_StartAndWaitForCompletion_ServerFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandler := &mocks.MockServerHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandler,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError)
}
