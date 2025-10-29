// Copyright 2025 The MathWorks, Inc.

// This package intentionally does not use the LoggerFactory, or any logging framework.
// All messages are read and written as raw strings, so that the main MCP Server process can log them correctly.

package watchdog

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/stdio"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
)

type OSLayer interface {
	Getppid() int
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

type ProcessHandler interface {
	WatchProcessAndGetTerminationChan(processPid int) <-chan struct{}
	KillProcess(processPid int) error
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type TransportFactory interface {
	NewReceiver(osStdio entities.OSStdio) (transport.Receiver, error)
}

type Watchdog struct {
	osLayer          OSLayer
	processHandler   ProcessHandler
	osSignaler       OSSignaler
	transportFactory TransportFactory

	parentPID         int
	processPIDsToKill map[int]struct{}
	lock              *sync.Mutex
}

func New(
	osLayer OSLayer,
	processHandler ProcessHandler,
	osSignaler OSSignaler,
	transportFactory TransportFactory,
) *Watchdog {
	return &Watchdog{
		osLayer:          osLayer,
		processHandler:   processHandler,
		osSignaler:       osSignaler,
		transportFactory: transportFactory,

		processPIDsToKill: make(map[int]struct{}),
		lock:              new(sync.Mutex),
	}
}

func (w *Watchdog) StartAndWaitForCompletion(_ context.Context) error {
	receiver, err := w.transportFactory.NewReceiver(stdio.NewOSStdio(
		w.osLayer.Stdin(),
		w.osLayer.Stdout(),
		w.osLayer.Stderr(),
	))
	if err != nil {
		return err
	}

	receiver.SendDebugMessage("Watchdog process has started")
	defer receiver.SendDebugMessage("Watchdog process has exited")

	w.parentPID = w.osLayer.Getppid()

	shutdownMessageProcessingC := make(chan struct{})
	defer close(shutdownMessageProcessingC)

	shutdownC := make(chan struct{})
	go func() {
		c := receiver.C()
		for {
			select {
			case <-shutdownMessageProcessingC:
				return
			case rawMessage, ok := <-c:
				if !ok {
					receiver.SendErrorMessage("Receiver channel closed unexpectedly")
					return
				}
				if abort := w.processIncomingMessage(receiver, rawMessage); abort {
					close(shutdownC)
					return
				}
			}
		}
	}()

	select {
	case <-shutdownC:
		defer func() {
			receiver.SendDebugMessage("Graceful shutdown completed")
			err := receiver.SendGracefulShutdownCompleted()
			if err != nil {
				receiver.SendErrorMessage("Failed to send graceful shutdown completed signal")
			}
		}()
		receiver.SendDebugMessage("Graceful shutdown signal received")

	case <-w.processHandler.WatchProcessAndGetTerminationChan(w.parentPID):
		receiver.SendDebugMessage("Lost connection to parent, shutting down")

	case <-w.osSignaler.InterruptSignalChan():
		receiver.SendDebugMessage("Received unexpected graceful shutdown OS signal")
	}

	w.terminateAllProcesses(receiver)

	return nil
}

func (w *Watchdog) processIncomingMessage(receiver transport.Receiver, rawMessage transport.Message) (abort bool) {
	w.lock.Lock()
	defer w.lock.Unlock()

	abort = false

	switch message := rawMessage.(type) {
	case transport.ProcessToKill:
		receiver.SendDebugMessage(fmt.Sprintf("Adding process %d to kill", message.PID))
		w.processPIDsToKill[message.PID] = struct{}{}
	case transport.Shutdown:
		abort = true
	}

	return
}

func (w *Watchdog) terminateAllProcesses(receiver transport.Receiver) {
	w.lock.Lock()
	defer w.lock.Unlock()

	receiver.SendDebugMessage(fmt.Sprintf("Trying to terminate %d children", len(w.processPIDsToKill)))

	for pid := range w.processPIDsToKill {
		receiver.SendDebugMessage(fmt.Sprintf("Killing process with PID %d", pid))

		if err := w.processHandler.KillProcess(pid); err != nil {
			receiver.SendErrorMessage(fmt.Sprintf("Failed to kill child with PID: %d", pid))
		}
	}
}
