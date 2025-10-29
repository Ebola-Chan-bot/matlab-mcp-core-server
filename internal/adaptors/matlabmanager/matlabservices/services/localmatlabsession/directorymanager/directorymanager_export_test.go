// Copyright 2025 The MathWorks, Inc.

package directorymanager

import (
	"time"
)

func NewDirectoryManager(sessionDir string, osLayer OSLayer) *directoryManager {
	return newDirectoryManager(sessionDir, osLayer)
}

func (m *directoryManager) SecurePortFile() string {
	return m.securePortFile()
}

func (m *directoryManager) SetEmbeddedConnectorDetailsTimeout(timeout time.Duration) {
	m.embeddedConnectorDetailsTimeout = timeout
}

func (m *directoryManager) SetEmbeddedConnectorDetailsRetry(retry time.Duration) {
	m.embeddedConnectorDetailsRetry = retry
}

func (m *directoryManager) SetCleanupTimeout(timeout time.Duration) {
	m.cleanupTimeout = timeout
}

func (m *directoryManager) SetCleanupRetry(retry time.Duration) {
	m.cleanupRetry = retry
}
