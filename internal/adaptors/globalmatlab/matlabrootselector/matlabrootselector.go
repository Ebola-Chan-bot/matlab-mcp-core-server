// Copyright 2025 The MathWorks, Inc.

package matlabrootselector

import (
	"context"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type Config interface {
	PreferredLocalMATLABRoot() string
}

type MATLABManager interface {
	ListEnvironments(ctx context.Context, sessionLogger entities.Logger) []entities.EnvironmentInfo
}

type MATLABRootSelector struct {
	config        Config
	matlabManager MATLABManager
}

func New(
	config Config,
	matlabManager MATLABManager,
) *MATLABRootSelector {
	return &MATLABRootSelector{
		config:        config,
		matlabManager: matlabManager,
	}
}

func (m *MATLABRootSelector) SelectFirstMATLABVersionOnPath(ctx context.Context, logger entities.Logger) (string, error) {
	if preferredLocalMATLABRoot := m.config.PreferredLocalMATLABRoot(); preferredLocalMATLABRoot != "" {
		return preferredLocalMATLABRoot, nil
	}

	environments := m.matlabManager.ListEnvironments(ctx, logger)
	if len(environments) == 0 {
		return "", fmt.Errorf("no valid MATLAB environments found")
	}

	return environments[0].MATLABRoot, nil
}
