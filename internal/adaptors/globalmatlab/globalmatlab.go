// Copyright 2025-2026 The MathWorks, Inc.

package globalmatlab

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type MATLABManager interface {
	StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error)
	StopMATLABSession(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) error
	GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error)
}

type MATLABRootSelector interface {
	SelectMATLABRoot(ctx context.Context, logger entities.Logger) (string, error)
}

type MATLABStartingDirSelector interface {
	SelectMatlabStartingDir(logger entities.Logger) (string, error)
}

type SessionDiscovery interface {
	DiscoverExistingSession(logger entities.Logger) *embeddedconnector.ConnectionDetails
}

type EmbeddedConnectorClientFactory interface {
	New(endpoint embeddedconnector.ConnectionDetails) (entities.MATLABSessionClient, error)
}

type MATLABFiles interface {
	GetAll() map[string][]byte
}

type GlobalMATLAB struct {
	matlabManager             MATLABManager
	matlabRootSelector        MATLABRootSelector
	matlabStartingDirSelector MATLABStartingDirSelector
	sessionDiscovery          SessionDiscovery
	clientFactory             EmbeddedConnectorClientFactory
	configFactory             ConfigFactory
	matlabFiles               MATLABFiles

	lock *sync.Mutex

	initOnce  *sync.Once
	initError error

	matlabRoot        string
	matlabStartingDir string
	sessionID         entities.SessionID
	discoveredClient  entities.MATLABSessionClient
	helperDir         string // directory where +matlab_mcp helper files are deployed
}

func New(
	matlabManager MATLABManager,
	matlabRootSelector MATLABRootSelector,
	matlabStartingDirSelector MATLABStartingDirSelector,
	sessionDiscovery SessionDiscovery,
	clientFactory EmbeddedConnectorClientFactory,
	configFactory ConfigFactory,
	matlabFiles MATLABFiles,
) *GlobalMATLAB {
	return &GlobalMATLAB{
		matlabManager:             matlabManager,
		matlabRootSelector:        matlabRootSelector,
		matlabStartingDirSelector: matlabStartingDirSelector,
		sessionDiscovery:          sessionDiscovery,
		clientFactory:             clientFactory,
		configFactory:             configFactory,
		matlabFiles:               matlabFiles,

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
		// 检查是否为仅连接已有会话模式
		existingOnly := false
		if cfg, cfgErr := g.configFactory.Config(); cfgErr == nil {
			existingOnly = cfg.ExistingSessionOnly()
		}

		// 先尝试发现已存在的会话
		if connectionDetails := g.sessionDiscovery.DiscoverExistingSession(logger); connectionDetails != nil {
			logger.With("host", connectionDetails.Host).
				With("port", connectionDetails.Port).
				Info("Found existing MATLAB session, attempting to connect")

			client, err := g.clientFactory.New(*connectionDetails)
			if err != nil {
				if existingOnly {
					g.initError = fmt.Errorf("existing-session-only mode: found session at %s:%s but failed to create client: %w", connectionDetails.Host, connectionDetails.Port, err)
					return
				}
				logger.WithError(err).Warn("Failed to create client for discovered session, will start new session")
			} else {
				// 验证连接是否工作
				pingResult := client.Ping(ctx, logger)
				if pingResult.IsAlive {
					logger.Info("Successfully connected to existing MATLAB session")
					// Deploy +matlab_mcp helper files for EvalWithCapture support
					if err := g.deployHelperFiles(ctx, logger, client); err != nil {
						logger.WithError(err).Warn("Failed to deploy helper files to existing session, EvalWithCapture may not work")
					}
					g.discoveredClient = client
					return
				}
				if existingOnly {
					g.initError = fmt.Errorf("existing-session-only mode: found session at %s:%s but it is not responding to ping", connectionDetails.Host, connectionDetails.Port)
					return
				}
				logger.Warn("Discovered session not responding, will start new session")
			}
		} else if existingOnly {
			g.initError = fmt.Errorf("existing-session-only mode: no existing MATLAB session found. Please run 'RegisterMatlabSession' in MATLAB first to register a session for discovery")
			return
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
	config, messagesErr := g.configFactory.Config()
	if messagesErr != nil {
		return messagesErr
	}

	sessionID, err := g.matlabManager.StartMATLABSession(ctx, logger, entities.LocalSessionDetails{
		MATLABRoot:             g.matlabRoot,
		IsStartingDirectorySet: g.matlabStartingDir != "",
		StartingDirectory:      g.matlabStartingDir,
		ShowMATLABDesktop:      config.ShouldShowMATLABDesktop(),
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

	matlabStartingDirectory, err := g.matlabStartingDirSelector.SelectMatlabStartingDir(logger)
	if err != nil {
		logger.WithError(err).Warn("failed to determine MATLAB starting directory, proceeding without one")
		return nil
	}

	g.matlabStartingDir = matlabStartingDirectory
	return nil
}

// deployHelperFiles writes the +matlab_mcp helper files (mcpEval.m, getOrStashExceptions.m, etc.)
// to a temporary directory and adds it to the MATLAB path via client.Eval.
// This is needed for existing-session-only mode where MATLAB wasn't started by the server.
func (g *GlobalMATLAB) deployHelperFiles(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) error {
	// Create temp dir for helper files
	helperDir, err := os.MkdirTemp("", "matlab-mcp-helpers-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory for helper files: %w", err)
	}

	packageDir := filepath.Join(helperDir, "+matlab_mcp")
	if err := os.Mkdir(packageDir, 0o700); err != nil {
		return fmt.Errorf("failed to create +matlab_mcp package directory: %w", err)
	}

	for fileName, fileContent := range g.matlabFiles.GetAll() {
		filePath := filepath.Join(packageDir, fileName)
		if err := os.WriteFile(filePath, fileContent, 0o600); err != nil {
			return fmt.Errorf("failed to write %s: %w", fileName, err)
		}
	}

	g.helperDir = helperDir
	logger.With("path", helperDir).Info("Deployed +matlab_mcp helper files")

	// Add the helper directory to MATLAB path
	addpathCode := fmt.Sprintf("addpath('%s');", strings.ReplaceAll(helperDir, "'", "''"))
	_, err = client.Eval(ctx, logger, entities.EvalRequest{Code: addpathCode})
	if err != nil {
		return fmt.Errorf("failed to addpath for helper files: %w", err)
	}

	logger.Info("Successfully added helper files to MATLAB path")
	return nil
}
