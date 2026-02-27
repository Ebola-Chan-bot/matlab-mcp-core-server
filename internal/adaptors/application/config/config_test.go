// Copyright 2025-2026 The MathWorks, Inc.

package config_test

import (
	"path/filepath"
	"runtime/debug"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func configDefaultParsedArgs() map[string]any {
	params := []entities.Parameter{
		defaultparameters.HelpMode(),
		defaultparameters.VersionMode(),
		defaultparameters.DisableTelemetry(),
		defaultparameters.UseSingleMATLABSession(),
		defaultparameters.LogLevel(),
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.BaseDir(),
		defaultparameters.WatchdogMode(),
		defaultparameters.ServerInstanceID(),
		defaultparameters.InitializeMATLABOnStartup(),
		defaultparameters.MATLABDisplayMode(),
	}

	result := make(map[string]any)
	for _, p := range params {
		result[p.GetID()] = p.GetDefaultValue()
	}
	return result
}

func TestNewConfig_InvalidParameterType(t *testing.T) {
	testCases := []struct {
		name         string
		key          string
		invalidValue any
		expectedType string
	}{
		{name: "LogLevel wrong type", key: defaultparameters.LogLevel().GetID(), invalidValue: 123, expectedType: "string"},
		{name: "UseSingleMATLABSession wrong type", key: defaultparameters.UseSingleMATLABSession().GetID(), invalidValue: "true", expectedType: "bool"},
		{name: "InitializeMATLABOnStartup wrong type", key: defaultparameters.InitializeMATLABOnStartup().GetID(), invalidValue: "false", expectedType: "bool"},
		{name: "VersionMode wrong type", key: defaultparameters.VersionMode().GetID(), invalidValue: "false", expectedType: "bool"},
		{name: "HelpMode wrong type", key: defaultparameters.HelpMode().GetID(), invalidValue: "false", expectedType: "bool"},
		{name: "DisableTelemetry wrong type", key: defaultparameters.DisableTelemetry().GetID(), invalidValue: "false", expectedType: "bool"},
		{name: "PreferredLocalMATLABRoot wrong type", key: defaultparameters.PreferredLocalMATLABRoot().GetID(), invalidValue: 123, expectedType: "string"},
		{name: "PreferredMATLABStartingDirectory wrong type", key: defaultparameters.PreferredMATLABStartingDirectory().GetID(), invalidValue: 123, expectedType: "string"},
		{name: "BaseDir wrong type", key: defaultparameters.BaseDir().GetID(), invalidValue: 123, expectedType: "string"},
		{name: "WatchdogMode wrong type", key: defaultparameters.WatchdogMode().GetID(), invalidValue: "false", expectedType: "bool"},
		{name: "ServerInstanceID wrong type", key: defaultparameters.ServerInstanceID().GetID(), invalidValue: 123, expectedType: "string"},
		{name: "MATLABDisplayMode wrong type", key: defaultparameters.MATLABDisplayMode().GetID(), invalidValue: 123, expectedType: "string"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			parsedArgs[tc.key] = tc.invalidValue

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, nil).
				Once()

			expectedError := messages.New_StartupErrors_InvalidParameterType_Error(tc.key, tc.expectedType)

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser)

			// Assert
			require.Equal(t, expectedError, err)
			assert.Nil(t, cfg)
		})
	}
}

func TestNewConfig_MissingParameter(t *testing.T) {
	testCases := []struct {
		name       string
		missingKey string
	}{
		{name: "missing LogLevel", missingKey: defaultparameters.LogLevel().GetID()},
		{name: "missing UseSingleMATLABSession", missingKey: defaultparameters.UseSingleMATLABSession().GetID()},
		{name: "missing InitializeMATLABOnStartup", missingKey: defaultparameters.InitializeMATLABOnStartup().GetID()},
		{name: "missing VersionMode", missingKey: defaultparameters.VersionMode().GetID()},
		{name: "missing HelpMode", missingKey: defaultparameters.HelpMode().GetID()},
		{name: "missing DisableTelemetry", missingKey: defaultparameters.DisableTelemetry().GetID()},
		{name: "missing PreferredLocalMATLABRoot", missingKey: defaultparameters.PreferredLocalMATLABRoot().GetID()},
		{name: "missing PreferredMATLABStartingDirectory", missingKey: defaultparameters.PreferredMATLABStartingDirectory().GetID()},
		{name: "missing BaseDir", missingKey: defaultparameters.BaseDir().GetID()},
		{name: "missing WatchdogMode", missingKey: defaultparameters.WatchdogMode().GetID()},
		{name: "missing ServerInstanceID", missingKey: defaultparameters.ServerInstanceID().GetID()},
		{name: "missing MATLABDisplayMode", missingKey: defaultparameters.MATLABDisplayMode().GetID()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			parsedArgs := configDefaultParsedArgs()
			delete(parsedArgs, tc.missingKey)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, parsedArgs, nil).
				Once()

			expectedError := messages.New_StartupErrors_InvalidParameterKey_Error(tc.missingKey)

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser)

			// Assert
			require.Equal(t, expectedError, err)
			assert.Nil(t, cfg)
		})
	}
}

func TestNewConfig_ParseError(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(nil, nil, messages.AnError).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser)

	// Assert
	require.ErrorIs(t, err, messages.AnError)
	assert.Nil(t, cfg, "Config should be nil")
}

func TestNewConfig_InvalidLogLevel(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.LogLevel().GetID()] = "invalid-level"

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, nil).
		Once()

	expectedError := messages.New_StartupErrors_InvalidLogLevel_Error("invalid-level")

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser)

	// Assert
	require.Equal(t, expectedError, err)
	assert.Nil(t, cfg, "Config should be nil")
}

func TestConfig_Version(t *testing.T) {
	modulePath := "github.com/matlab/matlab-mcp-core-server"

	testCases := []struct {
		name            string
		buildInfoOK     bool
		moduleVersion   string
		expectedVersion string
	}{
		{
			name:            "version from build info",
			buildInfoOK:     true,
			moduleVersion:   "v1.2.3",
			expectedVersion: modulePath + " v1.2.3",
		},
		{
			name:            "devel fallback",
			buildInfoOK:     true,
			moduleVersion:   "(devel)",
			expectedVersion: modulePath + " (devel)",
		},
		{
			name:            "build info unavailable",
			buildInfoOK:     false,
			moduleVersion:   "",
			expectedVersion: "(unknown)",
		},
		{
			name:            "empty version string",
			buildInfoOK:     true,
			moduleVersion:   "",
			expectedVersion: modulePath + " (devel)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return([]entities.Parameter{}, configDefaultParsedArgs(), nil).
				Once()

			var buildInfo *debug.BuildInfo
			if tc.buildInfoOK {
				buildInfo = &debug.BuildInfo{
					Main: debug.Module{
						Path:    modulePath,
						Version: tc.moduleVersion,
					},
				}
			}

			mockOSLayer.EXPECT().
				ReadBuildInfo().
				Return(buildInfo, tc.buildInfoOK).
				Once()

			// Act
			cfg, err := config.NewConfig(mockOSLayer, mockParser)
			require.NoError(t, err)

			version := cfg.Version()

			// Assert
			require.Equal(t, tc.expectedVersion, version)
		})
	}
}

func TestConfig_InitializeMATLABOnStartup_DisabledWhenNotSingleSession(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.UseSingleMATLABSession().GetID()] = false
	parsedArgs[defaultparameters.InitializeMATLABOnStartup().GetID()] = true

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return([]entities.Parameter{}, parsedArgs, nil).
		Once()

	// Act
	cfg, err := config.NewConfig(mockOSLayer, mockParser)

	// Assert
	require.NoError(t, err)
	assert.False(t, cfg.InitializeMATLABOnStartup(), "InitializeMATLABOnStartup should be false when UseSingleMATLABSession is false")
}

func TestConfig_RecordToLogger_HappyPath(t *testing.T) {
	// Arrange
	parsedArgs := configDefaultParsedArgs()
	parsedArgs[defaultparameters.DisableTelemetry().GetID()] = true
	parsedArgs[defaultparameters.PreferredMATLABStartingDirectory().GetID()] = filepath.Join("home", "user")
	parsedArgs[defaultparameters.LogLevel().GetID()] = string(entities.LogLevelDebug)
	parsedArgs[defaultparameters.PreferredLocalMATLABRoot().GetID()] = filepath.Join("home", "matlab")
	parsedArgs[defaultparameters.UseSingleMATLABSession().GetID()] = false

	expectedLogMessage := "Configuration state"
	expectedConfigField := map[string]any{
		defaultparameters.DisableTelemetry().GetID():                 true,
		defaultparameters.PreferredMATLABStartingDirectory().GetID(): filepath.Join("home", "user"),
		defaultparameters.LogLevel().GetID():                         string(entities.LogLevelDebug),
		defaultparameters.PreferredLocalMATLABRoot().GetID():         filepath.Join("home", "matlab"),
		defaultparameters.UseSingleMATLABSession().GetID():           false,
		defaultparameters.InitializeMATLABOnStartup().GetID():        false,
	}

	parameters := []entities.Parameter{
		defaultparameters.DisableTelemetry(),
		defaultparameters.UseSingleMATLABSession(),
		defaultparameters.LogLevel(),
		defaultparameters.PreferredLocalMATLABRoot(),
		defaultparameters.PreferredMATLABStartingDirectory(),
		defaultparameters.InitializeMATLABOnStartup(),
		defaultparameters.HelpMode(),
		defaultparameters.VersionMode(),
		defaultparameters.BaseDir(),
		defaultparameters.WatchdogMode(),
		defaultparameters.ServerInstanceID(),
	}

	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(parameters, parsedArgs, nil)

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	cfg, err := config.NewConfig(mockOSLayer, mockParser)
	require.NoError(t, err)

	testLogger := testutils.NewInspectableLogger()

	// Act
	cfg.RecordToLogger(testLogger)

	// Assert
	infoLogs := testLogger.InfoLogs()
	require.Len(t, infoLogs, 1)

	fields, found := infoLogs[expectedLogMessage]
	require.True(t, found, "Expected log message not found")

	for expectedField, expectedValue := range expectedConfigField {
		actualValue, exists := fields[expectedField]
		require.True(t, exists, "%s field not found in log", expectedField)
		assert.Equal(t, expectedValue, actualValue, "%s field has incorrect value", expectedField)
	}
}
