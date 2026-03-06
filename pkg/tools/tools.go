// Copyright 2026 The MathWorks, Inc.

package tools

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/tools"
)

type Definition = publictypes.ToolDefinition

func NewDefinition(name, title, description string, annotations Annotations) Definition {
	return tools.NewDefinition(name, title, description, annotations)
}

type CallRequest = publictypes.ToolCallRequest

type RichContent = publictypes.RichContent
