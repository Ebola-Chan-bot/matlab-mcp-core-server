// Copyright 2025 The MathWorks, Inc.

package httpserverfactory

import (
	"context"
	"net"
	"net/http"
	"time"
)

const defaultReadHeaderTimeout = 10 * time.Second

type HttpServer interface {
	Serve(socketPath string) error
	Shutdown(ctx context.Context) error
}

type OSLayer interface {
	RemoveAll(name string) error
}

type HTTPServerFactory struct {
	osLayer OSLayer
}

func New(osLayer OSLayer) *HTTPServerFactory {
	return &HTTPServerFactory{
		osLayer: osLayer,
	}
}

func (f *HTTPServerFactory) NewServerOverUDS(handlers map[string]http.HandlerFunc) (HttpServer, error) {
	mux := http.NewServeMux()
	for pattern, handler := range handlers {
		mux.HandleFunc(pattern, handler)
	}

	return &udsServer{
		httpServer: &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		},
		osLayer: f.osLayer,
	}, nil
}

type udsServer struct {
	httpServer *http.Server
	osLayer    OSLayer
	socketPath string
}

func (s *udsServer) Serve(socketPath string) error {
	if err := s.osLayer.RemoveAll(socketPath); err != nil {
		return err
	}

	s.socketPath = socketPath
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *udsServer) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	if s.socketPath == "" {
		return nil
	}

	if err := s.osLayer.RemoveAll(s.socketPath); err != nil {
		return err
	}

	return nil
}
