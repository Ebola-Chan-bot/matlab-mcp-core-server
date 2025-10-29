// Copyright 2025 The MathWorks, Inc.

package transport

import "time"

func (c *stdioClient) SetShutdownTimeout(timeout time.Duration) {
	c.shutdownTimeout = timeout
}
