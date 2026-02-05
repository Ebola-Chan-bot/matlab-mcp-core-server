// Copyright 2026 The MathWorks, Inc.

package logger

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)

	With(key string, value any) Logger
	WithError(err error) Logger
}
