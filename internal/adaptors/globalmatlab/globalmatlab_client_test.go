// Copyright 2025 The MathWorks, Inc.

package globalmatlab_test

import (
	"context"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGlobalMATLAB_Client_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	mockPreferredMATLABRoot := ""
	mockPreferredMATLABStartingDir := ""

	mockLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        mockPreferredMATLABRoot,
		StartingDirectory: mockPreferredMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), mockLocalSessionDetails).
		Return(sessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(mockSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mockSessionClient, client)
}

func TestGlobalMATLAB_Client_StartMATLABSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	const sessionIDThatShouldBeUnused = entities.SessionID(0)
	expectedError := assert.AnError

	mockPreferredMATLABRoot := ""
	mockPreferredMATLABStartingDir := ""

	mockLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        mockPreferredMATLABRoot,
		StartingDirectory: mockPreferredMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger, mockLocalSessionDetails).
		Return(sessionIDThatShouldBeUnused, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_GetMATLABSessionClientError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	mockPreferredMATLABRoot := ""
	mockPreferredMATLABStartingDir := ""
	expectedError := assert.AnError

	mockLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        mockPreferredMATLABRoot,
		StartingDirectory: mockPreferredMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), mockLocalSessionDetails).
		Return(sessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_ReturnsStartMATLABCachedErrorFromInitialize(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	matlabRoot := ""
	startingDir := ""
	expectedError := assert.AnError

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return(matlabRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(startingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger, entities.LocalSessionDetails{
			MATLABRoot:        matlabRoot,
			StartingDirectory: startingDir,
			ShowMATLABDesktop: true,
		}).
		Return(entities.SessionID(0), expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	err := globalMATLABSession.Initialize(ctx, mockLogger)
	require.ErrorIs(t, err, expectedError)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	assert.Nil(t, client)
	require.ErrorIs(t, err, expectedError)
}

func TestGlobalMATLAB_Client_ReturnsCachedErrorOnSubsequentClientCalls(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	matlabRoot := ""
	startingDir := ""
	expectedError := assert.AnError

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger, entities.LocalSessionDetails{
			MATLABRoot:        matlabRoot,
			StartingDirectory: startingDir,
			ShowMATLABDesktop: true,
		}).
		Return(entities.SessionID(0), expectedError).
		Once()

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	assert.Nil(t, client1)
	require.ErrorIs(t, err1, expectedError)

	assert.Nil(t, client2)
	require.ErrorIs(t, err2, expectedError)
}

func TestGlobalMATLAB_Client_ConcurrentCallsWaitForCompletion(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSelectedMATLABRoot := ""
	expectedMATLABStartingDir := ""
	mockSessionID := entities.SessionID(123)

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        expectedSelectedMATLABRoot,
		StartingDirectory: expectedMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	blockStartMATLAB := make(chan struct{})
	startMATLABCalled := make(chan struct{})

	type clientResult struct {
		client entities.MATLABSessionClient
		err    error
	}

	firstCallCompleted := make(chan clientResult)
	secondCallCompleted := make(chan clientResult)

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger, expectedLocalSessionDetails).
		Run(func(ctx context.Context, logger entities.Logger, details entities.SessionDetails) {
			close(startMATLABCalled)
			<-blockStartMATLAB
		}).
		Return(mockSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger, mockSessionID).
		Return(mockSessionClient, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger, mockSessionID).
		Return(mockSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	go func() {
		client, err := globalMATLABSession.Client(ctx, mockLogger)
		firstCallCompleted <- clientResult{client: client, err: err}
	}()

	<-startMATLABCalled

	go func() {
		client, err := globalMATLABSession.Client(ctx, mockLogger)
		secondCallCompleted <- clientResult{client: client, err: err}
	}()

	select {
	case <-secondCallCompleted:
		t.Fatal("Second Client call completed before first call was unblocked")
	case <-time.After(100 * time.Millisecond):
		// Second call is still waiting
	}

	close(blockStartMATLAB)
	result1 := <-firstCallCompleted
	result2 := <-secondCallCompleted

	// Assert
	require.NoError(t, result1.err)
	assert.Equal(t, mockSessionClient, result1.client)

	require.NoError(t, result2.err)
	assert.Equal(t, mockSessionClient, result2.client)
}

func TestGlobalMATLAB_Client_WaitUntilInitializeFinish(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSelectedMATLABRoot := "/mock/matlab/path"
	expectedMATLABStartingDir := "/home/myuser"
	mockSessionID := entities.SessionID(123)

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        expectedSelectedMATLABRoot,
		StartingDirectory: expectedMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	blockStartMATLAB := make(chan struct{})
	startMATLABCalled := make(chan struct{})

	type clientResult struct {
		client entities.MATLABSessionClient
		err    error
	}

	firstCallCompleted := make(chan error)
	secondCallCompleted := make(chan clientResult)

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return(expectedSelectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger, expectedLocalSessionDetails).
		Run(func(ctx context.Context, logger entities.Logger, details entities.SessionDetails) {
			close(startMATLABCalled)
			<-blockStartMATLAB
		}).
		Return(mockSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger, mockSessionID).
		Return(mockSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	go func() {
		firstCallCompleted <- globalMATLABSession.Initialize(ctx, mockLogger)
	}()

	<-startMATLABCalled

	go func() {
		client, err := globalMATLABSession.Client(ctx, mockLogger)
		secondCallCompleted <- clientResult{client: client, err: err}
	}()

	select {
	case <-secondCallCompleted:
		t.Fatal("Client call completed before Initialized")
	case <-time.After(100 * time.Millisecond):
		// client call is still waiting
	}

	close(blockStartMATLAB)
	result := <-secondCallCompleted

	// Assert
	require.NoError(t, <-firstCallCompleted)

	require.NoError(t, result.err)
	assert.Equal(t, mockSessionClient, result.client)
}
