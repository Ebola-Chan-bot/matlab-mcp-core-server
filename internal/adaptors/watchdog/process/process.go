// Copyright 2025 The MathWorks, Inc.

package process

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

type OSLayer interface {
	Command(name string, arg ...string) osfacade.Cmd
	Executable() (string, error)
}

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type Directory interface {
	BaseDir() string
	ID() string
}

type Config interface {
	LogLevel() entities.LogLevel
}

type Process struct {
	osLayer OSLayer
	cmd     osfacade.Cmd
	logger  entities.Logger
}

func New(
	osLayer OSLayer,
	loggerFactory LoggerFactory,
	directory Directory,
	config Config,
) (*Process, error) {
	logger := loggerFactory.GetGlobalLogger()

	programPath, err := osLayer.Executable()
	if err != nil {
		logger.WithError(err).Error("Failed to get executable path")
		return nil, err
	}
	cmd := osLayer.Command(programPath,
		"--"+flags.WatchdogMode,
		"--"+flags.BaseDir, directory.BaseDir(),
		"--"+flags.ServerInstanceID, directory.ID(),
		"--"+flags.LogLevel, string(config.LogLevel()),
	)

	cmd.SetSysProcAttr(getSysProcAttrForDetachingAProcess())

	process := &Process{
		osLayer: osLayer,
		cmd:     cmd,
		logger:  logger,
	}

	return process, nil
}

func (p *Process) Start() error {
	if err := p.cmd.Start(); err != nil {
		p.logger.WithError(err).Error("Failed to start watchdog process")
		return err
	}

	return nil
}
