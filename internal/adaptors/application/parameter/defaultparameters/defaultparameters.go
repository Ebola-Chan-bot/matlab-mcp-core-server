// Copyright 2026 The MathWorks, Inc.

package defaultparameters

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

func HelpMode() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		"HelpMode",
		"help",
		false,
		"",
		messages.CLIMessages_HelpDescription,
		false,
		false,
	)
}

func VersionMode() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		"VersionMode",
		"version",
		false,
		"",
		messages.CLIMessages_VersionDescription,
		false,
		false,
	)
}

func DisableTelemetry() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		"DisableTelemetry",
		"disable-telemetry",
		false,
		"",
		messages.CLIMessages_DisableTelemetryDescription,
		false,
		true,
	)
}

func PreferredLocalMATLABRoot() *parameter.Parameter[string] {
	return parameter.NewParameter(
		"PreferredLocalMATLABRoot",
		"matlab-root",
		false,
		"",
		messages.CLIMessages_PreferredLocalMATLABRootDescription,
		"",
		true,
	)
}

func PreferredMATLABStartingDirectory() *parameter.Parameter[string] {
	return parameter.NewParameter(
		"PreferredMATLABStartingDirectory",
		"initial-working-folder",
		false,
		"",
		messages.CLIMessages_PreferredMATLABStartingDirectoryDescription,
		"",
		true,
	)
}

func BaseDir() *parameter.Parameter[string] {
	return parameter.NewParameter(
		"BaseDir",
		"log-folder",
		false,
		"",
		messages.CLIMessages_BaseDirDescription,
		"",
		false,
	)
}

func LogLevel() *parameter.Parameter[string] {
	return parameter.NewParameter(
		"LogLevel",
		"log-level",
		false,
		"",
		messages.CLIMessages_LogLevelDescription,
		"info",
		true,
	)
}

func InitializeMATLABOnStartup() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		"InitializeMATLABOnStartup",
		"initialize-matlab-on-startup",
		false,
		"",
		messages.CLIMessages_InitializeMATLABOnStartupDescription,
		false,
		true,
	)
}

func MATLABDisplayMode() *parameter.Parameter[string] {
	return parameter.NewParameter(
		"MATLABDisplayMode",
		"matlab-display-mode",
		false,
		"",
		messages.CLIMessages_DisplayModeDescription,
		"desktop",
		true,
	)
}

func UseSingleMATLABSession() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		"UseSingleMATLABSession",
		"use-single-matlab-session",
		true,
		"",
		messages.CLIMessages_UseSingleMATLABSessionDescription,
		true,
		true,
	)
}

func WatchdogMode() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		"WatchdogMode",
		"watchdog",
		true,
		"",
		messages.CLIMessages_InternalUseDescription,
		false,
		false,
	)
}

func ExistingSessionOnly() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		"ExistingSessionOnly",
		"existing-session-only",
		false,
		"MATLAB_EXISTING_SESSION_ONLY",
		messages.CLIMessages_ExistingSessionOnlyDescription,
		false,
		true,
	)
}

func ServerInstanceID() *parameter.Parameter[string] {
	return parameter.NewParameter(
		"ServerInstanceID",
		"server-instance-id",
		true,
		"",
		messages.CLIMessages_InternalUseDescription,
		"",
		false,
	)
}
