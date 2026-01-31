// Copyright 2025-2026 The MathWorks, Inc.

package globalmatlab

import (
	"context"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type MATLABManager interface {
	StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error)
	StopMATLABSession(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) error
	GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error)
}

type MATLABRootSelector interface {
	SelectMATLABRoot(ctx context.Context, logger entities.Logger) (string, error)
}

type MATLABStartingDirSelector interface {
	SelectMatlabStartingDir() (string, error)
}

type SessionDiscovery interface {
	DiscoverExistingSession(logger entities.Logger) *embeddedconnector.ConnectionDetails
}

type EmbeddedConnectorClientFactory interface {
	New(endpoint embeddedconnector.ConnectionDetails) (entities.MATLABSessionClient, error)
}

type GlobalMATLAB struct {
	matlabManager             MATLABManager
	matlabRootSelector        MATLABRootSelector
	matlabStartingDirSelector MATLABStartingDirSelector
	sessionDiscovery          SessionDiscovery
	clientFactory             EmbeddedConnectorClientFactory

	lock *sync.Mutex

	initOnce  *sync.Once
	initError error

	matlabRoot        string
	matlabStartingDir string
	sessionID         entities.SessionID
	discoveredClient  entities.MATLABSessionClient
}

func New(
	matlabManager MATLABManager,
	matlabRootSelector MATLABRootSelector,
	matlabStartingDirSelector MATLABStartingDirSelector,
	sessionDiscovery SessionDiscovery,
	clientFactory EmbeddedConnectorClientFactory,
) *GlobalMATLAB {
	return &GlobalMATLAB{
		matlabManager:             matlabManager,
		matlabRootSelector:        matlabRootSelector,
		matlabStartingDirSelector: matlabStartingDirSelector,
		sessionDiscovery:          sessionDiscovery,
		clientFactory:             clientFactory,

		lock:     &sync.Mutex{},
		initOnce: &sync.Once{},
	}
}

func (g *GlobalMATLAB) Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	// 如果已有发现的客户端，直接返回
	if g.discoveredClient != nil {
		return g.discoveredClient, nil
	}

	// 首次调用时尝试发现已有会话
	g.initOnce.Do(func() {
		// 先尝试发现已存在的会话
		if connectionDetails := g.sessionDiscovery.DiscoverExistingSession(logger); connectionDetails != nil {
			logger.With("host", connectionDetails.Host).
				With("port", connectionDetails.Port).
				Info("Found existing MATLAB session, attempting to connect")

			client, err := g.clientFactory.New(*connectionDetails)
			if err != nil {
				logger.WithError(err).Warn("Failed to create client for discovered session, will start new session")
			} else {
				// 验证连接是否工作
				pingResult := client.Ping(ctx, logger)
				if pingResult.IsAlive {
					logger.Info("Successfully connected to existing MATLAB session")
					g.discoveredClient = client
					return
				}
				logger.Warn("Discovered session not responding, will start new session")
			}
		}

		// 没有发现的会话或连接失败，初始化启动配置
		err := g.initializeStartupConfig(ctx, logger)
		if err != nil {
			g.initError = err
		}
	})

	if g.discoveredClient != nil {
		return g.discoveredClient, nil
	}

	if g.initError != nil {
		return nil, g.initError
	}

	return g.getOrCreateClient(ctx, logger)
}

func (g *GlobalMATLAB) getOrCreateClient(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	var sessionIDZeroValue entities.SessionID

	// Start MATLAB if we don't have a session
	if g.sessionID == sessionIDZeroValue {
		if err := g.startNewSession(ctx, logger); err != nil {
			g.initError = err
			return nil, err
		}
	}

	// Try to get the client
	client, err := g.matlabManager.GetMATLABSessionClient(ctx, logger, g.sessionID)
	if err != nil {
		// Retry: stop old session and start a new one
		if stopErr := g.matlabManager.StopMATLABSession(ctx, logger, g.sessionID); stopErr != nil {
			logger.WithError(stopErr).Warn("failed to stop MATLAB session")
		}

		if err := g.startNewSession(ctx, logger); err != nil {
			g.initError = err
			return nil, err
		}

		return g.matlabManager.GetMATLABSessionClient(ctx, logger, g.sessionID)
	}

	return client, nil
}

func (g *GlobalMATLAB) startNewSession(ctx context.Context, logger entities.Logger) error {
	sessionID, err := g.matlabManager.StartMATLABSession(ctx, logger, entities.LocalSessionDetails{
		MATLABRoot:             g.matlabRoot,
		IsStartingDirectorySet: g.matlabStartingDir != "",
		StartingDirectory:      g.matlabStartingDir,
		ShowMATLABDesktop:      true,
	})
	if err != nil {
		return err
	}

	g.sessionID = sessionID
	return nil
}

func (g *GlobalMATLAB) initializeStartupConfig(ctx context.Context, logger entities.Logger) error {
	matlabRoot, err := g.matlabRootSelector.SelectMATLABRoot(ctx, logger)
	if err != nil {
		return err
	}

	g.matlabRoot = matlabRoot

	matlabStartingDirectory, err := g.matlabStartingDirSelector.SelectMatlabStartingDir()
	if err != nil {
		logger.WithError(err).Warn("failed to determine MATLAB starting directory, proceeding without one")
		return nil
	}

	g.matlabStartingDir = matlabStartingDirectory
	return nil
}
