// Copyright 2025 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

func NewServer(
	httpServerFactory HTTPServerFactory,
	logger entities.Logger,
	handler Handler,
) (*Server, error) {
	return newServer(
		httpServerFactory,
		logger,
		handler,
	)
}
