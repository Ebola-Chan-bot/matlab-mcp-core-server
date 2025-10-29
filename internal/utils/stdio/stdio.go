// Copyright 2025 The MathWorks, Inc.

package stdio

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type osStdio struct {
	stdin  entities.Reader
	stdout entities.Writer
	stderr entities.Writer
}

type subProcessStdio struct {
	stdin  entities.Writer
	stdout entities.Reader
	stderr entities.Reader
}

func NewOSStdio(
	stdin entities.Reader,
	stdout entities.Writer,
	stderr entities.Writer,
) *osStdio {
	return &osStdio{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func NewSubProcessStdio(
	stdin entities.Writer,
	stdout entities.Reader,
	stderr entities.Reader,
) *subProcessStdio {
	return &subProcessStdio{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (s *osStdio) Stdin() entities.Reader {
	return s.stdin
}

func (s *osStdio) Stdout() entities.Writer {
	return s.stdout
}

func (s *osStdio) Stderr() entities.Writer {
	return s.stderr
}

func (s *subProcessStdio) Stdin() entities.Writer {
	return s.stdin
}

func (s *subProcessStdio) Stdout() entities.Reader {
	return s.stdout
}

func (s *subProcessStdio) Stderr() entities.Reader {
	return s.stderr
}
