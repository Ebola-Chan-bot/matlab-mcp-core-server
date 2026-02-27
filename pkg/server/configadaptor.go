// Copyright 2026 The MathWorks, Inc.

package server

import (
	"reflect"

	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/pkg/config"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
)

type configAdaptor struct {
	internalConfig internalconfig.GenericConfig
	errorFactory   i18nErrorFactory
}

func newConfigAdaptor(internalConfig internalconfig.GenericConfig, messageCatalog definition.MessageCatalog) config.Config {
	return &configAdaptor{
		internalConfig: internalConfig,
		errorFactory:   newI18nErrorFactory(messageCatalog),
	}
}

func (c *configAdaptor) Get(key string, expectedType any) (any, i18n.Error) {
	value, err := c.internalConfig.Get(key)
	if err != nil {
		return nil, c.errorFactory.FromInternalError(err)
	}

	if reflect.TypeOf(value) != reflect.TypeOf(expectedType) {
		internalError := messages.New_StartupErrors_InvalidParameterType_Error(key, reflect.TypeOf(expectedType).String())
		return nil, c.errorFactory.FromInternalError(internalError)
	}

	return value, nil
}
