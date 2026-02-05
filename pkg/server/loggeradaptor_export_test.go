// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

func NewLoggerAdaptor(logger entities.Logger) logger.Logger {
	return newLoggerAdaptor(logger)
}
