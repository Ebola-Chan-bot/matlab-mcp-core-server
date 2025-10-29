// Copyright 2025 The MathWorks, Inc.

package processhandler

import (
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

type OSWrapper interface {
	WaitForProcessToComplete(processPid int)
	FindProcess(processPid int) osfacade.Process
}

type ProcessHandler struct {
	osWrapper OSWrapper
}

func New(
	osWrapper OSWrapper,
) *ProcessHandler {
	return &ProcessHandler{
		osWrapper: osWrapper,
	}
}

func (f *ProcessHandler) WatchProcessAndGetTerminationChan(processPid int) <-chan struct{} {
	parentTerminatedC := make(chan struct{})

	go func() {
		f.osWrapper.WaitForProcessToComplete(processPid)
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
