// Copyright 2025 The MathWorks, Inc.

package transport_test

import (
	"io"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFactory_NewClient_HappyPath(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	blockUntilStdoutIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdoutIsClosed
	}()

	blockUntilStderrIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStderrIsClosed
	}()

	mockStdout.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) { close(blockUntilStdoutIsClosed) }).
		Once()

	mockStderr.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) { close(blockUntilStderrIsClosed) }).
		Once()

	factory := transport.NewFactory()

	// Act
	client, err := factory.NewClient(mockSubProcessStdio)

	// Assert
	require.NoError(t, err, "NewClient should not return an error")
	assert.NotNil(t, client, "Client should not be nil")
}

func TestFactory_NewClient_NilStdin(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(nil).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	factory := transport.NewFactory()

	// Act
	client, err := factory.NewClient(mockSubProcessStdio)

	// Assert
	require.Error(t, err, "NewClient should return an error when stdin is nil")
	assert.Nil(t, client, "Client should be nil when error occurs")
}

func TestFactory_NewClient_NilStdout(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(nil).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	factory := transport.NewFactory()

	// Act
	client, err := factory.NewClient(mockSubProcessStdio)

	// Assert
	require.Error(t, err, "NewClient should return an error when stdout is nil")
	assert.Nil(t, client, "Client should be nil when error occurs")
}

func TestFactory_NewClient_NilStderr(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(nil).
		Once()

	factory := transport.NewFactory()

	// Act
	client, err := factory.NewClient(mockSubProcessStdio)

	// Assert
	require.Error(t, err, "NewClient should return an error when stderr is nil")
	assert.Nil(t, client, "Client should be nil when error occurs")
}

func TestFactory_NewReceiver_HappyPath(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) { close(blockUntilStdinIsClosed) }).
		Once()

	factory := transport.NewFactory()

	// Act
	receiver, err := factory.NewReceiver(mockOSStdio)

	// Assert
	require.NoError(t, err, "NewReceiver should not return an error")
	assert.NotNil(t, receiver, "Receiver should not be nil")
}

func TestFactory_NewReceiver_NilStdin(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(nil).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	factory := transport.NewFactory()

	// Act
	receiver, err := factory.NewReceiver(mockOSStdio)

	// Assert
	require.Error(t, err, "NewReceiver should return an error when stdin is nil")
	assert.Nil(t, receiver, "Receiver should be nil when error occurs")
}

func TestFactory_NewReceiver_NilStdout(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(nil).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	factory := transport.NewFactory()

	// Act
	receiver, err := factory.NewReceiver(mockOSStdio)

	// Assert
	require.Error(t, err, "NewReceiver should return an error when stdout is nil")
	assert.Nil(t, receiver, "Receiver should be nil when error occurs")
}

func TestFactory_NewReceiver_NilStderr(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(nil).
		Once()

	factory := transport.NewFactory()

	// Act
	receiver, err := factory.NewReceiver(mockOSStdio)

	// Assert
	require.Error(t, err, "NewReceiver should return an error when stderr is nil")
	assert.Nil(t, receiver, "Receiver should be nil when error occurs")
}
