// Copyright 2025 The MathWorks, Inc.

package httpclientfactory

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type HTTPClientFactory struct{}

func New() *HTTPClientFactory {
	return &HTTPClientFactory{}
}

func (f *HTTPClientFactory) NewClientForSelfSignedTLSServer(certificatePEM []byte) (HttpClient, error) {
	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM(certificatePEM); !ok {
		return nil, fmt.Errorf("failed to append certificate to pool")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		},
	}

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &http.Client{
		Transport: transport,
		Jar:       jar,
	}, nil
}
