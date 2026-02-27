// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/pkg/config"
)

func NewConfigAdaptor(internalConfig internalconfig.GenericConfig, messageCatalog definition.MessageCatalog) config.Config {
	return newConfigAdaptor(internalConfig, messageCatalog)
}
