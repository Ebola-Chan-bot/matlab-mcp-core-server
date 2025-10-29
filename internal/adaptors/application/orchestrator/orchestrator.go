// Copyright 2025 The MathWorks, Inc.

package orchestrator

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type LifecycleSignaler interface {
	RequestShutdown()
	WaitForShutdownToComplete() error
}

type Config interface {
	UseSingleMATLABSession() bool
	RecordToLogger(logger entities.Logger)
}

type Server interface {
	Run() error
}

type WatchdogClient interface {
	Start() error
	Stop() error
}

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type GlobalMATLAB interface {
	Initialize(ctx context.Context, logger entities.Logger) error
}

type Directory interface {
	BaseDir() string
}

// Orchestrator
type Orchestrator struct {
	lifecycleSignaler LifecycleSignaler
	config            Config
	server            Server
	watchdogClient    WatchdogClient
	logger            entities.Logger
	osSignaler        OSSignaler
	globalMATLAB      GlobalMATLAB
}

func New(
	lifecycleSignaler LifecycleSignaler,
	config Config,
	server Server,
	watchdogClient WatchdogClient,
	loggerFactory LoggerFactory,
	osSignaler OSSignaler,
	globalMATLAB GlobalMATLAB,
	directory Directory,
) *Orchestrator {
	orchestrator := &Orchestrator{
		lifecycleSignaler: lifecycleSignaler,
		config:            config,
		server:            server,
		watchdogClient:    watchdogClient,
		logger:            loggerFactory.GetGlobalLogger().With("log-dir", directory.BaseDir()),
		osSignaler:        osSignaler,
		globalMATLAB:      globalMATLAB,
	}
	return orchestrator
}

func (o *Orchestrator) StartAndWaitForCompletion(ctx context.Context) error {
	defer func() {
		o.logger.Info("Initiating MATLAB MCP Core Server application shutdown")
		o.lifecycleSignaler.RequestShutdown()

		err := o.lifecycleSignaler.WaitForShutdownToComplete()
		if err != nil {
			o.logger.WithError(err).Warn("MATLAB MCP Core Server application shutdown failed")
		}

		o.logger.Debug("Shutdown functions have all completed, stopping the watchdog")
		err = o.watchdogClient.Stop()
		if err != nil {
			o.logger.WithError(err).Warn("Watchdog shutdown failed")
		}

		o.logger.Info("MATLAB MCP Core Server application shutdown complete")
	}()

	o.logger.Info("Initiating MATLAB MCP Core Server application startup")
	o.config.RecordToLogger(o.logger)

	err := o.watchdogClient.Start()
	if err != nil {
		return err
	}

	serverErrC := make(chan error, 1)
	go func() {
		serverErrC <- o.server.Run()
	}()

	if o.config.UseSingleMATLABSession() {
		err := o.globalMATLAB.Initialize(ctx, o.logger)
		if err != nil {
			o.logger.WithError(err).Warn("MATLAB global initialization failed")
		}
	}

	o.logger.Info("MATLAB MCP Core Server application startup complete")

	select {
	case <-o.osSignaler.InterruptSignalChan():
		o.logger.Info("Received termination signal")
		return nil
	case err := <-serverErrC:
		return err
	}
}
