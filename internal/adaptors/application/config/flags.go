// Copyright 2025 The MathWorks, Inc.

package config

import (
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/spf13/pflag"
)

const (
	versionMode             = "version"
	versionModeDefaultValue = false

	disableTelemetry             = "disable-telemetry"
	disableTelemetryDefaultValue = false

	useSingleMATLABSession             = "use-single-matlab-session"
	useSingleMATLABSessionDefaultValue = true

	preferredLocalMATLABRoot             = "matlab-root"
	preferredLocalMATLABRootDefaultValue = ""

	preferredMATLABStartingDirectory             = "initial-working-folder"
	preferredMATLABStartingDirectoryDefaultValue = ""

	logLevel             = "log-level"
	logLevelDefaultValue = "info"

	watchdogMode             = "watchdog"
	watchdogModeDefaultValue = false
)

func setupFlags(flagSet *pflag.FlagSet) error {
	flagSet.Bool(versionMode, versionModeDefaultValue,
		"Display the version of the MATLAB MCP Core Server.",
	)

	flagSet.Bool(disableTelemetry, disableTelemetryDefaultValue,
		"Disable collection of usage data. By default, this software may collect information about you and your usage and send it to MathWorks. This data helps us improve our products and services.",
	)

	flagSet.Bool(useSingleMATLABSession, useSingleMATLABSessionDefaultValue,
		"When true, a MATLAB session is started when a MATLAB MCP Core Server starts, and stopped when the server is shut down. When false, the server can manage multiple MATLAB sessions.",
	)

	flagSet.String(logLevel, logLevelDefaultValue,
		"The log level to use for the global logger (for session logs, the clients sets the log level). Valid values are: debug, info, warn, error.",
	)

	flagSet.String(preferredLocalMATLABRoot, preferredLocalMATLABRootDefaultValue,
		fmt.Sprintf("When %s is true, if this is set, defines which local MATLAB installation to use. If not set, the first MATLAB installation on the PATH will be used.", useSingleMATLABSession),
	)

	flagSet.String(preferredMATLABStartingDirectory, preferredMATLABStartingDirectoryDefaultValue,
		fmt.Sprintf("When %s is true, if this is set, defines which starting directory MATLAB will use. If not set, MATLAB will use the default MATLAB's starting directory.", useSingleMATLABSession))

	// Hidden flags, for internal use only
	flagSet.Bool(watchdogMode, watchdogModeDefaultValue,
		"INTERNAL USE ONLY.",
	)
	err := flagSet.MarkHidden(watchdogMode)
	if err != nil {
		return err
	}

	return nil
}

func createConfigWithFlagValues(osLayer OSLayer, flagSet *pflag.FlagSet, args []string) (*Config, error) {
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	versionMode, err := flagSet.GetBool(versionMode)
	if err != nil {
		return nil, err
	}

	disableTelemetry, err := flagSet.GetBool(disableTelemetry)
	if err != nil {
		return nil, err
	}

	useSingleMATLABSession, err := flagSet.GetBool(useSingleMATLABSession)
	if err != nil {
		return nil, err
	}

	logLevel, err := flagSet.GetString(logLevel)
	if err != nil {
		return nil, err
	}

	switch logLevel {
	case string(entities.LogLevelDebug), string(entities.LogLevelInfo), string(entities.LogLevelWarn), string(entities.LogLevelError):
		break
	default:
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}

	preferredLocalMATLABRoot, err := flagSet.GetString(preferredLocalMATLABRoot)
	if err != nil {
		return nil, err
	}

	preferredMATLABStartingDirectory, err := flagSet.GetString(preferredMATLABStartingDirectory)
	if err != nil {
		return nil, err
	}

	watchdogMode, err := flagSet.GetBool(watchdogMode)
	if err != nil {
		return nil, err
	}

	return &Config{
		osLayer: osLayer,

		versionMode:                      versionMode,
		disableTelemetry:                 disableTelemetry,
		useSingleMATLABSession:           useSingleMATLABSession,
		logLevel:                         entities.LogLevel(logLevel),
		preferredLocalMATLABRoot:         preferredLocalMATLABRoot,
		preferredMATLABStartingDirectory: preferredMATLABStartingDirectory,
		watchdogMode:                     watchdogMode,
	}, nil
}
