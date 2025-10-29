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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGlobalMATLAB_Initialize_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedSelectedMATLABRoot := "/mock/matlab/path"
	expectedMATLABStartingDir := "/home/myuser"

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        expectedSelectedMATLABRoot,
		StartingDirectory: expectedMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	mockSessionID := entities.SessionID(123)

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
		Return(mockSessionID, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	err := globalMATLABSession.Initialize(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
}

func TestGlobalMATLAB_Initialize_SelectFirstMATLABVersionOnPathError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return("", expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	err := globalMATLABSession.Initialize(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestGlobalMATLAB_Initialize_SelectMatlabStartingDirError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	sessionId := entities.SessionID(123)
	ctx := t.Context()
	expectedSelectedMATLABRoot := "/mock/matlab/path"
	expectedMATLABStartingDir := ""

	mockLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        expectedSelectedMATLABRoot,
		StartingDirectory: expectedMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return(expectedSelectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, assert.AnError).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger, mockLocalSessionDetails).
		Return(sessionId, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	err := globalMATLABSession.Initialize(ctx, mockLogger)

	// Arrange
	require.NoError(t, err)
}

func TestGlobalMATLAB_Initialize_StartMATLABSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	var sessionIDZeroValue entities.SessionID
	expectedError := assert.AnError

	ctx := t.Context()
	expectedSelectedMATLABRoot := "/mock/matlab/path"
	expectedMATLABStartingDir := "/home/myuser"

	mockLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        expectedSelectedMATLABRoot,
		StartingDirectory: expectedMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return(expectedSelectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger, mockLocalSessionDetails).
		Return(sessionIDZeroValue, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	err := globalMATLABSession.Initialize(ctx, mockLogger)

	// Arrange
	require.ErrorIs(t, err, expectedError)
}

func TestGlobalMATLAB_Initialize_ReturnsCachedErrorOnSubsequentInitializeCalls(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	var sessionIDZeroValue entities.SessionID
	expectedError := assert.AnError

	ctx := t.Context()
	expectedSelectedMATLABRoot := "/mock/matlab/path"
	expectedMATLABStartingDir := "/home/myuser"

	mockLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:        expectedSelectedMATLABRoot,
		StartingDirectory: expectedMATLABStartingDir,
		ShowMATLABDesktop: true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return(expectedSelectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger, mockLocalSessionDetails).
		Return(sessionIDZeroValue, expectedError).
		Once()

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return(expectedSelectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	err1 := globalMATLABSession.Initialize(ctx, mockLogger)
	err2 := globalMATLABSession.Initialize(ctx, mockLogger)

	// Arrange
	require.ErrorIs(t, err1, expectedError)
	require.ErrorIs(t, err2, expectedError)
}

func TestGlobalMATLAB_Initialize_ConcurrentCallsWaitForCompletion(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

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

	firstCallCompleted := make(chan error)
	secondCallCompleted := make(chan error)

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

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger).
		Return(expectedSelectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
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
		secondCallCompleted <- globalMATLABSession.Initialize(ctx, mockLogger)
	}()

	select {
	case <-secondCallCompleted:
		t.Fatal("Second Initialize call completed before first call was unblocked")
	case <-time.After(100 * time.Millisecond):
		// Second call is still waiting
	}

	close(blockStartMATLAB)

	// Assert
	require.NoError(t, <-firstCallCompleted)
	require.NoError(t, <-secondCallCompleted)
}
