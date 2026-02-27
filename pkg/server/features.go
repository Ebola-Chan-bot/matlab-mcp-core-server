// Copyright 2026 The MathWorks, Inc.

package server

import "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"

type Features struct {
	MATLAB MATLABFeature
}

type MATLABFeature struct {
	Enabled bool
}

func (f Features) toInternal() definition.Features {
	return definition.Features{
		MATLAB: definition.MATLABFeature{
			Enabled: f.MATLAB.Enabled,
		},
	}
}
