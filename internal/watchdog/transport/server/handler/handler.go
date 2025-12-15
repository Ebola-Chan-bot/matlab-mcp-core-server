// Copyright 2025 The MathWorks, Inc.

package handler

import (
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
)

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type ProcessHandler interface {
	WatchProcessAndGetTerminationChan(processPid int) <-chan struct{}
	KillProcess(processPid int) error
}

type Handler struct {
	logger         entities.Logger
	processHandler ProcessHandler

	lock              *sync.Mutex
	processPIDsToKill map[int]struct{}
	shutdownFuncs     []func()
}

func New(
	loggerFactory LoggerFactory,
	processHandler ProcessHandler,
) *Handler {
	return &Handler{
		logger:         loggerFactory.GetGlobalLogger(),
		processHandler: processHandler,

		lock:              &sync.Mutex{},
		processPIDsToKill: make(map[int]struct{}),
		shutdownFuncs:     make([]func(), 0),
	}
}

func (h *Handler) HandleProcessToKill(req messages.ProcessToKillRequest) (messages.ProcessToKillResponse, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.logger.
		With("pid", req.PID).
		Info("Adding process to kill")
	h.processPIDsToKill[req.PID] = struct{}{}

	return messages.ProcessToKillResponse{}, nil
}

func (h *Handler) RegisterShutdownFunction(fn func()) {
	h.shutdownFuncs = append(h.shutdownFuncs, fn)
}

func (h *Handler) HandleShutdown(_ messages.ShutdownRequest) (messages.ShutdownResponse, error) {
	for _, fn := range h.shutdownFuncs {
		fn()
	}
	return messages.ShutdownResponse{}, nil
}

func (h *Handler) TerminateAllProcesses() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.logger.
		With("count", len(h.processPIDsToKill)).
		Info("Trying to terminate children")

	for pid := range h.processPIDsToKill {
		h.logger.
			With("pid", pid).
			Debug("Killing process")

		if err := h.processHandler.KillProcess(pid); err != nil {
			h.logger.
				WithError(err).
				With("pid", pid).
				Error("Failed to kill child")
		}
	}
}
