// Copyright 2025 The MathWorks, Inc.

package embeddedconnector

import "github.com/matlab/matlab-mcp-core-server/internal/utils/httpclientfactory"

func (c *Client) SetHttpClient(httpClient httpclientfactory.HttpClient) {
	c.httpClient = httpClient
}
