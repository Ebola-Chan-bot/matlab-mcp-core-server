// Copyright 2025 The MathWorks, Inc.

package processhandler

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type OSWrapper interface {
	WaitForProcessToComplete(processPid int)
	FindProcess(processPid int) osfacade.Process
}

type ProcessHandler struct {
	logger    entities.Logger
	osWrapper OSWrapper
}

func New(
	loggerFactory LoggerFactory,
	osWrapper OSWrapper,
) *ProcessHandler {
	return &ProcessHandler{
		logger:    loggerFactory.GetGlobalLogger(),
		osWrapper: osWrapper,
	}
}

func (f *ProcessHandler) WatchProcessAndGetTerminationChan(processPid int) <-chan struct{} {
	logger := f.logger.With("process-pid", processPid)
	logger.Debug("Watching process and notifying if it terminates")

	parentTerminatedC := make(chan struct{})

	go func() {
		f.osWrapper.WaitForProcessToComplete(processPid)
		logger.Debug("Process terminated")
		close(parentTerminatedC)
	}()

	return parentTerminatedC
}

func (f *ProcessHandler) KillProcess(processPid int) error {
	if process := f.osWrapper.FindProcess(processPid); process != nil {
		return process.Kill()
	}
	return nil
}
