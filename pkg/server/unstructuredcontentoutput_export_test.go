// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
)

func (t *ToolWithUnstructuredContentOutput[ToolInput]) ToInternal(loggerFactoryInstance basetool.LoggerFactory) basetool.ToolWithUnstructuredContentOutput[ToolInput] {
	return t.toInternal(loggerFactoryInstance).(basetool.ToolWithUnstructuredContentOutput[ToolInput])
}
