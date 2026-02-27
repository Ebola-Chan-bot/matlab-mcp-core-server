// Copyright 2026 The MathWorks, Inc.

package parser

import (
	"strconv"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

const internalErrorText = "Unimplemented parameter type"

func (p *Parser) parseEnvVars(specifiedArgs map[string]any) messages.Error {
	for _, parameter := range p.parameters {
		if !parameter.GetActive() {
			continue
		}

		envVarName := parameter.GetEnvVarName()
		if envVarName == "" {
			continue
		}

		val, ok := p.osLayer.LookupEnv(envVarName)
		if !ok {
			continue
		}

		switch parameter.GetDefaultValue().(type) {
		case bool:
			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				return messages.New_StartupErrors_BadValueForEnvVar_Error(val, envVarName)
			}
			specifiedArgs[parameter.GetID()] = boolVal
		case string:
			specifiedArgs[parameter.GetID()] = val
		default:
			// If you hit this error, it means this switch is not implementing a supported type in `pkg/config`
			return messages.New_StartupErrors_ParseFailed_Error("\n", internalErrorText)
		}
	}
	return nil
}
