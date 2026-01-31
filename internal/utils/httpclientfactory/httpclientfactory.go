// Copyright 2025 The MathWorks, Inc.

package httpclientfactory

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
	CloseIdleConnections()
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

// NewClientInsecureSkipVerify 创建一个跳过 TLS 证书验证的 HTTP 客户端。
// 用于连接到手动启动的 MATLAB 会话（没有证书文件时）。
func (f *HTTPClientFactory) NewClientInsecureSkipVerify() (HttpClient, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true, // nolint:gosec // 用于手动会话连接，用户已知风险
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

func (f *HTTPClientFactory) NewClientOverUDS(socketPath string) HttpClient {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", socketPath)
		},
	}

	return &http.Client{
		Transport: transport,
	}
}
