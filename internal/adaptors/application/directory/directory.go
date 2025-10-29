// Copyright 2025 The MathWorks, Inc.

package directory

const defaultLogDirPattern = "matlab-mcp-core-server-"

type OSLayer interface {
	MkdirTemp(dir string, pattern string) (string, error)
}

type Directory struct {
	logDir string

	osFacade OSLayer
}

func New(
	osFacade OSLayer,
) (*Directory, error) {
	logDir, err := osFacade.MkdirTemp("", defaultLogDirPattern)
	if err != nil {
		return nil, err
	}

	return &Directory{
		logDir:   logDir,
		osFacade: osFacade,
	}, nil
}

func (d *Directory) BaseDir() string {
	return d.logDir
}

func (d *Directory) MkdirTemp(pattern string) (string, error) {
	return d.osFacade.MkdirTemp(d.logDir, pattern)
}
