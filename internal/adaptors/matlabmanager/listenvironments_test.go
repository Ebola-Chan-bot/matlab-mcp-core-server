// Copyright 2025 The MathWorks, Inc.

package matlabmanager_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMATLABManager_ListEnvironments_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABServices{}
	defer mockMATLABManager.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	dummyMatlabInfos := []datatypes.MatlabInfo{{
		Location: "/path/to/matlab/R2023a",
		Version: datatypes.MatlabVersionInfo{
			ReleaseFamily: "R2023a",
			ReleasePhase:  "release",
			UpdateLevel:   0,
		},
	}, {
		Location: "/path/to/matlab/R2022b",
		Version: datatypes.MatlabVersionInfo{
			ReleaseFamily: "R2022b",
			ReleasePhase:  "release",
			UpdateLevel:   1,
		},
	},
	}

	mockResponse := datatypes.ListMatlabInfo{
		MatlabInfo: dummyMatlabInfos,
	}
	mockMATLABManager.EXPECT().
		ListDiscoveredMatlabInfo(mockLogger.AsMockArg()).
		Return(mockResponse).
		Once()

	manager := matlabmanager.New(mockMATLABManager, mockSessionStore, mockClientFactory)
	ctx := t.Context()

	// Act
	result := manager.ListEnvironments(ctx, mockLogger)

	// Assert
	require.Len(t, result, 2)

	// Verify the outputs match the mock data
	for i := range dummyMatlabInfos {
		assert.Equal(t, dummyMatlabInfos[i].Location, result[i].MATLABRoot, "Output MATLAB root does not match input dummy data")
		assert.Equal(t, dummyMatlabInfos[i].Version.ReleaseFamily, result[i].Version, "Output MATLAB version does not match input dummy data")
	}
}

func TestMATLABManager_ListEnvironments_EmptyList(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABServices{}
	defer mockMATLABManager.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockResponse := datatypes.ListMatlabInfo{
		MatlabInfo: []datatypes.MatlabInfo{},
	}
	mockMATLABManager.EXPECT().
		ListDiscoveredMatlabInfo(mockLogger).
		Return(mockResponse).
		Once()

	manager := matlabmanager.New(mockMATLABManager, mockSessionStore, mockClientFactory)
	ctx := t.Context()

	// Act
	result := manager.ListEnvironments(ctx, mockLogger)

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result)
}
