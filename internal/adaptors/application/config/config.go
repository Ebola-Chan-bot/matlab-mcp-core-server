// Copyright 2025-2026 The MathWorks, Inc.

package config

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type validatedArguments struct {
	logLevel         entities.LogLevel
	disableTelemetry bool

	versionMode  bool
	helpMode     bool
	watchdogMode bool

	useSingleMATLABSession           bool
	initializeMATLABOnStartup        bool
	existingSessionOnly              bool
	preferredLocalMATLABRoot         string
	preferredMATLABStartingDirectory string
	displayMode                      entities.DisplayMode

	baseDirectory    string
	serverInstanceID string
}

type rawConfig struct {
	parameters []entities.Parameter
	parsedArgs map[string]any
}

func (c *rawConfig) Get(key string) (any, messages.Error) {
	return getForKey(c.parsedArgs, key)
}

type config struct {
	buildInfo BuildInfo

	*rawConfig
	validatedArguments
}

func newConfig(osLayer OSLayer, parser Parser, buildInfo BuildInfo) (*config, messages.Error) {
	parameters, parsedArgs, err := parser.Parse(osLayer.Args()[1:])
	if err != nil {
		return nil, err
	}

	rawCfg := &rawConfig{
		parameters: parameters,
		parsedArgs: parsedArgs,
	}

	validated, err := validateArguments(rawCfg)
	if err != nil {
		return nil, err
	}

	return &config{
		buildInfo:          buildInfo,
		rawConfig:          rawCfg,
		validatedArguments: validated,
	}, nil
}

func (c *config) Get(key string) (any, messages.Error) {
	return getForKey(c.parsedArgs, key)
}

func (c *config) Version() string {
	return c.buildInfo.FullVersion()
}

func (c *config) LogLevel() entities.LogLevel {
	return c.logLevel
}

func (c *config) DisableTelemetry() bool {
	return c.disableTelemetry
}

func (c *config) VersionMode() bool {
	return c.versionMode
}

func (c *config) HelpMode() bool {
	return c.helpMode
}

func (c *config) WatchdogMode() bool {
	return c.watchdogMode
}

func (c *config) UseSingleMATLABSession() bool {
	return c.useSingleMATLABSession
}

func (c *config) PreferredLocalMATLABRoot() string {
	return c.preferredLocalMATLABRoot
}

func (c *config) PreferredMATLABStartingDirectory() string {
	return c.preferredMATLABStartingDirectory
}

func (c *config) InitializeMATLABOnStartup() bool {
	return c.initializeMATLABOnStartup
}

func (c *config) ExistingSessionOnly() bool {
	return c.existingSessionOnly
}

func (c *config) ShouldShowMATLABDesktop() bool {
	switch c.displayMode {
	case entities.DisplayModeDesktop:
		return true
	case entities.DisplayModeNoDesktop:
		return false
	default:
		return true
	}
}

func (c *config) BaseDir() string {
	return c.baseDirectory
}

func (c *config) ServerInstanceID() string {
	return c.serverInstanceID
}

func (c *config) RecordToLogger(logger entities.Logger) {
	for _, param := range c.parameters {
		if param.GetRecordToLog() {
			value, err := c.Get(param.GetID())
			if err == nil {
				logger = logger.With(param.GetID(), value)
			}
		}
	}
	logger.Info("Configuration state")
}

func validateArguments(rawCfg *rawConfig) (validatedArguments, messages.Error) {
	logLevel, err := get(rawCfg, defaultparameters.LogLevel())
	if err != nil {
		return validatedArguments{}, err
	}

	switch logLevel {
	case string(entities.LogLevelDebug), string(entities.LogLevelInfo), string(entities.LogLevelWarn), string(entities.LogLevelError):
	default:
		return validatedArguments{}, messages.New_StartupErrors_InvalidLogLevel_Error(logLevel)
	}

	disableTelemetry, err := get(rawCfg, defaultparameters.DisableTelemetry())
	if err != nil {
		return validatedArguments{}, err
	}

	versionMode, err := get(rawCfg, defaultparameters.VersionMode())
	if err != nil {
		return validatedArguments{}, err
	}

	helpMode, err := get(rawCfg, defaultparameters.HelpMode())
	if err != nil {
		return validatedArguments{}, err
	}

	watchdogMode, err := get(rawCfg, defaultparameters.WatchdogMode())
	if err != nil {
		return validatedArguments{}, err
	}

	useSingleMATLABSession, err := get(rawCfg, defaultparameters.UseSingleMATLABSession())
	if err != nil {
		return validatedArguments{}, err
	}

	initializeMATLABOnStartup, err := get(rawCfg, defaultparameters.InitializeMATLABOnStartup())
	if err != nil {
		return validatedArguments{}, err
	}

	if !useSingleMATLABSession {
		initializeMATLABOnStartup = false
	}

	existingSessionOnly, err := get(rawCfg, defaultparameters.ExistingSessionOnly())
	if err != nil {
		return validatedArguments{}, err
	}

	preferredLocalMATLABRoot, err := get(rawCfg, defaultparameters.PreferredLocalMATLABRoot())
	if err != nil {
		return validatedArguments{}, err
	}

	preferredMATLABStartingDirectory, err := get(rawCfg, defaultparameters.PreferredMATLABStartingDirectory())
	if err != nil {
		return validatedArguments{}, err
	}

	displayMode, err := get(rawCfg, defaultparameters.MATLABDisplayMode())
	if err != nil {
		return validatedArguments{}, err
	}

	switch displayMode {
	case string(entities.DisplayModeDesktop), string(entities.DisplayModeNoDesktop):
		break
	default:
		return validatedArguments{}, messages.New_StartupErrors_InvalidDisplayMode_Error(displayMode)
	}

	baseDirectory, err := get(rawCfg, defaultparameters.BaseDir())
	if err != nil {
		return validatedArguments{}, err
	}

	serverInstanceID, err := get(rawCfg, defaultparameters.ServerInstanceID())
	if err != nil {
		return validatedArguments{}, err
	}

	return validatedArguments{
		logLevel:         entities.LogLevel(logLevel),
		disableTelemetry: disableTelemetry,

		versionMode:  versionMode,
		helpMode:     helpMode,
		watchdogMode: watchdogMode,

		useSingleMATLABSession:           useSingleMATLABSession,
		initializeMATLABOnStartup:        initializeMATLABOnStartup,
		existingSessionOnly:              existingSessionOnly,
		preferredLocalMATLABRoot:         preferredLocalMATLABRoot,
		preferredMATLABStartingDirectory: preferredMATLABStartingDirectory,
		displayMode:                      entities.DisplayMode(displayMode),

		baseDirectory:    baseDirectory,
		serverInstanceID: serverInstanceID,
	}, nil
}

func getForKey(args map[string]any, key string) (any, messages.Error) {
	if value, ok := args[key]; ok {
		return value, nil
	}
	return nil, messages.New_StartupErrors_InvalidParameterKey_Error(key)
}
