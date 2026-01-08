// Copyright 2025-2026 The MathWorks, Inc.

package config

import (
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

func New(osLayer OSLayer, parser Parser) (*config, messages.Error) {
	return newConfig(osLayer, parser)
}
