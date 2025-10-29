// Copyright 2025 The MathWorks, Inc.

package entities

import (
	"context"
	"io"
)

type Mode interface {
	StartAndWaitForCompletion(ctx context.Context) error
}

type Reader interface {
	io.Reader
}

type Writer interface {
	io.Writer
}

type OSStdio interface {
	Stdin() Reader
	Stdout() Writer
	Stderr() Writer
}

type SubProcessStdio interface {
	Stdin() Writer
	Stdout() Reader
	Stderr() Reader
}
