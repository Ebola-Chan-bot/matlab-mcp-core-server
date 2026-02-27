// Copyright 2026 The MathWorks, Inc.

package server

func NewI18nErrorFactory(messageCatalog i18nMessageCatalog) i18nErrorFactory {
	return newI18nErrorFactory(messageCatalog)
}
