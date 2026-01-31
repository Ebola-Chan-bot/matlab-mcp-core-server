// Copyright 2025-2026 The MathWorks, Inc.

package matlabsessionstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type MATLABSessionClientWithCleanup interface {
	entities.MATLABSessionClient
	StopSession(ctx context.Context, sessionLogger entities.Logger) error
}

type LifecycleSignaler interface {
	AddShutdownFunction(shutdownFcn func() error)
}

type Store struct {
	l       *sync.RWMutex
	next    entities.SessionID
	clients map[entities.SessionID]MATLABSessionClientWithCleanup
}

func New(
	loggerFactory LoggerFactory,
	lifecycleSignaler LifecycleSignaler,
) *Store {
	store := &Store{
		l:       new(sync.RWMutex),
		next:    1,
		clients: map[entities.SessionID]MATLABSessionClientWithCleanup{},
	}

	// 不再在 MCP 服务器关闭时自动关闭 MATLAB 会话
	// 这样用户可以保持 MATLAB 运行并重新连接
	_ = loggerFactory     // 保留参数以维持接口兼容性
	_ = lifecycleSignaler // 保留参数以维持接口兼容性

	return store
}

func (s *Store) Add(client MATLABSessionClientWithCleanup) entities.SessionID {
	s.l.Lock()
	defer s.l.Unlock()

	sessionID := s.next
	s.clients[sessionID] = client
	s.next++
	return entities.SessionID(sessionID)
}

func (s *Store) Get(sessionID entities.SessionID) (MATLABSessionClientWithCleanup, error) {
	s.l.RLock()
	defer s.l.RUnlock()

	client, exists := s.clients[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %v", sessionID)
	}

	return client, nil
}

func (s *Store) Remove(sessionID entities.SessionID) {
	s.l.Lock()
	defer s.l.Unlock()

	delete(s.clients, sessionID)
}
