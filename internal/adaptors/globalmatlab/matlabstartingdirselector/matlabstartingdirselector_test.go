// Copyright 2025 The MathWorks, Inc.

package matlabstartingdirselector_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/matlabstartingdirselector"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab/matlabstartingdirselector"
	osFacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	// Act
	selector := matlabstartingdirselector.New(mockConfig, mockOSLayer)

	// Assert
	assert.NotNil(t, selector)
}

func TestMATLABStartingDirSelector_GetMatlabStartDir_HappyPath(t *testing.T) {
	testCases := []struct {
		name        string
		os          string
		homeDir     string
		expectedDir string
	}{
		{
			name:        "Windows",
			os:          "windows",
			homeDir:     "/Users/testuser",
			expectedDir: filepath.Join("/Users/testuser", "Documents"),
		},
		{
			name:        "Darwin",
			os:          "darwin",
			homeDir:     "/Users/testuser",
			expectedDir: filepath.Join("/Users/testuser", "Documents"),
		},
		{
			name:        "Linux",
			os:          "linux",
			homeDir:     "/home/testuser",
			expectedDir: "/home/testuser",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &mocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockConfig := &mocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockFileInfo := &osFacademocks.MockFileInfo{}
			defer mockFileInfo.AssertExpectations(t)

			selector := matlabstartingdirselector.New(mockConfig, mockOSLayer)

			mockConfig.EXPECT().
				PreferredMATLABStartingDirectory().
				Return("").
				Once()

			mockOSLayer.EXPECT().
				UserHomeDir().
				Return(tc.homeDir, nil).
				Once()

			mockOSLayer.EXPECT().
				GOOS().
				Return(tc.os).
				Once()

			mockOSLayer.EXPECT().
				Stat(tc.expectedDir).
				Return(mockFileInfo, nil).
				Once()

			// Act
			result, err := selector.SelectMatlabStartingDir()

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expectedDir, result)
		})
	}
}

func TestMATLABStartingDirSelector_GetMatlabStartDir_PreferredStartingDirectorySet_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockPreferredMATLABStartingDir := "/custom/preferred/directory"

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return(mockPreferredMATLABStartingDir).
		Once()

	mockOSLayer.EXPECT().
		Stat(mockPreferredMATLABStartingDir).
		Return(mockFileInfo, nil).
		Once()

	selector := matlabstartingdirselector.New(mockConfig, mockOSLayer)

	// Act
	result, err := selector.SelectMatlabStartingDir()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mockPreferredMATLABStartingDir, result)
}

func TestMATLABStartingDirSelector_GetMatlabStartDir_UnknownOS_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileInfo := &osFacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	selector := matlabstartingdirselector.New(mockConfig, mockOSLayer)

	homeDir := "/home/testuser"
	unknownOS := "freebsd"

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return(homeDir, nil).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return(unknownOS).
		Once()

	mockOSLayer.EXPECT().
		Stat(homeDir).
		Return(mockFileInfo, nil).
		Once()

	// Act
	result, err := selector.SelectMatlabStartingDir()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, homeDir, result)
}

func TestMATLABStartingDirSelector_GetMatlabStartDir_UserHomeDirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	selector := matlabstartingdirselector.New(mockConfig, mockOSLayer)
	expectedError := assert.AnError

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		UserHomeDir().
		Return("", expectedError).
		Once()

	// Act
	result, err := selector.SelectMatlabStartingDir()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestMATLABStartingDirSelector_GetMatlabStartDir_StatErrorOnHomeDir(t *testing.T) {
	testCases := []struct {
		name    string
		os      string
		homeDir string
	}{
		{
			name:    "Windows - Stat Error",
			os:      "windows",
			homeDir: "C:\\Users\\testuser",
		},
		{
			name:    "Darwin - Stat Error",
			os:      "darwin",
			homeDir: "/Users/testuser",
		},
		{
			name:    "Linux - Stat Error",
			os:      "linux",
			homeDir: "/home/testuser",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &mocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockConfig := &mocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			selector := matlabstartingdirselector.New(mockConfig, mockOSLayer)

			expectedDir := tc.homeDir
			expectedError := assert.AnError
			if tc.os == "windows" || tc.os == "darwin" {
				expectedDir = filepath.Join(tc.homeDir, "Documents")
			}

			mockConfig.EXPECT().
				PreferredMATLABStartingDirectory().
				Return("").
				Once()

			mockOSLayer.EXPECT().
				UserHomeDir().
				Return(tc.homeDir, nil).
				Once()

			mockOSLayer.EXPECT().
				GOOS().
				Return(tc.os).
				Once()

			mockOSLayer.EXPECT().
				Stat(expectedDir).
				Return(nil, expectedError).
				Once()

			// Act
			result, err := selector.SelectMatlabStartingDir()

			// Assert
			require.ErrorIs(t, err, expectedError)
			assert.Empty(t, result)
		})
	}
}

func TestMATLABStartingDirSelector_GetMatlabStartDir_StatErrorOnPreferredMATLABStartingDir(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	selector := matlabstartingdirselector.New(mockConfig, mockOSLayer)
	expectedPreferredMATLABStartingDir := "/some/path/that/doesnt/exist"
	expectedError := assert.AnError

	mockConfig.EXPECT().
		PreferredMATLABStartingDirectory().
		Return(expectedPreferredMATLABStartingDir).
		Once()

	mockOSLayer.EXPECT().
		Stat(expectedPreferredMATLABStartingDir).
		Return(nil, expectedError).
		Once()

	// Act
	result, err := selector.SelectMatlabStartingDir()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}
