// Copyright 2025 The MathWorks, Inc.

package transport

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const defaultWatchdogProcessStopTimeout = 20 * time.Second

type stdioClient struct {
	stdin  entities.Writer
	stdout entities.Reader
	stderr entities.Reader

	debugMessagesC chan string
	errorMessagesC chan string

	shutdownTimeout time.Duration
	shutdownC       chan struct{}
}

func NewStdioClient(subProcessStdio entities.SubProcessStdio) (*stdioClient, error) {
	stdin := subProcessStdio.Stdin()
	stdout := subProcessStdio.Stdout()
	stderr := subProcessStdio.Stderr()

	if stdin == nil || stdout == nil || stderr == nil {
		return nil, fmt.Errorf("invalid subprocess stdio: stdin, stdout, or stderr is nil")
	}

	client := &stdioClient{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,

		debugMessagesC: make(chan string),
		errorMessagesC: make(chan string),

		shutdownTimeout: defaultWatchdogProcessStopTimeout,
		shutdownC:       make(chan struct{}),
	}

	go client.readMessagesFromStdout()
	go client.readMessagesFromStderr()

	return client, nil
}

func (c *stdioClient) SendProcessPID(processPID int) error {
	_, err := fmt.Fprintf(c.stdin, "%d\n", processPID)
	return err
}

func (c *stdioClient) SendStop() error {
	if _, err := fmt.Fprintf(c.stdin, "%s\n", gracefulShutdownSignal); err != nil {
		return err
	}

	select {
	case <-time.After(c.shutdownTimeout):
		return fmt.Errorf("timedout waiting for watchdog process to shutdown")
	case <-c.shutdownC:
		return nil
	}
}

func (c *stdioClient) DebugMessagesC() <-chan string {
	return c.debugMessagesC
}

func (c *stdioClient) ErrorMessagesC() <-chan string {
	return c.errorMessagesC
}

func (c *stdioClient) readMessagesFromStdout() {
	scanner := bufio.NewScanner(c.stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch line {
		case gracefulShutdownCompletedSignal:
			// A shutdown was initiated and completed, close channel and exit
			close(c.shutdownC)
			return
		default:
			// Just a regular message
			c.debugMessagesC <- line
		}
	}
}

func (c *stdioClient) readMessagesFromStderr() {
	scanner := bufio.NewScanner(c.stderr)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch line {
		default:
			c.errorMessagesC <- line
		}
	}
}
