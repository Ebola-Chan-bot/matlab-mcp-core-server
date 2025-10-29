// Copyright 2025 The MathWorks, Inc.

package transport

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type stdioReceiver struct {
	stdin  entities.Reader
	stdout entities.Writer
	stderr entities.Writer

	messagesC chan Message
}

func NewStdioReceiver(osStdio entities.OSStdio) (*stdioReceiver, error) {
	stdin := osStdio.Stdin()
	stdout := osStdio.Stdout()
	stderr := osStdio.Stderr()

	if stdin == nil || stdout == nil || stderr == nil {
		return nil, fmt.Errorf("invalid os stdio: stdin, stdout, or stderr is nil")
	}

	client := &stdioReceiver{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,

		messagesC: make(chan Message),
	}

	go client.readMessagesFromStdin()

	return client, nil
}

func (r *stdioReceiver) SendDebugMessage(message string) {
	_, err := fmt.Fprintf(r.stdout, "%s\n", message)
	if err != nil {
		r.SendErrorMessage(err.Error())
	}
}

func (r *stdioReceiver) SendErrorMessage(message string) {
	_, err := fmt.Fprintf(r.stderr, "%s\n", message)
	if err != nil { //nolint:staticcheck // Readability
		// We're out of luck
	}
}

func (r *stdioReceiver) C() <-chan Message {
	return r.messagesC
}

func (r *stdioReceiver) SendGracefulShutdownCompleted() error {
	_, err := fmt.Fprintf(r.stdout, "%s\n", gracefulShutdownCompletedSignal)
	return err
}

func (r *stdioReceiver) readMessagesFromStdin() {
	scanner := bufio.NewScanner(r.stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch line {
		case gracefulShutdownSignal:
			// A shutdown was requested, signal and exit
			r.messagesC <- Shutdown{}
			return
		default:
			// Expect process PIDs
			processPid, err := strconv.Atoi(line)
			if err != nil {
				r.SendErrorMessage(fmt.Errorf("failed to cast message \"%s\" to int. %w", line, err).Error())
				continue
			}

			r.messagesC <- ProcessToKill{PID: processPid}
		}
	}
}
