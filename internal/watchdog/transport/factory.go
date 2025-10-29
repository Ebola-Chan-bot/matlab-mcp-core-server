// Copyright 2025 The MathWorks, Inc.

package transport

import "github.com/matlab/matlab-mcp-core-server/internal/entities"

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewClient(subProcessStdio entities.SubProcessStdio) (Client, error) {
	return NewStdioClient(subProcessStdio)
}

func (f *Factory) NewReceiver(osStdio entities.OSStdio) (Receiver, error) {
	return NewStdioReceiver(osStdio)
}
