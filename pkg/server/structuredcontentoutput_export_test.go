// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
)

func (t *ToolWithStructuredContentOutput[ToolInput, ToolOutput]) ToInternal(loggerFactoryInstance basetool.LoggerFactory) basetool.ToolWithStructuredContentOutput[ToolInput, ToolOutput] {
	return t.toInternal(loggerFactoryInstance).(basetool.ToolWithStructuredContentOutput[ToolInput, ToolOutput])
}
