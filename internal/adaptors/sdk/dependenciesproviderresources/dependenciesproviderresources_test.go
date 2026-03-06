// Copyright 2026 The MathWorks, Inc.

package dependenciesproviderresources_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sdk/dependenciesproviderresources"
	internalconfigmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	dependenciesproviderresourcesmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/sdk/dependenciesproviderresources"
	publictypesmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/sdk/publictypes"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &dependenciesproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &dependenciesproviderresourcesmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	// Act
	factory := dependenciesproviderresources.NewFactory(mockLoggerFactory, mockConfigFactory)

	// Assert
	require.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &dependenciesproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &dependenciesproviderresourcesmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &internalconfigmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(nil).
		Once()

	mockConfigFactory.EXPECT().
		New(mockInternalConfig, mockMessageCatalog).
		Return(nil).
		Once()

	internalResources := definition.NewDependenciesProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
	)

	// Act
	resources := dependenciesproviderresources.NewFactory(mockLoggerFactory, mockConfigFactory).New(internalResources)

	// Assert
	require.NotNil(t, resources)
}

func TestFactory_New_Logger(t *testing.T) {
	// Arrange
	mockLoggerFactory := &dependenciesproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &dependenciesproviderresourcesmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &internalconfigmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedLogger := &publictypesmocks.MockLogger{}
	defer expectedLogger.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(expectedLogger).
		Once()

	mockConfigFactory.EXPECT().
		New(mockInternalConfig, mockMessageCatalog).
		Return(nil).
		Once()

	internalResources := definition.NewDependenciesProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
	)

	// Act
	resources := dependenciesproviderresources.NewFactory(mockLoggerFactory, mockConfigFactory).New(internalResources)

	// Assert
	require.Equal(t, expectedLogger, resources.Logger())
}

func TestFactory_New_Config(t *testing.T) {
	// Arrange
	mockLoggerFactory := &dependenciesproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &dependenciesproviderresourcesmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &internalconfigmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedConfig := &publictypesmocks.MockConfig{}
	defer expectedConfig.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(nil).
		Once()

	mockConfigFactory.EXPECT().
		New(mockInternalConfig, mockMessageCatalog).
		Return(expectedConfig).
		Once()

	internalResources := definition.NewDependenciesProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
	)

	// Act
	resources := dependenciesproviderresources.NewFactory(mockLoggerFactory, mockConfigFactory).New(internalResources)

	// Assert
	require.Equal(t, expectedConfig, resources.Config())
}
