// Copyright 2025 The MathWorks, Inc.

package matlabsessionclient

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpclientfactory"
)

type HttpClientFactory interface {
	NewClientForSelfSignedTLSServer(certificatePEM []byte) (httpclientfactory.HttpClient, error)
}

type Factory struct {
	httpClientFactory HttpClientFactory
}

func NewFactory(
	httpClientFactory HttpClientFactory,
) *Factory {
	return &Factory{
		httpClientFactory: httpClientFactory,
	}
}

func (f *Factory) New(endpoint embeddedconnector.ConnectionDetails) (entities.MATLABSessionClient, error) {
	return embeddedconnector.NewClient(endpoint, f.httpClientFactory)
}
