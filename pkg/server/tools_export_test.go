// Copyright 2026 The MathWorks, Inc.

package server

import (
	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type ToolArray []Tool

func (t ToolArray) ToInternal(loggerFactoryInstance basetool.LoggerFactory) []internaltools.Tool {
	return toolArray(t).toInternal(loggerFactoryInstance)
}

func NewToolCallRequest(logger logger.Logger) *tools.CallRequest {
	return newToolCallRequest(logger)
}
