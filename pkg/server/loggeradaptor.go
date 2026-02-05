// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

type loggerAdaptor struct {
	entities.Logger
}

func newLoggerAdaptor(logger entities.Logger) loggerAdaptor {
	return loggerAdaptor{
		Logger: logger,
	}
}

func (l loggerAdaptor) With(key string, value any) logger.Logger {
	newLogger := l.Logger.With(key, value)
	return newLoggerAdaptor(newLogger)
}

func (l loggerAdaptor) WithError(err error) logger.Logger {
	newLogger := l.Logger.WithError(err)
	return newLoggerAdaptor(newLogger)
}
