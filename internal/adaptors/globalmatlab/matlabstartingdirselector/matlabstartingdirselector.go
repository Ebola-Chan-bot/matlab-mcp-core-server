// Copyright 2025 The MathWorks, Inc.

package matlabstartingdirselector

import (
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

type Config interface {
	PreferredMATLABStartingDirectory() string
}

type OSLayer interface {
	UserHomeDir() (string, error)
	Stat(path string) (osfacade.FileInfo, error)
	Getenv(key string) string
	GOOS() string
}

type MATLABStartingDirSelector struct {
	config  Config
	osLayer OSLayer
}

func New(
	config Config,
	osLayer OSLayer,
) *MATLABStartingDirSelector {
	return &MATLABStartingDirSelector{
		config:  config,
		osLayer: osLayer,
	}
}

func (s *MATLABStartingDirSelector) SelectMatlabStartingDir() (string, error) {
	var err error
	// Try preferred directory first
	if preferredDir := s.config.PreferredMATLABStartingDirectory(); preferredDir != "" {
		if _, err := s.osLayer.Stat(preferredDir); err != nil {
			return "", err
		}
		return preferredDir, nil
	}

	// Fall back to documents directory
	dir, err := s.getDocumentsDir()
	if err != nil {
		return "", err
	}

	if _, err := s.osLayer.Stat(dir); err != nil {
		return "", err
	}

	return dir, nil
}

func (s *MATLABStartingDirSelector) getDocumentsDir() (string, error) {
	home, err := s.osLayer.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch s.osLayer.GOOS() {
	case "windows", "darwin":
		return filepath.Join(home, "Documents"), nil
	default: // Linux - Documents less commonly used
		return home, nil // Just return home for Linux
	}
}
