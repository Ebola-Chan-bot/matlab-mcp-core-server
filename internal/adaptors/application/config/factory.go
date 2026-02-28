// Copyright 2025-2026 The MathWorks, Inc.

package config

import (
	"runtime/debug"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type Parser interface {
	Parse(args []string) ([]entities.Parameter, map[string]any, messages.Error)
}

type OSLayer interface {
	Args() []string
	ReadBuildInfo() (info *debug.BuildInfo, ok bool)
}

type GenericConfig interface {
	Get(key string) (any, messages.Error)
}

type Config interface {
	GenericConfig

	Version() string
	HelpMode() bool
	VersionMode() bool
	WatchdogMode() bool
	BaseDir() string
	ServerInstanceID() string
	UseSingleMATLABSession() bool
	InitializeMATLABOnStartup() bool
	RecordToLogger(logger entities.Logger)
	LogLevel() entities.LogLevel
	PreferredLocalMATLABRoot() string
	PreferredMATLABStartingDirectory() string
	ShouldShowMATLABDesktop() bool
	ExistingSessionOnly() bool
}

type Factory struct {
	parser  Parser
	osLayer OSLayer

	initOnce       sync.Once
	initError      messages.Error
	configInstance *config
}

func NewFactory(parser Parser, osLayer OSLayer) *Factory {
	return &Factory{
		parser:  parser,
		osLayer: osLayer,
	}
}

func (f *Factory) Config() (Config, messages.Error) {
	f.initOnce.Do(func() {
		configInstance, err := newConfig(f.osLayer, f.parser)
		if err != nil {
			f.initError = err
			return
		}

		f.configInstance = configInstance
	})

	if f.initError != nil {
		return nil, f.initError
	}

	return f.configInstance, nil
}
