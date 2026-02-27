// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependenciesProvider_toInternal_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedMessage := "test message"
	expectedKey := "test-key"
	expectedValue := "test-value"

	mockLogger.EXPECT().
		Info(expectedMessage).
		Once()

	mockConfig.EXPECT().
		Get(expectedKey).
		Return(expectedValue, nil).
		Once()

	type TestDependencies struct{}
	expectedDependencies := &TestDependencies{}

	provider := server.DependenciesProvider[*TestDependencies](func(resources server.DependenciesProviderResources) (*TestDependencies, i18n.Error) {
		resources.Logger().Info(expectedMessage)

		result, err := resources.Config().Get(expectedKey, "")
		require.NoError(t, err)
		assert.Equal(t, expectedValue, result)

		return expectedDependencies, nil
	})

	// Act
	internalProvider := provider.ToInternal()
	dependencies, err := internalProvider(definition.NewDependenciesProviderResources(
		mockLogger,
		mockConfig,
		mockMessageCatalog,
	))

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedDependencies, dependencies)
}

func TestDependenciesProvider_toInternal_NilProvider(t *testing.T) {
	// Arrange
	var provider server.DependenciesProvider[struct{}]

	// Act
	internalProvider := provider.ToInternal()
	result, err := internalProvider(definition.DependenciesProviderResources{})

	// Assert
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestDependenciesProvider_toInternal_Error(t *testing.T) {
	// Arrange
	type TestDependencies struct{}
	expectedError := AnI18nError

	provider := server.DependenciesProvider[*TestDependencies](func(resources server.DependenciesProviderResources) (*TestDependencies, i18n.Error) {
		return &TestDependencies{}, expectedError // Returning an actual pointer, to make sure we mask it, and return nil
	})

	// Act
	internalProvider := provider.ToInternal()
	dependencies, err := internalProvider(definition.DependenciesProviderResources{})

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.Nil(t, dependencies)
}

var AnI18nError = &i18nError{} //nolint:gochecknoglobals // AnError is an error

type i18nError struct{}

func (e *i18nError) Error() string { return "" }

func (e *i18nError) MWMarker() {}
