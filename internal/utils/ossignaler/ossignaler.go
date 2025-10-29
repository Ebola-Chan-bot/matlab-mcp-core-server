// Copyright 2025 The MathWorks, Inc.

package ossignaler

import (
	"os"
	"os/signal"
	"syscall"
)

type OSSignaler struct {
}

func New() *OSSignaler {
	return &OSSignaler{}
}

// GetInterruptSignalChan wraps the signal.Notify function to get a signal interrupt.
func (osw *OSSignaler) InterruptSignalChan() <-chan os.Signal {
	interruptC := make(chan os.Signal, 1)
	signal.Notify(interruptC, os.Interrupt, syscall.SIGTERM)
	return interruptC
}
