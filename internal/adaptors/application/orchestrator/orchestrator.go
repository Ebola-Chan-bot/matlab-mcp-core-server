// Copyright 2025-2026 The MathWorks, Inc.

package orchestrator

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type LifecycleSignaler interface {
	RequestShutdown()
	WaitForShutdownToComplete() error
}

type ApplicationDefinition interface {
	Dependencies(resources definition.DependenciesProviderResources) (any, error)
	Tools(resources definition.ToolsProviderResources) []tools.Tool
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Server interface {
	Run(tools []tools.Tool) error
}

type WatchdogClient interface {
	Start() error
	Stop() error
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type GlobalMATLAB interface {
	Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error)
}

type DirectoryFactory interface {
	Directory() (directory.Directory, messages.Error)
}

// Orchestrator
type Orchestrator struct {
	lifecycleSignaler     LifecycleSignaler
	applicationDefinition ApplicationDefinition
	configFactory         ConfigFactory
	server                Server
	watchdogClient        WatchdogClient
	loggerFactory         LoggerFactory
	osSignaler            OSSignaler
	globalMATLAB          GlobalMATLAB
	directoryFactory      DirectoryFactory
}

func New(
	lifecycleSignaler LifecycleSignaler,
	applicationDefinition ApplicationDefinition,
	configFactory ConfigFactory,
	server Server,
	watchdogClient WatchdogClient,
	loggerFactory LoggerFactory,
	osSignaler OSSignaler,
	globalMATLAB GlobalMATLAB,
	directoryFactory DirectoryFactory,
) *Orchestrator {
	orchestrator := &Orchestrator{
		lifecycleSignaler:     lifecycleSignaler,
		applicationDefinition: applicationDefinition,
		configFactory:         configFactory,
		server:                server,
		watchdogClient:        watchdogClient,
		loggerFactory:         loggerFactory,
		osSignaler:            osSignaler,
		globalMATLAB:          globalMATLAB,
		directoryFactory:      directoryFactory,
	}
	return orchestrator
}

func (o *Orchestrator) StartAndWaitForCompletion(ctx context.Context) error {
	config, messagesErr := o.configFactory.Config()
	if messagesErr != nil {
		return messagesErr
	}

	logger, messagesErr := o.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}

	defer func() {
		logger.Info("Initiating application shutdown")
		o.lifecycleSignaler.RequestShutdown()

		err := o.lifecycleSignaler.WaitForShutdownToComplete()
		if err != nil {
			logger.WithError(err).Warn("Application shutdown failed")
		}

		logger.Debug("Shutdown functions have all completed, stopping the watchdog")
		err = o.watchdogClient.Stop()
		if err != nil {
			logger.WithError(err).Warn("Watchdog shutdown failed")
		}

		logger.Info("Application shutdown complete")
	}()

	logger.Info("Initiating application startup")
	config.RecordToLogger(logger)

	directory, messagesErr := o.directoryFactory.Directory()
	if messagesErr != nil {
		return messagesErr
	}
	directory.RecordToLogger(logger)

	err := o.watchdogClient.Start()
	if err != nil {
		return err
	}

	logger.Debug("Building SDK dependencies")
	dependencies, err := o.applicationDefinition.Dependencies(definition.NewDependenciesProviderResources(
		logger,
	))
	if err != nil {
		return err
	}

	logger.Debug("Building SDK tools")
	tools := o.applicationDefinition.Tools(definition.NewToolsProviderResources(
		logger,
		dependencies,
		o.loggerFactory,
	))

	serverErrC := make(chan error, 1)
	go func() {
		serverErrC <- o.server.Run(tools)
	}()

	if config.UseSingleMATLABSession() && config.InitializeMATLABOnStartup() {
		_, err := o.globalMATLAB.Client(ctx, logger)
		if err != nil {
			logger.WithError(err).Warn("MATLAB global initialization failed")
		}
	}

	logger.Info("Application startup complete")

	select {
	case <-o.osSignaler.InterruptSignalChan():
		logger.Info("Received termination signal")
		return nil
	case err := <-serverErrC:
		return err
	}
}
