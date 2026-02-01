// Copyright 2025-2026 The MathWorks, Inc.

package sessiondiscovery

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const (
	appDirPattern     = "matlab-mcp-core-server-"
	sessionDirPattern = "matlab-session-"
	securePortFile    = "connector.securePort"
	certificateFile   = "cert.pem"
	apiKeyFile        = "apikey"
)

type OSLayer interface {
	TempDir() string
	Getenv(key string) string
	ReadDir(name string) ([]os.DirEntry, error)
	ReadFile(name string) ([]byte, error)
}

type SessionDiscovery struct {
	osLayer OSLayer
}

func New(osLayer OSLayer) *SessionDiscovery {
	return &SessionDiscovery{
		osLayer: osLayer,
	}
}

// DiscoverExistingSession 在临时目录中搜索已有的 MATLAB MCP 会话。
// 返回找到的第一个有效会话，如果不存在则返回 nil。
func (s *SessionDiscovery) DiscoverExistingSession(logger entities.Logger) *embeddedconnector.ConnectionDetails {
	// 使用 LOCALAPPDATA\Temp 作为共享临时目录（Windows）
	// 不同会话的 TEMP 环境变量可能不同，但 LOCALAPPDATA 是固定的
	tempDir := s.osLayer.Getenv("LOCALAPPDATA")
	if tempDir != "" {
		tempDir = filepath.Join(tempDir, "Temp")
	} else {
		// 非 Windows 系统回退到默认临时目录
		tempDir = s.osLayer.TempDir()
	}

	logger.With("temp_dir", tempDir).Debug("Searching for existing MATLAB sessions")

	// 列出临时目录中所有匹配应用目录模式的目录
	entries, err := s.osLayer.ReadDir(tempDir)
	if err != nil {
		logger.WithError(err).Debug("Failed to read temp directory")
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), appDirPattern) {
			continue
		}

		appDir := filepath.Join(tempDir, entry.Name())
		logger.With("app_dir", appDir).Debug("Found potential app directory")

		// 在应用目录中搜索会话目录
		sessionDetails := s.searchSessionsInAppDir(logger, appDir)
		if sessionDetails != nil {
			return sessionDetails
		}
	}

	logger.Debug("No existing MATLAB sessions found")
	return nil
}

func (s *SessionDiscovery) searchSessionsInAppDir(logger entities.Logger, appDir string) *embeddedconnector.ConnectionDetails {
	entries, err := s.osLayer.ReadDir(appDir)
	if err != nil {
		logger.WithError(err).Debug("Failed to read app directory")
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), sessionDirPattern) {
			continue
		}

		sessionDir := filepath.Join(appDir, entry.Name())
		logger.With("session_dir", sessionDir).Debug("Found potential session directory")

		details := s.tryReadSessionDetails(logger, sessionDir)
		if details != nil {
			logger.With("session_dir", sessionDir).Info("Discovered existing MATLAB session")
			return details
		}
	}

	return nil
}

func (s *SessionDiscovery) tryReadSessionDetails(logger entities.Logger, sessionDir string) *embeddedconnector.ConnectionDetails {
	securePortPath := filepath.Join(sessionDir, securePortFile)
	certPath := filepath.Join(sessionDir, certificateFile)
	apiKeyPath := filepath.Join(sessionDir, apiKeyFile)

	// 读取必需的会话文件
	securePort, err := s.osLayer.ReadFile(securePortPath)
	if err != nil || len(securePort) == 0 {
		logger.Debug("Session missing or invalid securePort file")
		return nil
	}

	apiKey, err := s.osLayer.ReadFile(apiKeyPath)
	if err != nil || len(apiKey) == 0 {
		logger.Debug("Session missing or invalid apikey file")
		return nil
	}

	// 证书是可选的 - 如果缺失，将使用 InsecureSkipVerify 模式
	certificatePEM, err := s.osLayer.ReadFile(certPath)
	if err != nil || len(certificatePEM) == 0 {
		logger.Debug("Session missing certificate file, will use InsecureSkipVerify")
		certificatePEM = nil
	}

	return &embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           strings.TrimSpace(string(securePort)),
		APIKey:         strings.TrimSpace(string(apiKey)),
		CertificatePEM: certificatePEM,
	}
}
