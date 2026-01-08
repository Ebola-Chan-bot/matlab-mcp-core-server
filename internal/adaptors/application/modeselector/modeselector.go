// Copyright 2025-2026 The MathWorks, Inc.

package modeselector

import (
	"context"
	"fmt"
	"io"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Parser interface {
	Usage() string
}

type WatchdogProcessFactory interface { //nolint:iface // Intentional interface for deps injection
	Create() (entities.Mode, error)
}

type OrchestratorFactory interface { //nolint:iface // Intentional interface for deps injection
	Create() (entities.Mode, error)
}

type OSLayer interface {
	Stdout() io.Writer
}

// ModeSelector is the top level object of the MATLAB MCP Core Server.
// It will be imported in `main.go` to start the application, and wait for it's completion.
// It will select which mode to run in based on the configuration, and defer construction of the required objects until the mode is known.
type ModeSelector struct {
	configFactory          ConfigFactory
	watchdogProcessFactory WatchdogProcessFactory
	orchestratorFactory    OrchestratorFactory
	osLayer                OSLayer
	parser                 Parser
}

func New(
	configFactory ConfigFactory,
	parser Parser,
	watchdogProcessFactory WatchdogProcessFactory,
	orchestratorFactory OrchestratorFactory,
	osLayer OSLayer,
) *ModeSelector {
	return &ModeSelector{
		configFactory:          configFactory,
		parser:                 parser,
		watchdogProcessFactory: watchdogProcessFactory,
		orchestratorFactory:    orchestratorFactory,
		osLayer:                osLayer,
	}
}

func (a *ModeSelector) StartAndWaitForCompletion(ctx context.Context) error {
	config, err := a.configFactory.Config()
	if err != nil {
		return err
	}

	switch {
	case config.HelpMode():
		_, err := fmt.Fprintf(a.osLayer.Stdout(), "%s\n", a.parser.Usage())
		return err
	case config.VersionMode():
		_, err := fmt.Fprintf(a.osLayer.Stdout(), "%s\n", config.Version())
		return err
	case config.WatchdogMode():
		watchdogProcess, err := a.watchdogProcessFactory.Create()
		if err != nil {
			return err
		}

		return watchdogProcess.StartAndWaitForCompletion(ctx)
	default:
		orchestrator, err := a.orchestratorFactory.Create()
		if err != nil {
			return err
		}

		return orchestrator.StartAndWaitForCompletion(ctx)
	}
}
