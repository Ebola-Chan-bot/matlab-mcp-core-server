// Copyright 2025 The MathWorks, Inc.

package matlabmanager

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

func (m *MATLABManager) GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error) {
	return m.sessionStore.Get(sessionID)
}
