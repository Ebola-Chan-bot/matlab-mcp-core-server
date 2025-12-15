// Copyright 2025 The MathWorks, Inc.

package server

import (
	"net/http"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpserverfactory"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
)

type HTTPServerFactory interface {
	NewServerOverUDS(handlers map[string]http.HandlerFunc) (httpserverfactory.HttpServer, error)
}

type Handler interface {
	HandleProcessToKill(req messages.ProcessToKillRequest) (messages.ProcessToKillResponse, error)
	HandleShutdown(req messages.ShutdownRequest) (messages.ShutdownResponse, error)
}

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type Factory struct {
	httpServerFactory HTTPServerFactory
	loggerFactory     LoggerFactory
	handler           Handler
}

func NewFactory(
	httpServerFactory HTTPServerFactory,
	loggerFactory LoggerFactory,
	handler Handler,
) *Factory {
	return &Factory{
		httpServerFactory: httpServerFactory,
		loggerFactory:     loggerFactory,
		handler:           handler,
	}
}

func (f *Factory) New() (transport.Server, error) {
	return newServer(
		f.httpServerFactory,
		f.loggerFactory.GetGlobalLogger(),
		f.handler,
	)
}
