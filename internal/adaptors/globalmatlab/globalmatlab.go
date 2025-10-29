// Copyright 2025 The MathWorks, Inc.

package globalmatlab

import (
	"context"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type MATLABManager interface {
	StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error)
	GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error)
}

type MATLABRootSelector interface {
	SelectFirstMATLABVersionOnPath(ctx context.Context, logger entities.Logger) (string, error)
}

type MATLABStartingDirSelector interface {
	SelectMatlabStartingDir() (string, error)
}

type GlobalMATLAB struct {
	matlabManager             MATLABManager
	matlabRootSelector        MATLABRootSelector
	matlabStartingDirSelector MATLABStartingDirSelector

	lock              *sync.Mutex
	matlabRoot        string
	matlabStartingDir string
	sessionID         entities.SessionID
	cachedStartErr    error
}

func New(
	matlabManager MATLABManager,
	matlabRootSelector MATLABRootSelector,
	matlabStartingDirSelector MATLABStartingDirSelector,
) *GlobalMATLAB {
	return &GlobalMATLAB{
		matlabManager:             matlabManager,
		matlabRootSelector:        matlabRootSelector,
		matlabStartingDirSelector: matlabStartingDirSelector,

		lock: &sync.Mutex{},
	}
}

func (g *GlobalMATLAB) Initialize(ctx context.Context, logger entities.Logger) error {
	var err error
	g.matlabRoot, err = g.matlabRootSelector.SelectFirstMATLABVersionOnPath(ctx, logger)
	if err != nil {
		return err
	}

	g.matlabStartingDir, err = g.matlabStartingDirSelector.SelectMatlabStartingDir()
	if err != nil {
		logger.WithError(err).Warn("failed to determine MATLAB starting directory, proceeding without one")
	}

	err = g.ensureMATLABClientIsValid(ctx, logger)
	if err != nil {
		return err
	}

	return nil
}

func (g *GlobalMATLAB) Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	if err := g.ensureMATLABClientIsValid(ctx, logger); err != nil {
		return nil, err
	}

	client, err := g.matlabManager.GetMATLABSessionClient(ctx, logger, g.sessionID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (g *GlobalMATLAB) ensureMATLABClientIsValid(ctx context.Context, logger entities.Logger) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.cachedStartErr != nil {
		return g.cachedStartErr
	}

	var sessionIDZeroValue entities.SessionID
	if g.sessionID == sessionIDZeroValue {
		sessionID, err := g.matlabManager.StartMATLABSession(ctx, logger, entities.LocalSessionDetails{
			MATLABRoot:        g.matlabRoot,
			StartingDirectory: g.matlabStartingDir,
			ShowMATLABDesktop: true,
		})
		if err != nil {
			g.cachedStartErr = err
			return err
		}

		g.sessionID = sessionID
	}

	return nil
}
