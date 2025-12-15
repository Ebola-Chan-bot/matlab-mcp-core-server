// Copyright 2025 The MathWorks, Inc.

package client

import (
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

func NewClient(
	osLayer OSLayer,
	httpClientFactory HTTPClientFactory,
	logger entities.Logger,
) *Client {
	return newClient(
		osLayer,
		httpClientFactory,
		logger,
	)
}

func (c *Client) SetSocketWaitTimeout(socketWaitTimeout time.Duration) {
	c.socketWaitTimeout = socketWaitTimeout
}

func (c *Client) SetSocketRetryInterval(socketRetryInterval time.Duration) {
	c.socketRetryInterval = socketRetryInterval
}
