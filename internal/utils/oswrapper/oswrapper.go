// Copyright 2025 The MathWorks, Inc.

package oswrapper

import (
	"syscall"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

const defaultCheckParentAliveInterval = 1 * time.Second

type OSLayer interface {
	FindProcess(pid int) (osfacade.Process, error)
	GOOS() string
}

type OSWrapper struct {
	osLayer OSLayer

	goos                     string
	checkParentAliveInterval time.Duration
}

func New(
	osLayer OSLayer,
) *OSWrapper {
	return &OSWrapper{
		osLayer: osLayer,

		goos:                     osLayer.GOOS(),
		checkParentAliveInterval: defaultCheckParentAliveInterval,
	}
}

func (w *OSWrapper) FindProcess(processPid int) osfacade.Process {
	proc, err := w.osLayer.FindProcess(processPid)
	if err != nil {
		// Can't find the process, return nil
		return nil
	}

	// From the FindProcess doc:
	// On Unix systems, FindProcess always succeeds and returns a Process for the given pid,
	// regardless of whether the process exists.
	// To test whether the process actually exists, see whether p.Signal(syscall.Signal(0)) reports an error.
	if w.goos != "windows" {
		err = proc.Signal(syscall.Signal(0))
		if err != nil {
			// Process does not exist or is not accessible, return nil.
			return nil
		}
	}

	// The process was found and is responsive, return it.
	return proc
}

func (w *OSWrapper) WaitForProcessToComplete(processPid int) {
	// On Windows, we can poll FindProcess, because once it has been started, it may keep returning, even after being stopped:
	// https://github.com/golang/go/issues/33814
	//
	// Therefore, we Wait().
	//
	// On the other hand, on Unix system, you can only Wait() child processes, so in this case, we poll.

	if w.goos == "windows" {
		if process := w.FindProcess(processPid); process != nil {
			process.Wait() //nolint:gosec,errcheck // It doesn't matter why we stopped waiting
		}
	} else {
		ticker := time.NewTicker(w.checkParentAliveInterval)
		defer ticker.Stop()

		for range ticker.C {
			if w.FindProcess(processPid) != nil {
				continue
			}
			return
		}
	}
}
