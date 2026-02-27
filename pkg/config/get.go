// Copyright 2026 The MathWorks, Inc.

package config

import (
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
)

func Get[ParameterType supportedParameterValueType](cfg Config, parameter Parameter[ParameterType]) (ParameterType, i18n.Error) {
	var zeroValue ParameterType

	value, err := cfg.Get(parameter.GetID(), zeroValue)
	if err != nil {
		return zeroValue, err
	}

	castValue, ok := value.(ParameterType)
	if !ok {
		// This code path should be unreachable.
		// cfg.Get is expected to error on type mismatch.
		return zeroValue, nil
	}

	return castValue, nil
}
