// Copyright 2025 The MathWorks, Inc.

package basetool

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const UnexpectedErrorPrefixForLLM = "unexpected error occurred: "

type LoggerFactory interface {
	NewMCPSessionLogger(session *mcp.ServerSession) entities.Logger
	GetGlobalLogger() entities.Logger
}

type ToolAdder[ToolInput, ToolOutput any] interface {
	AddTool(server *mcp.Server, tool *mcp.Tool, handler mcp.ToolHandlerFor[ToolInput, ToolOutput])
}

type tool[ToolInput any, ToolOutput any] struct {
	name          string
	title         string
	description   string
	loggerFactory LoggerFactory
	toolAdder     ToolAdder[ToolInput, ToolOutput]
}

func (t tool[_, _]) Name() string {
	return t.name
}

func (t tool[_, _]) Title() string {
	return t.title
}

func (t tool[_, _]) Description() string {
	return t.description
}

func (_ tool[ToolInput, _]) GetInputSchema() (any, error) {
	return jsonschema.For[ToolInput](&jsonschema.ForOptions{})
}
