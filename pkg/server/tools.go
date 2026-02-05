// Copyright 2026 The MathWorks, Inc.

package server

import (
	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type Tool interface {
	toInternal(loggerFactory basetool.LoggerFactory) internaltools.Tool
}

type toolArray []Tool

func (t toolArray) toInternal(loggerFactoryInstance basetool.LoggerFactory) []internaltools.Tool {
	internalTools := make([]internaltools.Tool, len(t))

	for i, tool := range t {
		internalTools[i] = tool.toInternal(loggerFactoryInstance)
	}

	return internalTools
}

func newToolCallRequest(logger logger.Logger) *tools.CallRequest {
	return &tools.CallRequest{
		Logger: logger,
	}
}
