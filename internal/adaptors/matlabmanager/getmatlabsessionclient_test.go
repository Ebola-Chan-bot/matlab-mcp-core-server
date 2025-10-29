// Copyright 2025 The MathWorks, Inc.

package matlabmanager_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager"
	sessionstoremocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabsessionstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMATLABManager_GetMATLABSessionClient_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionClient := &sessionstoremocks.MockMATLABSessionClientWithCleanup{}

	sessionID := entities.SessionID(123)
	ctx := t.Context()

	mockSessionStore.EXPECT().
		Get(sessionID).
		Return(mockSessionClient, nil).
		Once()

	manager := matlabmanager.New(mockMATLABServices, mockSessionStore, mockClientFactory)

	// Act
	client, err := manager.GetMATLABSessionClient(ctx, mockLogger, sessionID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mockSessionClient, client)
}

func TestMATLABManager_GetMATLABSessionClient_SessionStoreError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	sessionID := entities.SessionID(123)
	ctx := t.Context()
	expectedError := assert.AnError

	mockSessionStore.EXPECT().
		Get(sessionID).
		Return(nil, expectedError).
		Once()

	manager := matlabmanager.New(mockMATLABServices, mockSessionStore, mockClientFactory)

	// Act
	client, err := manager.GetMATLABSessionClient(ctx, mockLogger, sessionID)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}
