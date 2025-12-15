// Copyright 2025 The MathWorks, Inc.

package socket

import (
	"errors"
	"path/filepath"
)

var ErrSocketPathTooLong = errors.New("socket path is too long")

type Directory interface {
	BaseDir() string
	ID() string
}

type Socket interface {
	Path() string
}

type Factory struct {
	directory Directory

	socketInstance Socket
	socketError    error
}

func NewFactory(
	directory Directory,
) *Factory {
	return &Factory{
		directory: directory,
	}
}

func (f *Factory) Socket() (Socket, error) {
	if f.socketError != nil {
		return nil, f.socketError
	}

	if f.socketInstance == nil {
		socket, err := newSocket(
			f.directory,
		)
		if err != nil {
			f.socketError = err
			return nil, err
		}

		f.socketInstance = socket
	}

	return f.socketInstance, nil
}

type socket struct {
	path string
}

func newSocket(
	directory Directory,
) (*socket, error) {
	socketPath := filepath.Join(directory.BaseDir(), "watchdog-"+directory.ID()+".sock")

	// Socket path max length is 108 characters, but for safety using 105
	if len(socketPath) > 105 {
		return nil, ErrSocketPathTooLong
	}

	return &socket{
		path: socketPath,
	}, nil
}

func (s *socket) Path() string {
	return s.path
}
