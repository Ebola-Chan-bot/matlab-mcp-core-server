// Copyright 2025 The MathWorks, Inc.

package config

import (
	"encoding/json"
	"runtime/debug"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/spf13/pflag"
)

type OSLayer interface {
	Args() []string
	ReadBuildInfo() (info *debug.BuildInfo, ok bool)
}

type Config struct {
	osLayer OSLayer

	versionMode                      bool
	disableTelemetry                 bool
	useSingleMATLABSession           bool
	logLevel                         entities.LogLevel
	preferredLocalMATLABRoot         string
	preferredMATLABStartingDirectory string
	watchdogMode                     bool
}

func New(
	osLayer OSLayer,
) (*Config, error) {
	flagSet := pflag.NewFlagSet(pflag.CommandLine.Name(), pflag.ContinueOnError)
	err := setupFlags(flagSet)
	if err != nil {
		return nil, err
	}
	return createConfigWithFlagValues(osLayer, flagSet, osLayer.Args()[1:])
}

// Version will return the application version string.
// The version returned will be `version` if set, using ldflags during build.
// Otherwise, it will return the version from the build info.
func (c *Config) Version() string {
	buildInfo, ok := c.osLayer.ReadBuildInfo()

	finalVersion := strings.TrimSpace(version)

	if ok && version == unsetVersion {
		finalVersion = buildInfo.Main.Version
	}

	return buildInfo.Main.Path + " " + finalVersion
}

func (c *Config) VersionMode() bool {
	return c.versionMode
}

func (c *Config) DisableTelemetry() bool {
	return c.disableTelemetry
}

func (c *Config) UseSingleMATLABSession() bool {
	return c.useSingleMATLABSession
}

func (c *Config) LogLevel() entities.LogLevel {
	return c.logLevel
}

func (c *Config) PreferredLocalMATLABRoot() string {
	return c.preferredLocalMATLABRoot
}

func (c *Config) PreferredMATLABStartingDirectory() string {
	return c.preferredMATLABStartingDirectory
}

func (c *Config) WatchdogMode() bool {
	return c.watchdogMode
}

func (c *Config) RecordToLogger(logger entities.Logger) {
	data, err := json.Marshal(map[string]any{
		disableTelemetry:                 c.disableTelemetry,
		useSingleMATLABSession:           c.useSingleMATLABSession,
		logLevel:                         c.logLevel,
		preferredLocalMATLABRoot:         c.preferredLocalMATLABRoot,
		preferredMATLABStartingDirectory: c.preferredMATLABStartingDirectory,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to serialize configuration")
		return
	}
	logger.With("config", string(data)).Info("Configuration state")
}
