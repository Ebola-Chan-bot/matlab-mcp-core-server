// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
)

func (p ToolsProvider[Dependencies]) ToInternal() definition.ToolsProvider {
	return p.toInternal()
}
