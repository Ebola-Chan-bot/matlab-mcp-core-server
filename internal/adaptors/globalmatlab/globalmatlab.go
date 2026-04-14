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

	// 如果已有发现的客户端，验证其是否仍然存活
	if g.discoveredClient != nil {
		pingResult := g.discoveredClient.Ping(ctx, logger)
		if pingResult.IsAlive {
			return g.discoveredClient, nil
		}
		logger.Warn("之前发现的 MATLAB 会话已无响应，尝试重新发现")
		g.discoveredClient = nil
	}

	// 首次调用时初始化启动配置（非 existing-session-only 模式需要）
	g.initOnce.Do(func() {
		existingOnly := g.isExistingSessionOnly()
		if !existingOnly {
			err := g.initializeStartupConfig(ctx, logger)
			if err != nil {
				g.initError = err
			}
		}
	})

	// 尝试发现已存在的会话
	if client, err := g.tryDiscoverExistingSession(ctx, logger); err != nil {
		return nil, err
	} else if client != nil {
		return client, nil
	}

	// 没有发现可用的已有会话
	if g.isExistingSessionOnly() {
		return nil, fmt.Errorf("existing-session-only mode: no existing MATLAB session found. Please run 'RegisterMatlabSession' in MATLAB first to register a session for discovery")
	}

	if g.initError != nil {
		return nil, g.initError
	}

	return g.getOrCreateClient(ctx, logger)
}

// tryDiscoverExistingSession 尝试发现并连接到一个已存在的 MATLAB 会话。
// 成功返回 (client, nil)，未找到返回 (nil, nil)，致命错误返回 (nil, err)。
func (g *GlobalMATLAB) tryDiscoverExistingSession(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	connectionDetails := g.sessionDiscovery.DiscoverExistingSession(logger)
	if connectionDetails == nil {
		return nil, nil
	}

	logger.With("host", connectionDetails.Host).
		With("port", connectionDetails.Port).
		Info("发现已有 MATLAB 会话，正在尝试连接")

	client, err := g.clientFactory.New(*connectionDetails)
	if err != nil {
		if g.isExistingSessionOnly() {
			return nil, fmt.Errorf("existing-session-only 模式：在 %s:%s 发现会话但创建客户端失败: %w", connectionDetails.Host, connectionDetails.Port, err)
		}
		logger.WithError(err).Warn("为已发现的会话创建客户端失败")
		return nil, nil
	}

	pingResult := client.Ping(ctx, logger)
	if !pingResult.IsAlive {
		if g.isExistingSessionOnly() {
			return nil, fmt.Errorf("existing-session-only 模式：在 %s:%s 发现会话但 ping 无响应", connectionDetails.Host, connectionDetails.Port)
		}
		logger.Warn("已发现的会话无响应")
		return nil, nil
	}

	logger.Info("已成功连接到已有 MATLAB 会话")
	if err := g.deployHelperFiles(ctx, logger, client); err != nil {
		logger.WithError(err).Warn("向已有会话部署辅助文件失败，EvalWithCapture 可能无法工作")
	}
	g.discoveredClient = client
	return client, nil
}

func (g *GlobalMATLAB) isExistingSessionOnly() bool {
	if cfg, cfgErr := g.configFactory.Config(); cfgErr == nil {
		return cfg.ExistingSessionOnly()
	}
	return false
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
