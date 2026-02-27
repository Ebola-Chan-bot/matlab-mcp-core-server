// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
)

type i18nMessageCatalog interface {
	GetFromError(err messages.Error) string
}

type i18nErrorFactory struct {
	messageCatalog i18nMessageCatalog
}

func newI18nErrorFactory(messageCatalog i18nMessageCatalog) i18nErrorFactory {
	return i18nErrorFactory{
		messageCatalog: messageCatalog,
	}
}

func (f i18nErrorFactory) FromInternalError(internalError messages.Error) i18n.Error {
	message := f.messageCatalog.GetFromError(internalError)
	return &i18nErrorFromInternalError{
		message: message,
	}
}

type i18nErrorFromInternalError struct {
	message string
}

func (e *i18nErrorFromInternalError) Error() string {
	return e.message
}

func (e *i18nErrorFromInternalError) MWMarker() {}
