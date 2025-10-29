// Copyright 2025 The MathWorks, Inc.

package matlablocator_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/matlablocator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ListDiscoveredMatlabInfo_SingleMatlab(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABRootGetter := &mocks.MockMATLABRootGetter{}
	defer mockMATLABRootGetter.AssertExpectations(t)

	mockMATLABVersionGetter := &mocks.MockMATLABVersionGetter{}
	defer mockMATLABVersionGetter.AssertExpectations(t)

	// Mock discovery service to return MATLAB locations
	const dummyPathToMATLAB = "/path/to/matlab/R2023a"
	mockMATLABRootGetter.EXPECT().
		GetAll(mockLogger).
		Return([]string{
			dummyPathToMATLAB,
		}).
		Once()

	// Mock version service to return version info
	dummyMatlabVersionInfo := datatypes.MatlabVersionInfo{
		ReleaseFamily: "R2023a",
		ReleasePhase:  "release",
		UpdateLevel:   0,
	}
	mockMATLABVersionGetter.EXPECT().
		Get(dummyPathToMATLAB).
		Return(dummyMatlabVersionInfo, nil).
		Once()

	service := matlablocator.New(mockMATLABRootGetter, mockMATLABVersionGetter)

	// Act
	result := service.ListDiscoveredMatlabInfo(mockLogger)

	// Assert
	require.NotNil(t, result)
	require.Len(t, result.MatlabInfo, 1)

	// Verify the MATLAB info
	assert.Equal(t, dummyPathToMATLAB, result.MatlabInfo[0].Location)
	assert.Equal(t, dummyMatlabVersionInfo, result.MatlabInfo[0].Version)
}

func TestService_ListDiscoveredMatlabInfo_MultipleMatlabs(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABRootGetter := &mocks.MockMATLABRootGetter{}
	defer mockMATLABRootGetter.AssertExpectations(t)

	mockMATLABVersionGetter := &mocks.MockMATLABVersionGetter{}
	defer mockMATLABVersionGetter.AssertExpectations(t)

	dummyPaths := []string{
		"C:\\Program Files\\MATLAB\\R2023a",
		"C:\\Program Files\\MATLAB\\R2022b",
	}
	dummyMatlabInfos := map[string]datatypes.MatlabVersionInfo{
		dummyPaths[0]: {
			ReleaseFamily: "R2023a",
			ReleasePhase:  "release",
			UpdateLevel:   1,
		},
		dummyPaths[1]: {
			ReleaseFamily: "R2022b",
			ReleasePhase:  "release",
			UpdateLevel:   2,
		},
	}

	// Mock discovery service to return multiple MATLAB locations
	mockMATLABRootGetter.EXPECT().
		GetAll(mockLogger).
		Return(dummyPaths).
		Once()

	// Mock version service to return version info for each MATLAB
	for _, path := range dummyPaths {
		mockMATLABVersionGetter.EXPECT().
			Get(path).
			Return(dummyMatlabInfos[path], nil).
			Once()
	}

	service := matlablocator.New(mockMATLABRootGetter, mockMATLABVersionGetter)

	// Act
	result := service.ListDiscoveredMatlabInfo(mockLogger)

	// Assert
	require.NotNil(t, result)
	require.Len(t, result.MatlabInfo, len(dummyPaths))

	for i := range dummyPaths {
		path := dummyPaths[i]
		expectedInfo := dummyMatlabInfos[path]

		assert.Equal(t, path, result.MatlabInfo[i].Location)
		assert.Equal(t, expectedInfo.ReleaseFamily, result.MatlabInfo[i].Version.ReleaseFamily)
		assert.Equal(t, expectedInfo.ReleasePhase, result.MatlabInfo[i].Version.ReleasePhase)
		assert.Equal(t, expectedInfo.UpdateLevel, result.MatlabInfo[i].Version.UpdateLevel)
	}
}

func TestService_ListDiscoveredMatlabInfo_NoMatlabsFound(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABRootGetter := &mocks.MockMATLABRootGetter{}
	defer mockMATLABRootGetter.AssertExpectations(t)

	mockMATLABVersionGetter := &mocks.MockMATLABVersionGetter{}
	defer mockMATLABVersionGetter.AssertExpectations(t)

	// Mock discovery service to return empty list
	mockMATLABRootGetter.EXPECT().
		GetAll(mockLogger).
		Return([]string{}).
		Once()

	service := matlablocator.New(mockMATLABRootGetter, mockMATLABVersionGetter)

	// Act
	result := service.ListDiscoveredMatlabInfo(mockLogger)

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result.MatlabInfo)
}

func TestService_ListDiscoveredMatlabInfo_VersionServiceError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABRootGetter := &mocks.MockMATLABRootGetter{}
	defer mockMATLABRootGetter.AssertExpectations(t)

	mockMATLABVersionGetter := &mocks.MockMATLABVersionGetter{}
	defer mockMATLABVersionGetter.AssertExpectations(t)

	// Mock discovery service to return MATLAB locations
	mockMATLABRootGetter.EXPECT().
		GetAll(mockLogger).
		Return([]string{
			"/path/to/matlab/R2023a",
		}).
		Once()

	// Mock version service to return an error
	expectedError := assert.AnError
	mockMATLABVersionGetter.EXPECT().
		Get("/path/to/matlab/R2023a").
		Return(datatypes.MatlabVersionInfo{}, expectedError).
		Once()

	service := matlablocator.New(mockMATLABRootGetter, mockMATLABVersionGetter)

	// Act
	result := service.ListDiscoveredMatlabInfo(mockLogger)

	// Assert
	assert.Empty(t, result.MatlabInfo, "Result should be empty when there's an error")
	assert.Len(t, mockLogger.WarnLogs(), 1, "Should log expected number of warning messages")
}
