// Copyright 2026 The MathWorks, Inc.

package server

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/tools"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	pkgtools "github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type Tool = publictypes.Tool

type HandlerForToolWithUnstructuredContentOutput[ToolInput any] func(ctx context.Context, request pkgtools.CallRequest, inputs ToolInput) (pkgtools.RichContent, i18n.Error)

func NewToolWithUnstructuredContentOutput[ToolInput any](definition pkgtools.Definition, handler HandlerForToolWithUnstructuredContentOutput[ToolInput]) Tool {
	return tools.NewUnstructured(definition, tools.UnstructuredHandler[ToolInput](handler))
}

type HandlerForToolWithStructuredContentOutput[ToolInput, ToolOutput any] func(ctx context.Context, request pkgtools.CallRequest, inputs ToolInput) (ToolOutput, i18n.Error)

func NewToolWithStructuredContentOutput[ToolInput, ToolOutput any](definition pkgtools.Definition, handler HandlerForToolWithStructuredContentOutput[ToolInput, ToolOutput]) Tool {
	return tools.NewStructured(definition, tools.StructuredHandler[ToolInput, ToolOutput](handler))
}
