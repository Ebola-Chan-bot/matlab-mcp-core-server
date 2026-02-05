// Copyright 2026 The MathWorks, Inc.

package integration

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

func NewEmptyApplication() *wire.Application {
	serverDefinition := definition.New("", "", "", nil, nil)
	return wire.Initialize(serverDefinition)
}
