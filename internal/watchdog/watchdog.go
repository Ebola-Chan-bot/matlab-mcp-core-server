// Copyright 2025 The MathWorks, Inc.

package watchdog

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
)

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type OSLayer interface {
	Getppid() int
}

type ProcessHandler interface {
	WatchProcessAndGetTerminationChan(processPid int) <-chan struct{}
	KillProcess(processPid int) error
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type ServerHandler interface {
	RegisterShutdownFunction(fn func())
	TerminateAllProcesses()
}

type ServerFactory interface {
	New() (transport.Server, error)
}

type SocketFactory interface {
	Socket() (socket.Socket, error)
}

type Watchdog struct {
	logger         entities.Logger
	osLayer        OSLayer
	processHandler ProcessHandler
	osSignaler     OSSignaler
	serverHandler  ServerHandler
	serverFactory  ServerFactory
	socketFactory  SocketFactory

	parentPID         int
	shutdownRequestC  chan struct{}
	shutdownResponseC chan struct{}
}

func New(
	loggerFactory LoggerFactory,
	osLayer OSLayer,
	processHandler ProcessHandler,
	osSignaler OSSignaler,
	serverHandler ServerHandler,
	serverFactory ServerFactory,
	socketFactory SocketFactory,
) *Watchdog {
	return &Watchdog{
		logger:         loggerFactory.GetGlobalLogger(),
		osLayer:        osLayer,
		processHandler: processHandler,
		osSignaler:     osSignaler,
		serverHandler:  serverHandler,
		serverFactory:  serverFactory,
		socketFactory:  socketFactory,

		shutdownRequestC:  make(chan struct{}),
		shutdownResponseC: make(chan struct{}),
	}
}

func (w *Watchdog) StartAndWaitForCompletion(_ context.Context) error {
	socket, err := w.socketFactory.Socket()
	if err != nil {
		return err
	}

	server, err := w.serverFactory.New()
	if err != nil {
		return err
	}

	go func() {
		if err := server.Start(socket.Path()); err != nil {
			w.logger.WithError(err).Error("Server Start method returned an error")
		}
	}()

	defer func() {
		if err := server.Stop(); err != nil {
			w.logger.WithError(err).Error("Failed to stop server")
		}
	}()

	w.logger.Info("Watchdog process has started")
	defer w.logger.Info("Watchdog process has exited")

	w.parentPID = w.osLayer.Getppid()

	w.serverHandler.RegisterShutdownFunction(func() {
		close(w.shutdownRequestC)
		// Make sure we broke out of the select, before returning
		<-w.shutdownResponseC
	})

	select {
	case <-w.shutdownRequestC:
		w.logger.Debug("Graceful shutdown signal received")
		// Ackownledge shutdown
		close(w.shutdownResponseC)

	case <-w.processHandler.WatchProcessAndGetTerminationChan(w.parentPID):
		w.logger.Debug("Lost connection to parent, shutting down")

	case <-w.osSignaler.InterruptSignalChan():
		w.logger.Debug("Received unexpected graceful shutdown OS signal")
	}

	w.serverHandler.TerminateAllProcesses()

	return nil
}
