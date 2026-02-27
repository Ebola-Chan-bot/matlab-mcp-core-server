// Copyright 2026 The MathWorks, Inc.

package config

import "github.com/matlab/matlab-mcp-core-server/pkg/i18n"

type Config interface {
	Get(key string, expectedType any) (any, i18n.Error)
}
