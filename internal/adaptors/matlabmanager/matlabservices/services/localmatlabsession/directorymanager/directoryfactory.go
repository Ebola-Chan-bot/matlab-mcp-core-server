// Copyright 2025 The MathWorks, Inc.

package directorymanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

type ApplicationDirectory interface {
	MkdirTemp(pattern string) (string, error)
}

type OSLayer interface {
	Mkdir(name string, perm os.FileMode) error
	RemoveAll(path string) error

	Stat(name string) (osfacade.FileInfo, error)
	ReadFile(filePath string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type MATLABFiles interface {
	GetAll() map[string][]byte
}

type Directory interface {
	Path() string
	CertificateFile() string
	CertificateKeyFile() string
	GetEmbeddedConnectorDetails() (string, []byte, error)
	Cleanup() error
}

type DirectoryFactory struct {
	osLayer              OSLayer
	applicationDirectory ApplicationDirectory
	matlabFiles          MATLABFiles
}

func NewFactory(
	osLayer OSLayer,
	applicationDirectory ApplicationDirectory,
	matlabFiles MATLABFiles,
) *DirectoryFactory {
	return &DirectoryFactory{
		osLayer:              osLayer,
		applicationDirectory: applicationDirectory,
		matlabFiles:          matlabFiles,
	}
}

func (f *DirectoryFactory) Create(logger entities.Logger) (Directory, error) {
	sessionDir, err := f.applicationDirectory.MkdirTemp("matlab-session-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary session directory: %w", err)
	}

	matlabMCPPackagePath := filepath.Join(sessionDir, "+matlab_mcp")

	err = f.osLayer.Mkdir(matlabMCPPackagePath, 0o700)
	if err != nil {
		return nil, fmt.Errorf("failed to create package directory: %w", err)
	}

	for fileName, fileContent := range f.matlabFiles.GetAll() {
		filePath := filepath.Join(matlabMCPPackagePath, fileName)
		if err := f.osLayer.WriteFile(filePath, fileContent, 0o600); err != nil {
			return nil, fmt.Errorf("failed to create %s file: %w", fileName, err)
		}
	}

	return newDirectoryManager(sessionDir, f.osLayer), nil
}
